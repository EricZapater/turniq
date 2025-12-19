<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useAuthStore } from "../stores/auth.store";
import { customersApi } from "../api/customers.api";
import Card from "primevue/card";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const authStore = useAuthStore();
const isAdmin = computed(() => authStore.isAdmin);

const customersUtils = ref({
  total: 0,
  active: 0,
  loading: false,
});

onMounted(async () => {
  if (isAdmin.value) {
    fetchStats();
  }
});

const fetchStats = async () => {
  customersUtils.value.loading = true;
  try {
    // Fetching all customers (assuming pagination might need handling if list is huge,
    // but for a dashboard summary often a specific stats endpoint is better.
    // For now, fetching first page/list is what we have)
    const response = await customersApi.list({ page: 1, page_size: 1000 });
    const customers = response.data || [];
    customersUtils.value.total = customers.length; // Approximate if paginated, ideally use response.total if available
    customersUtils.value.active = customers.filter(
      (c) => c.status === "active"
    ).length;
  } catch (e) {
    console.error("Failed to load customer stats", e);
  } finally {
    customersUtils.value.loading = false;
  }
};
</script>

<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-surface-900 mb-6">
      {{ t("dashboard.title") }}
    </h1>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- KPI: Active Customers (Admin Only) -->
      <div
        v-if="isAdmin"
        class="bg-white p-6 rounded-lg border border-surface-200 shadow-sm flex items-start justify-between"
      >
        <div>
          <span class="block text-surface-500 font-medium mb-1">{{
            t("dashboard.active_customers")
          }}</span>
          <div class="text-3xl font-bold text-surface-900">
            <span
              v-if="customersUtils.loading"
              class="pi pi-spin pi-spinner text-xl"
            ></span>
            <span v-else
              >{{ customersUtils.active }}
              <span class="text-surface-400 text-lg font-normal"
                >/ {{ customersUtils.total }}</span
              ></span
            >
          </div>
        </div>
        <div
          class="w-10 h-10 rounded-full bg-primary-50 flex items-center justify-center text-primary-600"
        >
          <i class="pi pi-briefcase text-xl"></i>
        </div>
      </div>

      <!-- Placeholder for other KPIs -->
      <div
        class="bg-white p-6 rounded-lg border border-surface-200 shadow-sm flex items-start justify-between opacity-60"
      >
        <div>
          <span class="block text-surface-500 font-medium mb-1">{{
            t("dashboard.work_orders")
          }}</span>
          <div class="text-3xl font-bold text-surface-900">-</div>
        </div>
        <div
          class="w-10 h-10 rounded-full bg-surface-100 flex items-center justify-center text-surface-500"
        >
          <i class="pi pi-file text-xl"></i>
        </div>
      </div>
    </div>

    <!-- Empty State / Welcome -->
    <div
      class="mt-8 p-12 bg-surface-50 rounded-lg border border-dashed border-surface-300 text-center"
    >
      <h2 class="text-lg font-semibold text-surface-700">
        {{ t("dashboard.welcome_user", { user: authStore.user?.username }) }}
      </h2>
      <p class="text-surface-500">
        {{ t("dashboard.select_option") }}
      </p>
    </div>
  </div>
</template>
