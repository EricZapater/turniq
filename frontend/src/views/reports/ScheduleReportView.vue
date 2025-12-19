<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useToast } from "primevue/usetoast";
import { scheduleApi, type ScheduleEntry } from "../../api/schedule.api";
import { useCustomersStore } from "../../stores/customers.store";
import { useAuthStore } from "../../stores/auth.store";
import { workcentersApi } from "../../api/workcenters.api";
import { operatorsApi } from "../../api/operators.api";
import { shopfloorsApi } from "../../api/shopfloors.api";
import { jobsApi } from "../../api/jobs.api";
import { shiftsApi } from "../../api/shifts.api";
import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";

// Components
import Button from "primevue/button";
import Calendar from "primevue/calendar";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Select from "primevue/select";
import Tag from "primevue/tag";
import Card from "primevue/card";
import Toolbar from "primevue/toolbar";

const toast = useToast();
const { t } = useI18n();
const customersStore = useCustomersStore();
const authStore = useAuthStore();

const loading = ref(false);
const entries = ref<ScheduleEntry[]>([]);
const dates = ref<Date[] | null>(null);

// Filters
const selectedCustomer = ref<string | null>(null);
const selectedShopfloor = ref<string | null>(null); // Optional helper filter
const selectedWorkcenter = ref<string | null>(null);
const selectedOperator = ref<string | null>(null);

// Resources options
const workcenters = ref<any[]>([]);
const operators = ref<any[]>([]);
const shopfloors = ref<any[]>([]);
const jobs = ref<any[]>([]);
const shifts = ref<any[]>([]);

const dt = ref();

onMounted(async () => {
  loading.value = true;
  if (authStore.isAdmin) {
    await customersStore.fetchCustomers({});
  }

  // Load resources based on context
  await loadResources();
  await search(); // Initial load
  loading.value = false;
});

const loadResources = async () => {
  const customerId =
    selectedCustomer.value ||
    (authStore.isAdmin ? null : authStore.user?.customer_id);

  // If we have a customerId, we filter.
  // If we don't (and we are admin), we want ALL resources to resolve names in the table working across customers.
  // If we are not admin and have no customerId, we shouldn't see anything (but previous logic handles that via authStore).

  const params = customerId ? { customer_id: customerId } : {};

  // For Admin observing everything (customerId=null): fetch ALL resources.
  // For Admin filtering (customerId=XYZ): fetch filtered.
  // For User (customerId=ABC): fetch filtered.

  if (authStore.isAdmin || customerId) {
    const sfRes = await shopfloorsApi.list(params);
    shopfloors.value = sfRes.data;

    const wcRes = await workcentersApi.list(params);
    workcenters.value = wcRes.data;

    const opRes = await operatorsApi.list(params);
    operators.value = opRes.data;

    const jRes = await jobsApi.list(params);
    jobs.value = jRes.data;

    const sRes = await shiftsApi.list(params);
    shifts.value = sRes.data;
  } else {
    workcenters.value = [];
    operators.value = [];
    shopfloors.value = [];
    jobs.value = [];
    shifts.value = [];
  }
};

// Watch Customer Change (Admin)
watch(selectedCustomer, async () => {
  selectedShopfloor.value = null;
  selectedWorkcenter.value = null;
  selectedOperator.value = null;
  await loadResources();
});

const customerOptions = computed(() => {
  return customersStore.customers.map((c) => ({
    label: c.name,
    value: c.id,
  }));
});

// Helper lookups for table display
const getCustomerName = (id: string) => {
  if (!id) return "-";
  const c = customersStore.customers.find((c) => c.id === id);
  return c ? c.name : "...";
};
const getWorkcenterName = (id: string) => {
  if (!id) return "-";
  const w = workcenters.value.find((w) => w.id === id);
  return w ? w.name : "...";
};
const getOperatorName = (id?: string | null) => {
  if (!id) return "-";
  const o = operators.value.find((o) => o.id === id);
  return o ? `${o.name} ${o.surname}` : "...";
};
const getShopfloorName = (id: string) => {
  if (!id) return "-";
  const sf = shopfloors.value.find((s) => s.id === id);
  return sf ? sf.name : "...";
};
const getShiftName = (id: string) => {
  if (!id) return "-";
  const s = shifts.value.find((sh) => sh.id === id);
  return s ? s.name : "...";
};
const getJobInfo = (id?: string | null) => {
  if (!id) return "-";
  const j = jobs.value.find((j) => j.id === id);
  return j ? `${j.job_code} (${j.product_code})` : "...";
};
const getEstDuration = (id?: string | null) => {
  if (!id) return "-";
  const j = jobs.value.find((j) => j.id === id);
  return j ? `${j.estimated_duration}m` : "-";
};

const search = async () => {
  loading.value = true;
  try {
    const params: any = {};
    if (selectedCustomer.value) params.customer_id = selectedCustomer.value;
    if (selectedShopfloor.value) params.shopfloor_id = selectedShopfloor.value;
    if (selectedWorkcenter.value)
      params.workcenter_id = selectedWorkcenter.value;
    if (selectedOperator.value && authStore.isAdmin)
      params.operator_id = selectedOperator.value; // Restricted param

    if (dates.value && dates.value[0]) {
      params.from = dates.value[0].toISOString().split("T")[0];
      if (dates.value[1]) {
        params.to = dates.value[1].toISOString().split("T")[0];
      }
    }

    const res = await scheduleApi.list(params);
    entries.value = res || [];
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
  selectedShopfloor.value = null;
  selectedWorkcenter.value = null;
  selectedOperator.value = null;
  search();
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleDateString("ca-ES");
};

const exportCSV = () => {
  dt.value.exportCSV();
};

const exportPDF = () => {
  const doc = new jsPDF();

  doc.setFontSize(18);
  doc.text(t("reports.schedule_title"), 14, 22);
  doc.setFontSize(11);
  doc.text(
    `${t("reports.generated_at")}: ${new Date().toLocaleString("ca-ES")}`,
    14,
    30
  );

  const head = [
    [
      t("reports.date"),
      t("customers.title"),
      t("shopfloors.title"),
      t("shifts.title"),
      t("workcenters.title"),
      t("operators.title"),
      t("jobs.title"),
      t("reports.check_in"),
      t("reports.check_out"),
      t("reports.expected"),
      t("common.status"),
    ],
  ];
  const body = entries.value.map((e) => [
    formatDate(e.date),
    getCustomerName(e.customer_id),
    getShopfloorName(e.shopfloor_id),
    getShiftName(e.shift_id),
    getWorkcenterName(e.workcenter_id),
    getOperatorName(e.operator_id),
    getJobInfo(e.job_id),
    e.start_time || "-",
    e.end_time || "-",
    getEstDuration(e.job_id),
    e.is_completed ? t("reports.completed") : t("reports.pending"),
  ]);

  autoTable(doc, {
    startY: 40,
    head: head,
    body: body,
  });

  doc.save("planificacio_turniq.pdf");
};
</script>

<template>
  <div class="flex flex-col h-full gap-4">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">
        {{ t("reports.schedule_title") }}
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

          <!-- Shopfloor (Optional Helper) -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-surface-700">{{
              t("shopfloors.title")
            }}</label>
            <Select
              v-model="selectedShopfloor"
              :options="shopfloors"
              optionLabel="name"
              optionValue="id"
              :placeholder="t('reports.all')"
              showClear
              class="w-full"
            />
          </div>

          <!-- Workcenter -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-surface-700">{{
              t("workcenters.title")
            }}</label>
            <Select
              v-model="selectedWorkcenter"
              :options="workcenters"
              optionLabel="name"
              optionValue="id"
              :placeholder="t('reports.all')"
              showClear
              filter
              class="w-full"
            />
          </div>

          <!-- Operator (Admin Only) -->
          <div v-if="authStore.isAdmin" class="flex flex-col gap-2">
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

          <div class="flex gap-2 col-span-1 md:col-span-4 justify-end">
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

          <Column field="date" :header="t('reports.date')" sortable>
            <template #body="slotProps">
              {{ formatDate(slotProps.data.date) }}
            </template>
          </Column>

          <Column
            v-if="authStore.isAdmin"
            field="customer_id"
            :header="t('customers.title')"
            sortable
          >
            <template #body="slotProps">
              {{ getCustomerName(slotProps.data.customer_id) }}
            </template>
          </Column>

          <Column field="shopfloor_id" :header="t('shopfloors.title')" sortable>
            <template #body="slotProps">
              {{ getShopfloorName(slotProps.data.shopfloor_id) }}
            </template>
          </Column>

          <Column field="shift_id" :header="t('shifts.title')" sortable>
            <template #body="slotProps">
              {{ getShiftName(slotProps.data.shift_id) }}
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

          <Column field="operator_id" :header="t('operators.title')" sortable>
            <template #body="slotProps">
              {{ getOperatorName(slotProps.data.operator_id) }}
            </template>
          </Column>

          <Column field="job_id" :header="t('jobs.title')" sortable>
            <template #body="slotProps">
              {{ getJobInfo(slotProps.data.job_id) }}
            </template>
          </Column>

          <Column field="start_time" :header="t('reports.check_in')" sortable>
            <template #body="slotProps">
              {{ slotProps.data.start_time || "-" }}
            </template>
          </Column>

          <Column field="end_time" :header="t('reports.check_out')" sortable>
            <template #body="slotProps">
              {{ slotProps.data.end_time || "-" }}
            </template>
          </Column>

          <Column :header="t('reports.expected')">
            <template #body="slotProps">
              {{ getEstDuration(slotProps.data.job_id) }}
            </template>
          </Column>

          <Column field="is_completed" :header="t('common.status')" sortable>
            <template #body="slotProps">
              <Tag
                :value="
                  slotProps.data.is_completed
                    ? t('reports.completed')
                    : t('reports.pending')
                "
                :severity="slotProps.data.is_completed ? 'success' : 'warn'"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>
