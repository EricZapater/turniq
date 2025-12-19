<script setup lang="ts">
import { ref, onMounted, reactive, computed } from "vue";
import { useShiftsStore } from "../../stores/shifts.store";
import { useShopfloorsStore } from "../../stores/shopfloors.store";
import { useCustomersStore } from "../../stores/customers.store";
import { useAuthStore } from "../../stores/auth.store";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";

// PrimeVue Components
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Button from "primevue/button";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import IconField from "primevue/iconfield";
import InputIcon from "primevue/inputicon";
import Select from "primevue/select";
import Tag from "primevue/tag";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const router = useRouter();
const shiftsStore = useShiftsStore();
const shopfloorsStore = useShopfloorsStore();
const customersStore = useCustomersStore();
const authStore = useAuthStore();

const { shifts, total, loading } = storeToRefs(shiftsStore);
const { shopfloors } = storeToRefs(shopfloorsStore);
const { customers } = storeToRefs(customersStore);

// Pagination
const lazyParams = ref({
  page: 0,
  rows: 10,
  sortField: null as string | null,
  sortOrder: null as number | null,
});

// Search and Filter
const searchQuery = ref("");
const showFilterDialog = ref(false);
const filters = reactive({
  customer_id: null as string | null,
  shop_floor_id: null as string | null,
});

onMounted(async () => {
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({ page_size: 100 });
  }
  await shopfloorsStore.fetchShopfloors({ page_size: 100 });
  loadShifts();
});

const loadShifts = () => {
  // If we are filtering by shopfloor, we can use fetchShiftsByShopfloor for efficiency OR just use the generic list with params if backend supports filtering by shopfloor in FindAll (which it currently doesn't, service.FindByCustomer returns all).
  // Wait, backend FindAll calls FindByCustomer. It doesn't filter by shopfloor query param.
  // Frontend filtering is safer for now if the list isn't huge, OR we rely on the specific endpoint for shopfloor.
  // BUT the user wants "same system", usually that means `list` with params.
  // Since backend FindAll ignores shopfloor param, we will do client side filtering effectively or just fetch all.

  // Actually, if a shopfloor is selected in filter, we COULD use fetchShiftsByShopfloor.
  // But for the main list, we use fetchShifts.

  if (filters.shop_floor_id) {
    shiftsStore.fetchShiftsByShopfloor(filters.shop_floor_id, {
      page: lazyParams.value.page + 1,
      page_size: lazyParams.value.rows,
      search: searchQuery.value,
      sort_by: lazyParams.value.sortField || undefined,
      sort_desc: lazyParams.value.sortOrder === -1,
    });
  } else {
    // If Admin & Customer selected? Backend FindAll supports customer_id query param?
    // Handler FindAll checks middleware, else query param.
    // So we can pass customer_id if admin.
    shiftsStore.fetchShifts({
      page: lazyParams.value.page + 1,
      page_size: lazyParams.value.rows,
      search: searchQuery.value,
      sort_by: lazyParams.value.sortField || undefined,
      sort_desc: lazyParams.value.sortOrder === -1,
      customer_id: filters.customer_id || undefined,
    });
  }
};

const onPage = (event: any) => {
  lazyParams.value = event;
  loadShifts();
};

const onSort = (event: any) => {
  lazyParams.value = event;
  loadShifts();
};

const onFilter = () => {
  loadShifts();
  showFilterDialog.value = false;
};

const clearFilters = () => {
  filters.shop_floor_id = null;
  filters.customer_id = null;
  searchQuery.value = "";
  loadShifts();
};

const navigateToCreate = () => {
  router.push("/shifts/new");
};

const navigateToDetail = (id: string) => {
  router.push(`/shifts/${id}`);
};

const getShopfloorName = (id?: string) => {
  if (!id) return "-";
  const sf = shopfloors.value.find((s) => s.id === id);
  return sf ? sf.name : id;
};

const filteredShopfloors = computed(() => {
  if (authStore.isAdmin && filters.customer_id) {
    return shopfloors.value.filter(
      (s) => s.customer_id === filters.customer_id
    );
  }
  return shopfloors.value;
});

// Format time string
const formatTime = (timeStr: string) => {
  if (!timeStr) return "";
  // Expected format: "0000-01-01T06:00:00Z"
  if (timeStr.includes("T")) {
    const timePart = timeStr.split("T")[1];
    // Remove 'Z' or offset if present and seconds
    return timePart.substring(0, 5);
  }
  return timeStr.substring(0, 5);
};
</script>

<template>
  <div class="card p-6 bg-white rounded-lg shadow-sm border border-surface-200">
    <!-- Header Actions -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 tracking-tight">
        {{ t("shifts.title") }}
      </h1>
      <div class="flex gap-2">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search" />
          <InputText
            v-model="searchQuery"
            :placeholder="t('common.search') + '...'"
            @keydown.enter="loadShifts"
            class="w-64"
          />
        </IconField>
        <Button
          icon="pi pi-filter"
          text
          severity="secondary"
          @click="showFilterDialog = true"
          class="border border-surface-300"
          :class="{
            'bg-primary-50 text-primary-600':
              filters.shop_floor_id || filters.customer_id,
          }"
        />
        <Button
          :label="t('shifts.new')"
          icon="pi pi-plus"
          @click="navigateToCreate"
        />
      </div>
    </div>

    <!-- Data Table -->
    <DataTable
      :value="shifts"
      :lazy="true"
      :paginator="true"
      :rows="lazyParams.rows"
      :totalRecords="total"
      :loading="loading"
      @page="onPage"
      @sort="onSort"
      tableStyle="min-width: 50rem"
      class="p-datatable-sm"
      stripedRows
    >
      <Column
        field="name"
        :header="t('shifts.name')"
        sortable
        style="width: 25%"
      >
        <template #body="slotProps">
          <div class="flex items-center gap-2">
            <div
              :style="{ backgroundColor: slotProps.data.color }"
              class="w-4 h-4 rounded-full border border-surface-200 shadow-sm"
            ></div>
            <span class="font-medium">{{ slotProps.data.name }}</span>
          </div>
        </template>
      </Column>
      <Column
        field="shopfloor_id"
        :header="t('shifts.shopfloor')"
        sortable
        style="width: 20%"
      >
        <template #body="slotProps">
          {{ getShopfloorName(slotProps.data.shopfloor_id) }}
        </template>
      </Column>
      <Column
        field="start_time"
        :header="t('shifts.start_time')"
        sortable
        style="width: 15%"
      >
        <template #body="slotProps">
          {{ formatTime(slotProps.data.start_time) }}
        </template>
      </Column>
      <Column
        field="end_time"
        :header="t('shifts.end_time')"
        sortable
        style="width: 15%"
      >
        <template #body="slotProps">
          {{ formatTime(slotProps.data.end_time) }}
        </template>
      </Column>
      <Column
        field="is_active"
        :header="t('common.status')"
        sortable
        style="width: 10%"
      >
        <template #body="slotProps">
          <Tag
            :severity="slotProps.data.is_active ? 'success' : 'secondary'"
            :value="
              slotProps.data.is_active
                ? t('common.active')
                : t('common.inactive')
            "
            rounded
          />
        </template>
      </Column>
      <Column
        :header="t('common.actions')"
        style="width: 15%"
        class="text-right"
      >
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            text
            rounded
            severity="secondary"
            @click="navigateToDetail(slotProps.data.id)"
          />
        </template>
      </Column>
      <template #empty>
        {{ t("common.no_records") || "No s'ha trobat cap torn." }}
      </template>
    </DataTable>

    <!-- Filter Dialog -->
    <Dialog
      v-model:visible="showFilterDialog"
      modal
      header="Filters"
      :style="{ width: '30rem' }"
    >
      <div class="flex flex-col gap-4 mb-4">
        <!-- Admin Customer Filter -->
        <div class="flex flex-col gap-2" v-if="authStore.isAdmin">
          <label for="filterCustomer" class="font-semibold">{{
            t("customers.title")
          }}</label>
          <Select
            id="filterCustomer"
            v-model="filters.customer_id"
            :options="customers"
            optionLabel="name"
            optionValue="id"
            placeholder="Tots els clients"
            showClear
            filter
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-2">
          <label for="filterShopfloor" class="font-semibold">{{
            t("shifts.shopfloor")
          }}</label>
          <Select
            id="filterShopfloor"
            v-model="filters.shop_floor_id"
            :options="filteredShopfloors"
            optionLabel="name"
            optionValue="id"
            placeholder="Totes les plantes"
            showClear
            filter
            class="w-full"
            :disabled="authStore.isAdmin && !filters.customer_id"
          />
        </div>
      </div>
      <div class="flex justify-end gap-2">
        <Button
          type="button"
          :label="t('common.clear')"
          severity="secondary"
          @click="clearFilters"
        ></Button>
        <Button
          type="button"
          :label="t('common.search')"
          @click="onFilter"
        ></Button>
      </div>
    </Dialog>
  </div>
</template>
