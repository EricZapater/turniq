<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useShiftsStore } from "../../stores/shifts.store";
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
import ColorPicker from "primevue/colorpicker";
import DatePicker from "primevue/datepicker"; // For time
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const shiftsStore = useShiftsStore();
const shopfloorsStore = useShopfloorsStore();
const authStore = useAuthStore();
const customersStore = useCustomersStore();

const { currentShift, loading } = storeToRefs(shiftsStore);
const { shopfloors } = storeToRefs(shopfloorsStore);
const { customers } = storeToRefs(customersStore);

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
  color: "3B82F6", // Default Blue
  start_time: null as Date | null,
  end_time: null as Date | null,
  customer_id: null as string | null,
  shopfloor_id: null as string | null,
  is_active: true,
});

onMounted(async () => {
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({ page_size: 100 });
  }
  await shopfloorsStore.fetchShopfloors({ page_size: 100 });

  if (!isCreate.value && id) {
    await shiftsStore.fetchShift(id);
    if (currentShift.value) {
      form.value.name = currentShift.value.name;
      form.value.color = currentShift.value.color.replace("#", ""); // ColorPicker handles hex without # usually or implies it
      form.value.customer_id = currentShift.value.customer_id;
      form.value.shopfloor_id = currentShift.value.shopfloor_id || null;
      form.value.is_active = currentShift.value.is_active;

      // Parse time strings 'HH:mm' or ISO to Date objects for DatePicker
      if (currentShift.value.start_time) {
        const timeStr = currentShift.value.start_time;
        // Backend returns 'YYYY-MM-DDTHH:mm:ssZ' or 'HH:mm'
        // We want to extract 'HH:mm' exactly as is, shielding from timezone shifts/year 0 issues
        let hours, minutes;
        if (timeStr.includes("T")) {
          const timePart = timeStr.split("T")[1]; // "06:00:00Z"
          const hhmm = timePart.substring(0, 5); // "06:00"
          [hours, minutes] = hhmm.split(":");
        } else {
          [hours, minutes] = timeStr.split(":");
        }

        const d = new Date();
        d.setHours(parseInt(hours), parseInt(minutes), 0, 0);
        form.value.start_time = d;
      }
      if (currentShift.value.end_time) {
        const timeStr = currentShift.value.end_time;
        let hours, minutes;
        if (timeStr.includes("T")) {
          const timePart = timeStr.split("T")[1];
          const hhmm = timePart.substring(0, 5);
          [hours, minutes] = hhmm.split(":");
        } else {
          [hours, minutes] = timeStr.split(":");
        }

        const d = new Date();
        d.setHours(parseInt(hours), parseInt(minutes), 0, 0);
        form.value.end_time = d;
      }
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
      if (form.value.shopfloor_id) {
        const sf = shopfloors.value.find(
          (s) => s.id === form.value.shopfloor_id
        );
        if (sf && sf.customer_id !== newVal) {
          form.value.shopfloor_id = null;
        }
      }
    }
  }
);

const goBack = () => {
  router.push("/shifts");
};

const goToEdit = () => {
  router.push(`/shifts/${id}/edit`);
};

const cancelEdit = () => {
  if (isCreate.value) {
    goBack();
  } else {
    router.push(`/shifts/${id}`);
  }
};

// Helper: Format Date to HH:mm
const formatTimeHHMM = (d: Date | null) => {
  if (!d) return "00:00";
  return d.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
};

const save = async () => {
  try {
    if (!form.value.shopfloor_id) {
      toast.add({
        severity: "warn",
        summary: t("shifts.missing_data"),
        detail: t("shifts.missing_data"), // Or specific "select shopfloor"
        life: 3000,
      });
      return;
    }
    if (!form.value.start_time || !form.value.end_time) {
      toast.add({
        severity: "warn",
        summary: t("shifts.missing_data"),
        detail: t("shifts.select_dates"),
        life: 3000,
      });
      return;
    }

    // Backend expects HH:mm string maybe? Or ISO? Service uses time.Parse("15:04", ...)
    const payload = {
      name: form.value.name,
      color: "#" + form.value.color, // Add hash back
      start_time: formatTimeHHMM(form.value.start_time),
      end_time: formatTimeHHMM(form.value.end_time),
      shopfloor_id: form.value.shopfloor_id,
      is_active: form.value.is_active,
      customer_id: authStore.isAdmin
        ? form.value.customer_id || "" // Backend requires string
        : authStore.user?.customer_id || "",
    };

    // Ensure customer_id is set if admin forgot or something (validation should catch it, but safe fallback)
    if (authStore.isAdmin && !payload.customer_id) {
      // Try to get from shopfloor?
      const sf = shopfloors.value.find((s) => s.id === form.value.shopfloor_id);
      if (sf) payload.customer_id = sf.customer_id;
    }

    if (isCreate.value) {
      await shiftsStore.createShift(payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      goBack();
    } else {
      await shiftsStore.updateShift(id, payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      router.push(`/shifts/${id}`);
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
  if (!confirm(t("shifts.delete_confirm"))) return;
  try {
    await shiftsStore.deleteShift(id);
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
            {{ isCreate ? t("shifts.new") : form.name || t("shifts.detail") }}
          </h1>
          <div class="text-sm text-surface-500 flex gap-2 items-center">
            <div
              v-if="!isCreate"
              :style="{ backgroundColor: '#' + form.color }"
              class="w-3 h-3 rounded-full border border-surface-200"
            ></div>
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
          <i class="pi pi-clock text-surface-400"></i>
          {{ t("customers.general_info") }}
        </h2>
        <div class="flex flex-col gap-4">
          <!-- Customer Selector (Admin) -->
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

          <!-- Shopfloor Selector -->
          <div class="flex flex-col gap-2">
            <label
              for="shopfloor"
              class="font-medium text-sm text-surface-600"
              >{{ t("shifts.shopfloor") }}</label
            >
            <Select
              id="shopfloor"
              v-model="form.shopfloor_id"
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

          <!-- Name & Color -->
          <div class="flex gap-4">
            <div class="flex flex-col gap-2 flex-grow">
              <label for="name" class="font-medium text-sm text-surface-600">{{
                t("shifts.name")
              }}</label>
              <InputText
                id="name"
                v-model="form.name"
                :disabled="!isEditable"
                class="w-full"
                placeholder="Ex: Matí, Tarda..."
              />
            </div>
            <div class="flex flex-col gap-2">
              <label for="color" class="font-medium text-sm text-surface-600">{{
                t("shifts.color")
              }}</label>
              <ColorPicker v-model="form.color" :disabled="!isEditable" />
            </div>
          </div>

          <!-- Times -->
          <div class="grid grid-cols-2 gap-4">
            <div class="flex flex-col gap-2">
              <label
                for="start_time"
                class="font-medium text-sm text-surface-600"
                >{{ t("shifts.start_time") }}</label
              >
              <DatePicker
                id="start_time"
                v-model="form.start_time"
                timeOnly
                fluid
                :disabled="!isEditable"
                hourFormat="24"
                showIcon
                iconDisplay="input"
              />
            </div>
            <div class="flex flex-col gap-2">
              <label
                for="end_time"
                class="font-medium text-sm text-surface-600"
                >{{ t("shifts.end_time") }}</label
              >
              <DatePicker
                id="end_time"
                v-model="form.end_time"
                timeOnly
                fluid
                :disabled="!isEditable"
                hourFormat="24"
                showIcon
                iconDisplay="input"
              />
            </div>
          </div>

          <div class="flex items-center gap-2 mt-2">
            <ToggleSwitch v-model="form.is_active" :disabled="!isEditable" />
            <label class="font-medium text-sm text-surface-600">{{
              t("common.active")
            }}</label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
