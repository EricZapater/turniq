<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useToast } from "primevue/usetoast";
import { timeentriesApi, type TimeEntry } from "../../api/timeentries.api";
import { useCustomersStore } from "../../stores/customers.store";
import { useAuthStore } from "../../stores/auth.store";
import { workcentersApi } from "../../api/workcenters.api";
import { operatorsApi } from "../../api/operators.api";
import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";

// Components
import Button from "primevue/button";
import Calendar from "primevue/calendar";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Select from "primevue/select";
import Card from "primevue/card";
import Toolbar from "primevue/toolbar";

const toast = useToast();
const { t } = useI18n();
const customersStore = useCustomersStore();
const authStore = useAuthStore();

const loading = ref(false);
const entries = ref<TimeEntry[]>([]);
const dates = ref<Date[] | null>(null);

// Filters
const selectedCustomer = ref<string | null>(null);
const selectedOperator = ref<string | null>(null);

// Resources options
const workcenters = ref<any[]>([]);
const operators = ref<any[]>([]);

const dt = ref();

onMounted(async () => {
  loading.value = true;
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({});
  }

  await loadResources();
  await search();
  loading.value = false;
});

const loadResources = async () => {
  const customerId =
    selectedCustomer.value ||
    (authStore.isAdmin ? null : authStore.user?.customer_id);

  const params = customerId ? { customer_id: customerId } : {};

  if (authStore.isAdmin || customerId) {
    const wcRes = await workcentersApi.list(params);
    workcenters.value = wcRes.data;

    const opRes = await operatorsApi.list(params);
    operators.value = opRes.data;
  } else {
    workcenters.value = [];
    operators.value = [];
  }
};

watch(selectedCustomer, async () => {
  selectedOperator.value = null;
  await loadResources();
});

const customerOptions = computed(() => {
  return customersStore.customers.map((c) => ({
    label: c.name,
    value: c.id,
  }));
});

const getWorkcenterName = (id?: string) => {
  if (!id) return "-";
  const w = workcenters.value.find((w) => w.id === id);
  return w ? w.name : "-";
};
const getOperatorName = (id?: string) => {
  if (!id) return "-";
  const o = operators.value.find((o) => o.id === id);
  return o ? `${o.name} ${o.surname}` : "...";
};

const search = async () => {
  loading.value = true;
  try {
    const params: any = {};
    if (selectedCustomer.value) params.customer_id = selectedCustomer.value;
    if (selectedOperator.value) params.operator_id = selectedOperator.value;

    if (dates.value && dates.value[0]) {
      params.from = dates.value[0].toISOString().split("T")[0];
      if (dates.value[1]) {
        params.to = dates.value[1].toISOString().split("T")[0];
      }
    }

    const res = await timeentriesApi.list(params);
    entries.value = res.data || [];
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: t("common.error"),
      detail: t("common.error"),
    });
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  dates.value = null;
  selectedCustomer.value = null;
  selectedOperator.value = null;
  search();
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleDateString("ca-ES");
};

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleString("ca-ES", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
};

const calculateDuration = (start: string, end?: string) => {
  if (!end) return "En curs";
  const diff = new Date(end).getTime() - new Date(start).getTime();
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

  return `${hours}h ${minutes}m`;
};

const exportCSV = () => {
  dt.value.exportCSV();
};

const exportPDF = () => {
  const doc = new jsPDF();

  doc.setFontSize(18);
  doc.text(t("reports.hours_title"), 14, 22);
  doc.setFontSize(11);
  doc.text(
    `${t("reports.generated_at")}: ${new Date().toLocaleString("ca-ES")}`,
    14,
    30
  );

  const head = [
    [
      t("reports.date"),
      t("operators.title"),
      t("workcenters.title"),
      t("reports.check_in"),
      t("reports.check_out"),
      t("reports.duration"),
    ],
  ];
  const body = entries.value.map((e) => [
    formatDate(e.check_in),
    getOperatorName(e.operator_id),
    getWorkcenterName(e.workcenter_id),
    formatDateTime(e.check_in),
    formatDateTime(e.check_out || ""),
    calculateDuration(e.check_in, e.check_out),
  ]);

  autoTable(doc, {
    startY: 40,
    head: head,
    body: body,
  });

  doc.save("registre_horari_turniq.pdf");
};
</script>

<template>
  <div class="flex flex-col h-full gap-4">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">
        {{ t("reports.hours_title") }}
      </h1>
    </div>

    <Card>
      <template #content>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6 items-end">
          <!-- Customer (Admin Only) -->
          <div v-if="authStore.isAdmin" class="flex flex-col gap-2">
            <label for="customer" class="font-medium text-surface-700">{{
              t("customers.title")
            }}</label>
            <Select
              v-model="selectedCustomer"
              :options="customerOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('reports.all')"
              showClear
              filter
              class="w-full"
            />
          </div>

          <!-- Date -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-surface-700">{{
              t("reports.period")
            }}</label>
            <Calendar
              v-model="dates"
              selectionMode="range"
              :manualInput="false"
              showIcon
              :placeholder="t('reports.select_range')"
            />
          </div>

          <!-- Operator -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-surface-700">{{
              t("operators.title")
            }}</label>
            <Select
              v-model="selectedOperator"
              :options="operators"
              :optionLabel="(opt) => `${opt.name} ${opt.surname}`"
              optionValue="id"
              :placeholder="t('reports.all')"
              showClear
              filter
              class="w-full"
            />
          </div>

          <div class="flex gap-2 justify-end col-span-1 md:col-span-1">
            <Button
              :label="t('common.search')"
              icon="pi pi-search"
              @click="search"
              :loading="loading"
            />
            <Button
              :label="t('common.clear')"
              icon="pi pi-filter-slash"
              severity="secondary"
              @click="clearFilters"
            />
          </div>
        </div>

        <Toolbar class="mb-4">
          <template #start>
            <span class="text-surface-500"
              >{{ t("reports.total_records") }}: {{ entries.length }}</span
            >
          </template>
          <template #end>
            <Button
              :label="t('reports.export_csv')"
              icon="pi pi-file-excel"
              severity="success"
              class="mr-2"
              @click="exportCSV"
            />
            <Button
              :label="t('reports.export_pdf')"
              icon="pi pi-file-pdf"
              severity="danger"
              @click="exportPDF"
            />
          </template>
        </Toolbar>

        <DataTable
          ref="dt"
          :value="entries"
          :loading="loading"
          paginator
          :rows="10"
          :rowsPerPageOptions="[10, 25, 50, 100]"
          stripedRows
          removableSort
          class="p-datatable-sm"
        >
          <template #empty>{{ t("reports.no_entries") }}</template>

          <Column field="check_in" :header="t('reports.date')" sortable>
            <template #body="slotProps">
              {{ formatDate(slotProps.data.check_in) }}
            </template>
          </Column>

          <Column field="operator_id" :header="t('operators.title')" sortable>
            <template #body="slotProps">
              {{ getOperatorName(slotProps.data.operator_id) }}
            </template>
          </Column>

          <Column
            field="workcenter_id"
            :header="t('workcenters.title')"
            sortable
          >
            <template #body="slotProps">
              {{ getWorkcenterName(slotProps.data.workcenter_id) }}
            </template>
          </Column>

          <Column field="check_in" :header="t('reports.check_in')" sortable>
            <template #body="slotProps">
              {{ formatDateTime(slotProps.data.check_in) }}
            </template>
          </Column>

          <Column field="check_out" :header="t('reports.check_out')" sortable>
            <template #body="slotProps">
              {{ formatDateTime(slotProps.data.check_out) }}
            </template>
          </Column>

          <Column :header="t('reports.duration')">
            <template #body="slotProps">
              {{
                calculateDuration(
                  slotProps.data.check_in,
                  slotProps.data.check_out
                )
              }}
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>
