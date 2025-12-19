<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useCustomersStore } from "../../stores/customers.store";
import { useAuthStore } from "../../stores/auth.store";
import { usePaymentsStore } from "../../stores/payments.store";
import { storeToRefs } from "pinia";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const authStore = useAuthStore();
const confirm = useConfirm();

// PrimeVue
import Tabs from "primevue/tabs";
import TabList from "primevue/tablist";
import Tab from "primevue/tab";
import TabPanels from "primevue/tabpanels";
import TabPanel from "primevue/tabpanel";
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import InputNumber from "primevue/inputnumber";
import Select from "primevue/select";
import Textarea from "primevue/textarea";
import Dialog from "primevue/dialog";
import Calendar from "primevue/calendar";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import ConfirmDialog from "primevue/confirmdialog";

// Custom Components
import CustomerUsersList from "./CustomerUsersList.vue";

// Props derived from route for reusability if used as component,
// but here we use route params mainly.
const props = defineProps<{
  customerIdProp?: string;
  forceViewMode?: boolean;
}>();

const route = useRoute();
const router = useRouter();
const toast = useToast();
const customersStore = useCustomersStore();
const paymentsStore = usePaymentsStore();
const { currentCustomer, loading } = storeToRefs(customersStore);
const { payments, loading: loadingPayments } = storeToRefs(paymentsStore);

// Modes: CREATE, VIEW, EDIT
const mode = computed(() => {
  if (props.forceViewMode) return "VIEW";
  if (route.path.includes("/new")) return "CREATE";
  if (route.path.includes("/edit")) return "EDIT";
  return "VIEW";
});

const isView = computed(() => mode.value === "VIEW");
const isEdit = computed(() => mode.value === "EDIT");
const isCreate = computed(() => mode.value === "CREATE");
const isEditable = computed(() => !isView.value);

const id = computed(() => props.customerIdProp || (route.params.id as string));

// Form Data
const form = ref({
  name: "",
  contact_name: "",
  email: "",
  phone: "",
  status: "active",
  vat_number: "",
  address: "",
  city: "",
  state: "",
  zip_code: "",
  country: "",
  language: "ca",
  plan: "standard",
  billing_cycle: "monthly",
  price: 0,
  trial_ends_at: null as Date | null,
  internal_notes: "",
  max_operators: 50,
  max_workcenters: 10,
  max_shop_floors: 2,
  max_users: 5,
  max_jobs: 1000,
});

const statuses = computed(() => [
  { label: t("common.active"), value: "active" },
  { label: t("common.inactive"), value: "inactive" },
  { label: "Suspès", value: "suspended" }, // Need key for suspended if used
]);

const languages = ref([
  { label: "Català", value: "ca" },
  { label: "Español", value: "es" },
  { label: "English", value: "en" },
  { label: "Français", value: "fr" },
  { label: "Deutsch", value: "de" },
]);

const plans = computed(() => [
  { label: "Bàsic", value: "basic" },
  { label: "Estàndard", value: "standard" },
  { label: "Premium", value: "premium" },
  { label: "Enterprise", value: "enterprise" },
]);

const cycles = computed(() => [
  { label: "Mensual", value: "monthly" },
  { label: "Anual", value: "yearly" },
]);

onMounted(async () => {
  if (!isCreate.value && id.value) {
    await customersStore.fetchCustomer(id.value);
    if (currentCustomer.value) {
      // Populate form
      Object.assign(form.value, currentCustomer.value);
      // Ensure specific fields
      if (currentCustomer.value.trial_ends_at) {
        form.value.trial_ends_at = new Date(
          currentCustomer.value.trial_ends_at
        );
      }
    }
    // Fetch payments if admin
    if (authStore.isAdmin) {
      paymentsStore.fetchByCustomer(id.value);
    }
  }
});

const goBack = () => {
  router.push("/customers");
};

const goToEdit = () => {
  router.push(`/customers/${id.value}/edit`);
};

const cancelEdit = () => {
  if (isCreate.value) {
    goBack();
  } else {
    router.push(`/customers/${id.value}`);
  }
};

const save = async () => {
  try {
    const payload = {
      ...form.value,
      trial_ends_at: form.value.trial_ends_at
        ? form.value.trial_ends_at.toISOString()
        : null,
    };

    if (isCreate.value) {
      await customersStore.createCustomer(payload as any);
      toast.add({
        severity: "success",
        summary: t("customers.creation"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      goBack();
    } else {
      await customersStore.updateCustomer(id.value, payload as any);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      router.push(`/customers/${id.value}`); // Go back to view mode
    }
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: t("common.error"),
      detail: e.message || t("common.error"),
      life: 3000,
    });
  }
};

// Payment Logic
const showPaymentDialog = ref(false);
const paymentForm = ref({
  date: new Date(),
  amount: 0,
});

const openPaymentDialog = () => {
  paymentForm.value.date = new Date();
  paymentForm.value.amount = form.value.price || 0;
  showPaymentDialog.value = true;
};

const savePayment = async () => {
  try {
    await paymentsStore.createPayment({
      customer_id: id.value,
      amount: paymentForm.value.amount,
      paid_at: paymentForm.value.date.toISOString(),
      status: "paid",
      currency: "EUR",
    });
    toast.add({
      severity: "success",
      summary: t("common.success"),
      detail: t("customers.save_success"), // Reusing generic success or add payment_success
      life: 2000,
    });
    showPaymentDialog.value = false;
    paymentsStore.fetchByCustomer(id.value);
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: e.message,
      life: 3000,
    });
  }
};

const confirmDeletePayment = (paymentId: string) => {
  confirm.require({
    message: t("customers.msg_delete_payment"),
    header: t("common.confirm_title"),
    icon: "pi pi-exclamation-triangle",
    accept: () => {
      paymentsStore
        .deletePayment(paymentId)
        .then(() => {
          toast.add({
            severity: "success",
            summary: t("common.success"),
            detail: "Payment deleted",
            life: 2000,
          });
        })
        .catch((e: any) => {
          toast.add({
            severity: "error",
            summary: "Error",
            detail: e.message,
            life: 3000,
          });
        });
    },
  });
};
</script>

<template>
  <div
    class="card bg-white rounded-lg shadow-sm border border-surface-200 flex flex-col h-full"
  >
    <!-- Header -->
    <div
      class="p-6 border-b border-surface-100 flex justify-between items-center bg-surface-50 rounded-t-lg"
    >
      <div class="flex items-center gap-4">
        <Button
          icon="pi pi-arrow-left"
          text
          rounded
          severity="secondary"
          @click="goBack"
        />
        <div>
          <h1 class="text-xl font-bold text-surface-900">
            {{
              isCreate
                ? t("customers.new_customer")
                : form.name || t("customers.customer_detail")
            }}
          </h1>
          <div class="text-sm text-surface-500 flex gap-2 items-center">
            <span v-if="isView">{{ t("customers.mode_view") }}</span>
            <span v-else-if="isEdit">{{ t("customers.mode_edit") }}</span>
            <span v-else>{{ t("customers.creation") }}</span>
          </div>
        </div>
      </div>

      <div class="flex gap-2">
        <Button
          v-if="isView && authStore.isAdmin"
          :label="t('common.edit')"
          icon="pi pi-pencil"
          severity="secondary"
          @click="goToEdit"
        />
        <Button
          v-if="isEditable"
          :label="t('common.cancel')"
          severity="secondary"
          text
          @click="cancelEdit"
        />
        <Button
          v-if="isEditable"
          :label="t('common.save')"
          icon="pi pi-save"
          @click="save"
          :loading="loading"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden flex flex-col">
      <Tabs value="data" class="flex-1 flex flex-col">
        <TabList>
          <Tab value="data">{{ t("customers.tabs.data") }}</Tab>
          <Tab value="users" :disabled="isCreate">{{
            t("customers.tabs.users")
          }}</Tab>
          <Tab value="payments" v-if="authStore.isAdmin" :disabled="isCreate">{{
            t("customers.tabs.payments")
          }}</Tab>
        </TabList>

        <TabPanels class="flex-1 overflow-auto p-6 bg-surface-50/50">
          <!-- DATA TAB -->
          <TabPanel value="data">
            <div
              class="grid grid-cols-1 lg:grid-cols-2 gap-8 max-w-6xl mx-auto"
            >
              <!-- BLOCK 1: GENERAL INFO -->
              <div
                class="bg-white p-6 rounded-lg shadow-sm border border-surface-200"
              >
                <h2
                  class="text-lg font-bold text-surface-800 mb-4 pb-2 border-b border-surface-100 flex items-center gap-2"
                >
                  <i class="pi pi-id-card text-surface-400"></i>
                  {{ t("customers.general_info") }}
                </h2>
                <div class="grid grid-cols-1 gap-4">
                  <div class="flex flex-col gap-2">
                    <label
                      for="name"
                      class="font-medium text-sm text-surface-600"
                      >{{ t("customers.fiscal_name") }}</label
                    >
                    <InputText
                      id="name"
                      v-model="form.name"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2">
                      <label
                        for="contact"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.contact_name") }}</label
                      >
                      <InputText
                        id="contact"
                        v-model="form.contact_name"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label
                        for="status"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.status") }}</label
                      >
                      <Select
                        id="status"
                        v-model="form.status"
                        :options="statuses"
                        optionLabel="label"
                        optionValue="value"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                  </div>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2">
                      <label
                        for="email"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.email") }}</label
                      >
                      <InputText
                        id="email"
                        v-model="form.email"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label
                        for="phone"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.phone") }}</label
                      >
                      <InputText
                        id="phone"
                        v-model="form.phone"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- BLOCK 2: BILLING & ADDRESS -->
              <div
                class="bg-white p-6 rounded-lg shadow-sm border border-surface-200"
              >
                <h2
                  class="text-lg font-bold text-surface-800 mb-4 pb-2 border-b border-surface-100 flex items-center gap-2"
                >
                  <i class="pi pi-receipt text-surface-400"></i>
                  {{ t("customers.billing_address") }}
                </h2>
                <div class="grid grid-cols-1 gap-4">
                  <div class="grid grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2">
                      <label
                        for="vat"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.vat") }}</label
                      >
                      <InputText
                        id="vat"
                        v-model="form.vat_number"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label
                        for="lang"
                        class="font-medium text-sm text-surface-600"
                        >{{ t("customers.language") }}</label
                      >
                      <Select
                        id="lang"
                        v-model="form.language"
                        :options="languages"
                        optionLabel="label"
                        optionValue="value"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                  </div>
                  <div class="flex flex-col gap-2">
                    <label
                      for="addr"
                      class="font-medium text-sm text-surface-600"
                      >Adreça</label
                    >
                    <InputText
                      id="addr"
                      v-model="form.address"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="grid grid-cols-3 gap-4">
                    <div class="flex flex-col gap-2">
                      <label
                        for="city"
                        class="font-medium text-sm text-surface-600"
                        >Ciutat</label
                      >
                      <InputText
                        id="city"
                        v-model="form.city"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label
                        for="zip"
                        class="font-medium text-sm text-surface-600"
                        >CP</label
                      >
                      <InputText
                        id="zip"
                        v-model="form.zip_code"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label
                        for="country"
                        class="font-medium text-sm text-surface-600"
                        >País</label
                      >
                      <InputText
                        id="country"
                        v-model="form.country"
                        :disabled="!isEditable"
                        class="w-full"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- BLOCK 3: SUBSCRIPTION (ADMIN ONLY EDIT) -->
              <div
                class="bg-white p-6 rounded-lg shadow-sm border border-surface-200"
              >
                <h2
                  class="text-lg font-bold text-surface-800 mb-4 pb-2 border-b border-surface-100 flex items-center gap-2"
                >
                  <i class="pi pi-credit-card text-surface-400"></i> Subscripció
                </h2>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="flex flex-col gap-2">
                    <label class="font-medium text-sm text-surface-600"
                      >Pla</label
                    >
                    <Select
                      v-model="form.plan"
                      :options="plans"
                      optionLabel="label"
                      optionValue="value"
                      :disabled="!isEditable || !authStore.isAdmin"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label class="font-medium text-sm text-surface-600"
                      >Cicle de Facturació</label
                    >
                    <Select
                      v-model="form.billing_cycle"
                      :options="cycles"
                      optionLabel="label"
                      optionValue="value"
                      :disabled="!isEditable || !authStore.isAdmin"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label class="font-medium text-sm text-surface-600"
                      >Preu</label
                    >
                    <InputNumber
                      v-model="form.price"
                      mode="currency"
                      currency="EUR"
                      locale="es-ES"
                      :disabled="!isEditable || !authStore.isAdmin"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label class="font-medium text-sm text-surface-600"
                      >Fi de Prova</label
                    >
                    <!-- Use Text or Calendar depending on type. Form uses string or null. Calendar handles Date or String usually if matched. 
                                 I'll just use a simple date input or PrimeVue Calendar. Calendar is safer. -->
                    <!-- Need to cast if form.trial_ends_at is string -->
                    <Calendar
                      v-model="form.trial_ends_at"
                      dateFormat="yy-mm-dd"
                      :disabled="!isEditable || !authStore.isAdmin"
                      class="w-full"
                      showIcon
                    />
                  </div>
                </div>
              </div>

              <!-- BLOCK 4: CONFIGURATION (LIMITS) -->
              <div
                class="bg-white p-6 rounded-lg shadow-sm border border-surface-200 lg:col-span-2"
              >
                <!-- ... -->

                <h2
                  class="text-lg font-bold text-surface-800 mb-4 pb-2 border-b border-surface-100 flex items-center gap-2"
                >
                  <i class="pi pi-cog text-surface-400"></i>
                  {{ t("customers.config_limits") }}
                </h2>
                <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
                  <div class="flex flex-col gap-2">
                    <label
                      for="max_op"
                      class="font-medium text-sm text-surface-600"
                      >{{ t("customers.max_operators") }}</label
                    >
                    <InputNumber
                      id="max_op"
                      v-model="form.max_operators"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label
                      for="max_wc"
                      class="font-medium text-sm text-surface-600"
                      >Max. Centres</label
                    >
                    <InputNumber
                      id="max_wc"
                      v-model="form.max_workcenters"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label
                      for="max_sf"
                      class="font-medium text-sm text-surface-600"
                      >Max. Plantes</label
                    >
                    <InputNumber
                      id="max_sf"
                      v-model="form.max_shop_floors"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label
                      for="max_u"
                      class="font-medium text-sm text-surface-600"
                      >Max. Usuaris</label
                    >
                    <InputNumber
                      id="max_u"
                      v-model="form.max_users"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-2">
                    <label
                      for="max_j"
                      class="font-medium text-sm text-surface-600"
                      >Max. Feines</label
                    >
                    <InputNumber
                      id="max_j"
                      v-model="form.max_jobs"
                      :disabled="!isEditable"
                      class="w-full"
                    />
                  </div>
                </div>
                <div class="mt-4 flex flex-col gap-2">
                  <label
                    for="notes"
                    class="font-medium text-sm text-surface-600"
                    >Notes Internes</label
                  >
                  <Textarea
                    id="notes"
                    v-model="form.internal_notes"
                    :disabled="!isEditable"
                    rows="3"
                    class="w-full"
                  />
                </div>
              </div>
            </div>
          </TabPanel>

          <!-- USERS TAB -->
          <TabPanel value="users">
            <CustomerUsersList
              v-if="currentCustomer && currentCustomer.id"
              :customerId="currentCustomer.id"
              :readOnly="!authStore.isAdmin"
            />
            <div v-else class="text-center p-6 text-surface-500">
              {{ t("customers.save_users_first") }}
            </div>
          </TabPanel>

          <!-- PAYMENTS TAB -->
          <TabPanel value="payments">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-bold text-surface-800">
                {{ t("customers.payments_history") }}
              </h3>
              <Button
                :label="t('customers.add_payment')"
                icon="pi pi-plus"
                class="p-button-sm"
                @click="openPaymentDialog"
                :disabled="!currentCustomer?.id"
              />
            </div>

            <DataTable
              :value="payments"
              :loading="loadingPayments"
              stripedRows
              responsiveLayout="scroll"
            >
              <Column :header="t('customers.payment_date')">
                <template #body="slotProps">
                  {{ new Date(slotProps.data.paid_at).toLocaleDateString() }}
                </template>
              </Column>
              <Column :header="t('customers.import')">
                <template #body="slotProps">
                  {{ slotProps.data.amount }} {{ slotProps.data.currency }}
                </template>
              </Column>
              <Column field="status" :header="t('common.status')">
                <template #body="slotProps">
                  <span
                    class="capitalize"
                    :class="
                      slotProps.data.status === 'paid'
                        ? 'text-green-600 font-bold'
                        : ''
                    "
                  >
                    {{ slotProps.data.status }}
                  </span>
                </template>
              </Column>
              <Column :header="t('common.actions')" style="width: 100px">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    @click="confirmDeletePayment(slotProps.data.id)"
                  />
                </template>
              </Column>
              <template #empty>{{ t("customers.no_payments") }}</template>
            </DataTable>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </div>

    <!-- Payment Dialog -->
    <Dialog
      v-model:visible="showPaymentDialog"
      :header="t('customers.add_payment')"
      :modal="true"
      class="w-full max-w-sm"
    >
      <div class="flex flex-col gap-4 pt-2">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm">{{
            t("customers.payment_date")
          }}</label>
          <Calendar v-model="paymentForm.date" showIcon dateFormat="dd/mm/yy" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-sm"
            >{{ t("customers.import") }} (EUR)</label
          >
          <InputNumber
            v-model="paymentForm.amount"
            mode="currency"
            currency="EUR"
            locale="es-ES"
          />
        </div>
        <div class="flex justify-end gap-2 mt-2">
          <Button
            :label="t('common.cancel')"
            severity="secondary"
            @click="showPaymentDialog = false"
            text
          />
          <Button
            :label="t('common.save')"
            @click="savePayment"
            :loading="loadingPayments"
          />
        </div>
      </div>
    </Dialog>

    <ConfirmDialog></ConfirmDialog>
  </div>
</template>
