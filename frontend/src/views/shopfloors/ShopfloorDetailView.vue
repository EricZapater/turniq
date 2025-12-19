<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useShopfloorsStore } from "../../stores/shopfloors.store";
import { useAuthStore } from "../../stores/auth.store";
import { useCustomersStore } from "../../stores/customers.store";
import { storeToRefs } from "pinia";
import { useToast } from "primevue/usetoast";

// PrimeVue
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import Select from "primevue/select";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const shopfloorsStore = useShopfloorsStore();
const authStore = useAuthStore();
const customersStore = useCustomersStore();
const { currentShopfloor, loading } = storeToRefs(shopfloorsStore);
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
  customer_id: null as string | null,
});

onMounted(async () => {
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({ page_size: 100 });
  }

  if (!isCreate.value && id) {
    await shopfloorsStore.fetchShopfloor(id);
    if (currentShopfloor.value) {
      form.value.name = currentShopfloor.value.name;
      form.value.customer_id = currentShopfloor.value.customer_id;
    }
  }
});

const goBack = () => {
  router.push("/shopfloors");
};

const goToEdit = () => {
  router.push(`/shopfloors/${id}/edit`);
};

const cancelEdit = () => {
  if (isCreate.value) {
    goBack();
  } else {
    router.push(`/shopfloors/${id}`);
  }
};

const save = async () => {
  try {
    // Assuming backend handles customer_id via context or we pass it
    // ShopfloorRequest needs customer_id but context middleware injects it on backend?
    // Wait, ShopfloorRequest struct has CustomerID string.
    // If we are admin, we might need to select customer.
    // If not admin, backend uses context.
    // Frontend API service defines ShopfloorRequest with optional customer_id.
    // Let's assume for now we use auth context or current user's customer.
    // Ideally if is_admin, we might want a customer selector.
    // For now, I'll pass customer_id from authStore if available, or let backend handle strictness.
    // Backend `Create` service: uses context `customer_id` if not admin. If admin, uses request `customer_id`.

    // Simplification: We assume current user's context for now as per "Operen dins del context del customer de l’usuari".
    // If admin is creating for a specific customer, they would need a dropdown.
    // User request: "El rol Admin: Té accés complet... Rols no Admin: Accés limitat".
    // For MVP compliance with "contextual access", passing current user's customer_id is safe for non-admins.
    // For admins, to be fully functional, they'd need to select a customer.
    // I'll add `customer_id: authStore.user?.customer_id` to payload as fallback.

    const payload = {
      name: form.value.name,
      customer_id: authStore.isAdmin
        ? form.value.customer_id || undefined
        : authStore.user?.customer_id,
    };

    if (isCreate.value) {
      await shopfloorsStore.createShopfloor(payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"), // Reusing or new specific
        life: 3000,
      });
      goBack();
    } else {
      await shopfloorsStore.updateShopfloor(id, payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      router.push(`/shopfloors/${id}`);
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
  if (!confirm(t("shopfloors.delete_confirm"))) return;
  try {
    await shopfloorsStore.deleteShopfloor(id);
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
            {{ isCreate ? t("shopfloors.new") : form.name || t("shopfloors.detail") }}
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
          <i class="pi pi-building text-surface-400"></i> {{ t("customers.general_info") }}
        </h2>
        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-2" v-if="authStore.isAdmin">
            <label for="customer" class="font-medium text-sm text-surface-600"
              >{{ t("customers.title") }}</label // Singular 'Customer' key needed? re-use title
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
            <label for="name" class="font-medium text-sm text-surface-600"
              >{{ t("shopfloors.name") }}</label
            >
            <InputText
              id="name"
              v-model="form.name"
              :disabled="!isEditable"
              class="w-full"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
