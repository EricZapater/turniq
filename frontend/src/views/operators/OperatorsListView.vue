<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { useOperatorsStore } from "../../stores/operators.store";
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
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const router = useRouter();
const operatorsStore = useOperatorsStore();
const { operators, total, loading } = storeToRefs(operatorsStore);

const lazyParams = ref({
  page: 0,
  rows: 10,
  sortField: null as string | null,
  sortOrder: null as number | null,
});

const searchQuery = ref("");
const showFilterDialog = ref(false);
const filters = reactive({
  name: "",
});

onMounted(() => {
  loadOperators();
});

const loadOperators = () => {
  operatorsStore.fetchOperators({
    page: lazyParams.value.page + 1,
    page_size: lazyParams.value.rows,
    search: searchQuery.value,
    sort_by: lazyParams.value.sortField || undefined,
    sort_desc: lazyParams.value.sortOrder === -1,
  });
};

const onPage = (event: any) => {
  lazyParams.value = event;
  loadOperators();
};

const onSort = (event: any) => {
  lazyParams.value = event;
  loadOperators();
};

const onFilter = () => {
  loadOperators();
  showFilterDialog.value = false;
};

const clearFilters = () => {
  filters.name = "";
  searchQuery.value = "";
  loadOperators();
};

const navigateToCreate = () => {
  router.push("/operators/new");
};

const navigateToDetail = (id: string) => {
  router.push(`/operators/${id}`);
};
</script>

<template>
  <div class="card p-6 bg-white rounded-lg shadow-sm border border-surface-200">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-surface-900 tracking-tight">
        {{ t("operators.title") }}
      </h1>
      <div class="flex gap-2">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search" />
          <InputText
            v-model="searchQuery"
            :placeholder="t('common.search') + '...'"
            @keydown.enter="loadOperators"
            class="w-64"
          />
        </IconField>
        <Button
          icon="pi pi-filter"
          text
          severity="secondary"
          @click="showFilterDialog = true"
          class="border border-surface-300"
        />
        <Button
          :label="t('operators.new')"
          icon="pi pi-plus"
          @click="navigateToCreate"
        />
      </div>
    </div>

    <DataTable
      :value="operators"
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
        field="code"
        :header="t('operators.code')"
        sortable
        style="width: 15%"
      ></Column>
      <Column
        field="name"
        :header="t('operators.name')"
        sortable
        style="width: 25%"
      ></Column>
      <Column
        field="surname"
        :header="t('operators.surname')"
        sortable
        style="width: 25%"
      ></Column>
      <Column field="is_active" :header="t('common.status')" style="width: 15%">
        <template #body="slotProps">
          <Tag
            :value="
              slotProps.data.is_active
                ? t('common.active')
                : t('common.inactive')
            "
            :severity="slotProps.data.is_active ? 'success' : 'danger'"
          />
        </template>
      </Column>
      <Column
        :header="t('common.actions')"
        style="width: 10%"
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
        {{ t("common.no_records") || "No s'ha trobat cap operari." }}
      </template>
    </DataTable>

    <Dialog
      v-model:visible="showFilterDialog"
      modal
      header="Filters"
      :style="{ width: '30rem' }"
    >
      <div class="flex flex-col gap-4 mb-4">
        <div class="flex flex-col gap-2">
          <label for="fname" class="font-semibold w-24">{{
            t("operators.name")
          }}</label>
          <InputText
            id="fname"
            v-model="filters.name"
            class="flex-auto"
            autocomplete="off"
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
