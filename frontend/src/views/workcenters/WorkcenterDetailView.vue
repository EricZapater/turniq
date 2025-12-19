<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useWorkcentersStore } from "../../stores/workcenters.store";
import { useShopfloorsStore } from "../../stores/shopfloors.store";
import { useAuthStore } from "../../stores/auth.store";
import { useCustomersStore } from "../../stores/customers.store";
import { storeToRefs } from "pinia";
import { useToast } from "primevue/usetoast";

// PrimeVue
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import Select from "primevue/select";
import ToggleSwitch from "primevue/toggleswitch";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const workcentersStore = useWorkcentersStore();
const shopfloorsStore = useShopfloorsStore();
const authStore = useAuthStore();
const customersStore = useCustomersStore();

const { currentWorkcenter, loading } = storeToRefs(workcentersStore);
const { shopfloors } = storeToRefs(shopfloorsStore);
const { customers } = storeToRefs(customersStore);

// Modes
const mode = computed(() => {
  if (route.path.includes("/new")) return "CREATE";
  if (route.path.includes("/edit")) return "EDIT";
  return "VIEW";
});

const isView = computed(() => mode.value === "VIEW");
const isEdit = computed(() => mode.value === "EDIT");
const isCreate = computed(() => mode.value === "CREATE");
const isEditable = computed(() => !isView.value);

const id = route.params.id as string;

const form = ref({
  name: "",
  shop_floor_id: null as string | null,
  customer_id: null as string | null,
  is_active: true,
});

onMounted(async () => {
  // Load common data
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({ page_size: 100 });
  }
  // Load Shopfloors for dropdown
  await shopfloorsStore.fetchShopfloors({ page_size: 100 }); // Fetch all (limited to 100 for now)

  if (!isCreate.value && id) {
    await workcentersStore.fetchWorkcenter(id);
    if (currentWorkcenter.value) {
      form.value.name = currentWorkcenter.value.name;
      form.value.shop_floor_id = currentWorkcenter.value.shop_floor_id;
      form.value.customer_id = currentWorkcenter.value.customer_id;
      form.value.is_active = currentWorkcenter.value.is_active;
    }
  }
});

// Computed filtered shopfloors based on selected customer (for admins)
const filteredShopfloors = computed(() => {
  if (authStore.isAdmin && form.value.customer_id) {
    return shopfloors.value.filter(
      (s) => s.customer_id === form.value.customer_id
    );
  }
  return shopfloors.value;
});

// Watch customer change to clear shopfloor if invalid
watch(
  () => form.value.customer_id,
  (newVal) => {
    if (authStore.isAdmin) {
      // Check if current shopfloor belongs to new customer
      if (form.value.shop_floor_id) {
        const sf = shopfloors.value.find(
          (s) => s.id === form.value.shop_floor_id
        );
        if (sf && sf.customer_id !== newVal) {
          form.value.shop_floor_id = null;
        }
      }
    }
  }
);

const goBack = () => {
  router.push("/workcenters");
};

const goToEdit = () => {
  router.push(`/workcenters/${id}/edit`);
};

const cancelEdit = () => {
  if (isCreate.value) {
    goBack();
  } else {
    router.push(`/workcenters/${id}`);
  }
};

const save = async () => {
  try {
    const payload = {
      name: form.value.name,
      shop_floor_id: form.value.shop_floor_id,
      is_active: form.value.is_active,
      customer_id: authStore.isAdmin
        ? form.value.customer_id || undefined
        : authStore.user?.customer_id,
    };

    if (isCreate.value) {
      await workcentersStore.createWorkcenter(payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      goBack();
    } else {
      await workcentersStore.updateWorkcenter(id, payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      router.push(`/workcenters/${id}`);
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

const confirmDelete = async () => {
  if (!confirm(t("workcenters.delete_confirm"))) return;
  try {
    await workcentersStore.deleteWorkcenter(id);
    toast.add({
      severity: "success",
      summary: t("common.success"),
      detail: "Deleted",
      life: 3000,
    });
    goBack();
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: t("common.error"),
      detail: e.message || t("common.error"),
      life: 3000,
    });
  }
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
                ? t("workcenters.new")
                : form.name || t("workcenters.detail")
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
          v-if="isView"
          :label="t('common.delete')"
          icon="pi pi-trash"
          severity="danger"
          text
          @click="confirmDelete"
        />
        <Button
          v-if="isView"
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
    <div class="p-6 max-w-4xl mx-auto w-full">
      <div class="bg-white p-6 rounded-lg shadow-sm border border-surface-200">
        <h2
          class="text-lg font-bold text-surface-800 mb-4 pb-2 border-b border-surface-100 flex items-center gap-2"
        >
          <i class="pi pi-cog text-surface-400"></i>
          {{ t("customers.general_info") }}
        </h2>
        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-2" v-if="authStore.isAdmin">
            <label
              for="customer"
              class="font-medium text-sm text-surface-600"
              >{{ t("customers.title") }}</label
            >
            <Select
              id="customer"
              v-model="form.customer_id"
              :options="customers"
              optionLabel="name"
              optionValue="id"
              placeholder="Selecciona un client"
              :disabled="!isEditable"
              class="w-full"
              showClear
              filter
            />
          </div>
          <div class="flex flex-col gap-2">
            <label for="name" class="font-medium text-sm text-surface-600">{{
              t("workcenters.name")
            }}</label>
            <InputText
              id="name"
              v-model="form.name"
              :disabled="!isEditable"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label
              for="shopfloor"
              class="font-medium text-sm text-surface-600"
              >{{ t("workcenters.shopfloor") }}</label
            >
            <Select
              id="shopfloor"
              v-model="form.shop_floor_id"
              :options="filteredShopfloors"
              optionLabel="name"
              optionValue="id"
              placeholder="Selecciona una planta"
              :disabled="
                !isEditable || (authStore.isAdmin && !form.customer_id)
              "
              class="w-full"
              filter
            />
          </div>
          <div class="flex items-center gap-2 mt-2">
            <ToggleSwitch
              v-model="form.is_active"
              :disabled="!isEditable"
              inputId="is_active"
            />
            <label for="is_active" class="cursor-pointer">{{
              t("common.active")
            }}</label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
