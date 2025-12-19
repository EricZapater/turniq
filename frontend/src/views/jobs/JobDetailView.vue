<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useJobsStore } from "../../stores/jobs.store";
import { useShopfloorsStore } from "../../stores/shopfloors.store";
import { useWorkcentersStore } from "../../stores/workcenters.store";
import { useAuthStore } from "../../stores/auth.store";
import { useCustomersStore } from "../../stores/customers.store";
import { storeToRefs } from "pinia";
import { useToast } from "primevue/usetoast";

// PrimeVue
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import InputNumber from "primevue/inputnumber";
import Select from "primevue/select";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const jobsStore = useJobsStore();
const shopfloorsStore = useShopfloorsStore();
const workcentersStore = useWorkcentersStore();
const authStore = useAuthStore();
const customersStore = useCustomersStore();

const { currentJob, loading } = storeToRefs(jobsStore);
const { shopfloors } = storeToRefs(shopfloorsStore);
const { workcenters } = storeToRefs(workcentersStore);
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
  job_code: "",
  product_code: "",
  description: "",
  estimated_duration: 0,
  shop_floor_id: null as string | null,
  workcenter_id: null as string | null,
  customer_id: null as string | null,
});

onMounted(async () => {
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({ page_size: 100 });
  }
  // Load Shopfloors for dropdown
  await shopfloorsStore.fetchShopfloors({ page_size: 100 });
  // Load Workcenters for dropdown (ideally filtered by shopfloor, but for now fetch all to populate)
  await workcentersStore.fetchWorkcenters({ page_size: 100 });

  if (!isCreate.value && id) {
    await jobsStore.fetchJob(id);
    if (currentJob.value) {
      form.value.job_code = currentJob.value.job_code;
      form.value.product_code = currentJob.value.product_code;
      form.value.description = currentJob.value.description;
      form.value.estimated_duration = currentJob.value.estimated_duration;
      form.value.shop_floor_id = currentJob.value.shop_floor_id;
      form.value.workcenter_id = currentJob.value.workcenter_id;
      form.value.customer_id = currentJob.value.customer_id;
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
          form.value.workcenter_id = null; // Also clear workcenter derived from shopfloor
        }
      }
    }
  }
);

// Computed filtered workcenters based on selected shopfloor
const filteredWorkcenters = computed(() => {
  if (!form.value.shop_floor_id) return [];
  return workcenters.value.filter(
    (wc) => wc.shop_floor_id === form.value.shop_floor_id
  );
});

// Watch shopfloor change to clear workcenter if invalid
watch(
  () => form.value.shop_floor_id,
  (newVal) => {
    if (newVal && form.value.workcenter_id) {
      // Check if current workcenter belongs to new shopfloor
      const wc = workcenters.value.find(
        (w) => w.id === form.value.workcenter_id
      );
      if (wc && wc.shop_floor_id !== newVal) {
        form.value.workcenter_id = null;
      }
    } else if (!newVal) {
      form.value.workcenter_id = null;
    }
  }
);

const goBack = () => {
  router.push("/jobs");
};

const goToEdit = () => {
  router.push(`/jobs/${id}/edit`);
};

const cancelEdit = () => {
  if (isCreate.value) {
    goBack();
  } else {
    router.push(`/jobs/${id}`);
  }
};

const save = async () => {
  try {
    if (!form.value.shop_floor_id || !form.value.workcenter_id) {
      toast.add({
        severity: "warn",
        summary: t("shifts.missing_data"),
        detail: t("shifts.missing_data"),
        life: 3000,
      });
      return;
    }

    const payload = {
      job_code: form.value.job_code,
      product_code: form.value.product_code,
      description: form.value.description,
      estimated_duration: form.value.estimated_duration,
      shop_floor_id: form.value.shop_floor_id,
      workcenter_id: form.value.workcenter_id,
      customer_id: authStore.isAdmin
        ? form.value.customer_id || undefined
        : authStore.user?.customer_id,
    };

    if (isCreate.value) {
      await jobsStore.createJob(payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      goBack();
    } else {
      await jobsStore.updateJob(id, payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("customers.save_success"),
        life: 3000,
      });
      router.push(`/jobs/${id}`);
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
  if (!confirm(t("jobs.delete_confirm"))) return;
  try {
    await jobsStore.deleteJob(id);
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
            {{ isCreate ? t("jobs.new") : form.job_code || t("jobs.detail") }}
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
          <i class="pi pi-briefcase text-surface-400"></i> {{ t("jobs.info") }}
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            class="flex flex-col gap-2 md:col-span-2"
            v-if="authStore.isAdmin"
          >
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
            <label for="jcode" class="font-medium text-sm text-surface-600">{{
              t("jobs.job_code")
            }}</label>
            <InputText
              id="jcode"
              v-model="form.job_code"
              :disabled="!isEditable"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label for="pcode" class="font-medium text-sm text-surface-600">{{
              t("jobs.product_code")
            }}</label>
            <InputText
              id="pcode"
              v-model="form.product_code"
              :disabled="!isEditable"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-2 md:col-span-2">
            <label for="desc" class="font-medium text-sm text-surface-600">{{
              t("jobs.description")
            }}</label>
            <Textarea
              id="desc"
              v-model="form.description"
              :disabled="!isEditable"
              rows="3"
              class="w-full"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label for="dur" class="font-medium text-sm text-surface-600">{{
              t("jobs.duration")
            }}</label>
            <InputNumber
              id="dur"
              v-model="form.estimated_duration"
              :disabled="!isEditable"
              class="w-full"
            />
          </div>
        </div>

        <h2
          class="text-lg font-bold text-surface-800 mb-4 mt-6 pb-2 border-b border-surface-100 flex items-center gap-2"
        >
          <i class="pi pi-cog text-surface-400"></i> {{ t("jobs.assignment") }}
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
          <div class="flex flex-col gap-2">
            <label
              for="workcenter"
              class="font-medium text-sm text-surface-600"
              >{{ t("jobs.workcenter") }}</label
            >
            <Select
              id="workcenter"
              v-model="form.workcenter_id"
              :options="filteredWorkcenters"
              optionLabel="name"
              optionValue="id"
              placeholder="Selecciona un centre"
              :disabled="!isEditable || !form.shop_floor_id"
              class="w-full"
              filter
            />
            <small v-if="!form.shop_floor_id" class="text-surface-400">{{
              t("jobs.select_shopfloor_first")
            }}</small>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
