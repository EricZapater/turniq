<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { useCustomersStore } from "../../stores/customers.store";
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
import Tag from "primevue/tag";

const router = useRouter();
const customersStore = useCustomersStore();
const { customers, total, loading, isForbidden } = storeToRefs(customersStore);

// Pagination & filters state (Server-side handling preferred, but keeping robust local state for API wrapping)
const lazyParams = ref({
  page: 0,
  rows: 10,
  sortField: null as string | null,
  sortOrder: null as number | null,
});

// Search and Dialog Filters
const searchQuery = ref("");
const showFilterDialog = ref(false);
const filters = reactive({
  name: "",
  email: "",
  city: "",
});

onMounted(() => {
  loadCustomers();
});

const loadCustomers = () => {
  customersStore.fetchCustomers({
    page: lazyParams.value.page + 1, // API usually 1-indexed
    page_size: lazyParams.value.rows,
    search: searchQuery.value,
    sort_by: lazyParams.value.sortField || undefined,
    sort_desc: lazyParams.value.sortOrder === -1,
  });
};

const onPage = (event: any) => {
  lazyParams.value = event;
  loadCustomers();
};

const onSort = (event: any) => {
  lazyParams.value = event;
  loadCustomers();
};

const onFilter = () => {
  // Implement specific field filtering logic if API supports it,
  // or concatenate to search query for generic search
  loadCustomers();
  showFilterDialog.value = false;
};

const clearFilters = () => {
  filters.name = "";
  filters.email = "";
  filters.city = "";
  searchQuery.value = "";
  loadCustomers();
};

const navigateToCreate = () => {
  router.push("/customers/new");
};

const navigateToDetail = (id: string) => {
  router.push(`/customers/${id}`);
};

const getSeverity = (status: string) => {
  switch (status.toLowerCase()) {
    case "active":
      return "success";
    case "inactive":
      return "warn";
    case "suspended":
      return "danger";
    default:
      return "info";
  }
};
</script>

<template>
  <div class="card p-6 bg-white rounded-lg shadow-sm border border-surface-200">
    <!-- Header Actions -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 tracking-tight">
        Clients
      </h1>
      <div class="flex gap-2">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search" />
          <InputText
            v-model="searchQuery"
            placeholder="Cerca global..."
            @keydown.enter="loadCustomers"
            class="w-64"
          />
        </IconField>
        <Button
          icon="pi pi-filter"
          text
          severity="secondary"
          aria-label="Filter"
          @click="showFilterDialog = true"
          class="border border-surface-300"
        />
        <Button
          label="Nou Client"
          icon="pi pi-plus"
          @click="navigateToCreate"
          :disabled="isForbidden"
        />
      </div>
    </div>

    <!-- Data Table -->
    <DataTable
      :value="customers"
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
      scrollable
      scrollHeight="flex"
    >
      <Column field="name" header="Nom" sortable style="width: 20%"></Column>
      <Column
        field="contact_name"
        header="Contacte"
        style="width: 20%"
      ></Column>
      <Column field="email" header="Email" style="width: 20%"></Column>
      <Column field="phone" header="Telèfon" style="width: 15%"></Column>
      <Column field="status" header="Estat" style="width: 10%">
        <template #body="slotProps">
          <Tag
            :value="slotProps.data.status"
            :severity="getSeverity(slotProps.data.status)"
          />
        </template>
      </Column>
      <Column header="Accions" style="width: 10%" class="text-right">
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
      <template #empty> No s'ha trobat cap client. </template>
    </DataTable>

    <!-- Filter Dialog -->
    <Dialog
      v-model:visible="showFilterDialog"
      modal
      header="Filtres avançats"
      :style="{ width: '30rem' }"
    >
      <div class="flex flex-col gap-4 mb-4">
        <!-- Placeholder for detailed specific filters if supported by API separately -->
        <div class="flex flex-col gap-2">
          <label for="fname" class="font-semibold w-24">Nom</label>
          <InputText
            id="fname"
            v-model="filters.name"
            class="flex-auto"
            autocomplete="off"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label for="femail" class="font-semibold w-24">Email</label>
          <InputText
            id="femail"
            v-model="filters.email"
            class="flex-auto"
            autocomplete="off"
          />
        </div>
      </div>
      <div class="flex justify-end gap-2">
        <Button
          type="button"
          label="Netejar"
          severity="secondary"
          @click="clearFilters"
        ></Button>
        <Button type="button" label="Aplicar" @click="onFilter"></Button>
      </div>
    </Dialog>
  </div>
</template>
