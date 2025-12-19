<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useToast } from "primevue/usetoast";
import { paymentsApi, type Payment } from "../../api/payments.api";
import { useCustomersStore } from "../../stores/customers.store";
import { useAuthStore } from "../../stores/auth.store";
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
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";

const toast = useToast();
const { t } = useI18n();
const customersStore = useCustomersStore();
const authStore = useAuthStore();

const loading = ref(false);
const payments = ref<Payment[]>([]);
const dates = ref<Date[] | null>(null);
const selectedCustomer = ref<string | null>(null);
const dt = ref();

onMounted(async () => {
  if (!authStore.isAdmin) {
    // Should be redirected by router guard, but safe check
    return;
  }
  loading.value = true;
  await customersStore.fetchCustomers({});
  await search(); // Initial load
  loading.value = false;
});

const customerOptions = computed(() => {
  return customersStore.customers.map((c) => ({
    label: c.name,
    value: c.id,
  }));
});

const getCustomerName = (id: string) => {
  const c = customersStore.customers.find((c) => c.id === id);
  return c ? c.name : id;
};

const search = async () => {
  loading.value = true;
  try {
    const params: any = {};
    if (selectedCustomer.value) {
      params.customer_id = selectedCustomer.value;
    }
    if (dates.value && dates.value[0]) {
      params.from = dates.value[0].toISOString().split("T")[0];
      if (dates.value[1]) {
        params.to = dates.value[1].toISOString().split("T")[0];
      }
    }

    const res = await paymentsApi.list(params);
    // Response wrapper has .data or is the list?
    // paymentsApi.list returns response.data which is PaymentListResponse { data: [], message: "" }
    payments.value = res.data || [];
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
  search();
};

const formatCurrency = (amount: number, currency: string) => {
  return new Intl.NumberFormat("ca-ES", {
    style: "currency",
    currency: currency || "EUR",
  }).format(amount);
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleDateString("ca-ES");
};

const getStatusSeverity = (status: string) => {
  switch (status) {
    case "succeeded":
    case "paid":
      return "success";
    case "pending":
      return "warn";
    case "failed":
      return "danger";
    default:
      return "info";
  }
};

const exportCSV = () => {
  dt.value.exportCSV();
};

const exportPDF = () => {
  const doc = new jsPDF();

  // Title
  doc.setFontSize(18);
  doc.text(t("reports.billing_title"), 14, 22);
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
      t("customers.quantity") || t("reports.amount"),
      t("reports.method"),
      t("common.status"),
    ],
  ];
  const body = payments.value.map((p) => [
    formatDate(p.paid_at),
    getCustomerName(p.customer_id),
    formatCurrency(p.amount, p.currency),
    p.payment_method,
    p.status,
  ]);

  autoTable(doc, {
    startY: 40,
    head: head,
    body: body,
  });

  doc.save("facturacio_turniq.pdf");
};
</script>

<template>
  <div class="flex flex-col h-full gap-4">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-surface-900 dark:text-surface-0">
        {{ t("reports.billing_title") }}
      </h1>
    </div>

    <Card>
      <template #content>
        <div class="flex flex-col md:flex-row gap-4 mb-6 items-end">
          <!-- Filters -->
          <div class="flex flex-col gap-2 flex-1">
            <label for="dates" class="font-medium text-surface-700">{{
              t("reports.period")
            }}</label>
            <Calendar
              v-model="dates"
              selectionMode="range"
              :manualInput="false"
              showIcon
              :placeholder="t('reports.select_range')"
              inputId="dates"
            />
          </div>

          <div class="flex flex-col gap-2 flex-1">
            <label for="customer" class="font-medium text-surface-700">{{
              t("customers.title")
            }}</label>
            <Select
              v-model="selectedCustomer"
              :options="customerOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('reports.all_customers')"
              showClear
              inputId="customer"
              filter
              class="w-full"
            />
          </div>

          <div class="flex gap-2">
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
            <!-- Add Summary Stats here if needed? -->
            <span class="text-surface-500"
              >{{ t("reports.total_records") }}: {{ payments.length }}</span
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
          :value="payments"
          :loading="loading"
          paginator
          :rows="10"
          :rowsPerPageOptions="[10, 25, 50, 100]"
          stripedRows
          removableSort
          class="p-datatable-sm"
        >
          <template #empty>{{ t("reports.no_payments") }}</template>

          <Column field="paid_at" :header="t('reports.payment_date')" sortable>
            <template #body="slotProps">
              {{ formatDate(slotProps.data.paid_at) }}
            </template>
          </Column>

          <Column field="customer_id" :header="t('customers.title')" sortable>
            <template #body="slotProps">
              {{ getCustomerName(slotProps.data.customer_id) }}
            </template>
          </Column>

          <Column field="amount" :header="t('reports.amount')" sortable>
            <template #body="slotProps">
              <span class="font-bold">
                {{
                  formatCurrency(slotProps.data.amount, slotProps.data.currency)
                }}
              </span>
            </template>
          </Column>

          <Column
            field="payment_method"
            :header="t('reports.method')"
            sortable
          ></Column>

          <Column field="status" :header="t('common.status')" sortable>
            <template #body="slotProps">
              <Tag
                :value="slotProps.data.status"
                :severity="getStatusSeverity(slotProps.data.status)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>
