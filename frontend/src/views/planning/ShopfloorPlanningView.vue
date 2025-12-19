<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";
import { usePlanningStore } from "../../stores/planning.store";
import { useShopfloorsStore } from "../../stores/shopfloors.store";
import { useAuthStore } from "../../stores/auth.store";
import { storeToRefs } from "pinia";
import DatePicker from "primevue/datepicker";
import Select from "primevue/select";
import Button from "primevue/button";

import { useToast } from "primevue/usetoast";

const planningStore = usePlanningStore();
const shopfloorsStore = useShopfloorsStore();
const authStore = useAuthStore();
const toast = useToast();
const { t } = useI18n();

const { resources, planning, loading, unsavedChanges } =
  storeToRefs(planningStore);
const { shopfloors } = storeToRefs(shopfloorsStore);

const savePlanning = async () => {
  if (!selectedShopfloor.value) return;
  try {
    const dateStr = formatDateForApi(selectedDate.value);
    await planningStore.saveChanges(selectedShopfloor.value, dateStr);
    toast.add({
      severity: "success",
      summary: t("planning.saved"),
      detail: t("planning.saved_success"),
      life: 2000,
    });
  } catch (error) {
    toast.add({
      severity: "error",
      summary: t("common.error"),
      detail: t("planning.error_saving"),
      life: 3000,
    });
  }
};

const selectedDate = ref(new Date());
const selectedShopfloor = ref<string | null>(null);

// Drag State
const draggedItem = ref<any>(null); // { type: 'operator' | 'workcenter' | 'job', data: Object, origin: 'palette' | 'planning' }

onMounted(async () => {
  await shopfloorsStore.fetchShopfloors({ page_size: 100 });
  if (shopfloors.value.length > 0) {
    selectedShopfloor.value = shopfloors.value[0].id; // Default
  }
});

watch(
  [selectedShopfloor, selectedDate],
  async ([sf, date]) => {
    if (sf && date) {
      await planningStore.loadResources(sf);
      await planningStore.loadPlanning(sf, date);
    }
  },
  { immediate: true }
);

// Drag Handlers for Palette Items
const onDragStart = (event: DragEvent, type: string, item: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = "copy";
    event.dataTransfer.setData("type", type);
    event.dataTransfer.setData("item", JSON.stringify(item));
    event.dataTransfer.setData("origin", "palette");
    draggedItem.value = { type, item, origin: "palette" };
  }
};

// Drop Handlers

const onDragOver = (event: DragEvent) => {
  event.preventDefault(); // Necessary to allow dropping
};

const formatTime = (timeStr?: string) => {
  if (!timeStr) return "--:--";
  if (timeStr.includes("T")) {
    return timeStr.split("T")[1].substring(0, 5);
  }
  return timeStr?.substring(0, 5) || "--:--";
};

const formatDateForApi = (date: Date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

// Drag Handlers for Board Items
const onDragStartFromBoard = (event: DragEvent, type: string, item: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("type", type);
    event.dataTransfer.setData("item", JSON.stringify(item));
    event.dataTransfer.setData("origin", "planning");
    draggedItem.value = { type, item, origin: "planning" };
  }
};

// Tree Helper
const getShiftPlanning = (shiftId: string) => {
  const entries = planning.value.filter((e) => e.shift_id === shiftId);

  const ops = new Map<string, any>();

  entries.forEach((e) => {
    if (!e.operator_id) return;
    if (!ops.has(e.operator_id)) {
      ops.set(e.operator_id, {
        id: e.operator_id,
        name: getResourceName("operator", e.operator_id),
        entry_id: e.id, // Root entry ID if mostly placeholders
        wcs: new Map<string, any>(),
      });
    }

    const opGroup = ops.get(e.operator_id);

    if (e.workcenter_id) {
      if (!opGroup.wcs.has(e.workcenter_id)) {
        opGroup.wcs.set(e.workcenter_id, {
          id: e.workcenter_id,
          name: getResourceName("workcenter", e.workcenter_id),
          jobs: [],
          order: 0, // Default
        });
      }
      const wcGroup = opGroup.wcs.get(e.workcenter_id);
      if (e.job_id) {
        const job = resources.value.jobs.find((j) => j.id === e.job_id);
        wcGroup.jobs.push({
          id: e.job_id,
          entry_id: e.id,
          name: job?.name || t("planning.unknown_job"),
          job_code: job?.job_code || "",
          product_code: job?.product_code || "",
          estimated_duration: job?.estimated_duration || 0,
          order: e.order,
        });
      } else {
        wcGroup.entry_id = e.id;
        wcGroup.order = e.order; // Capture order for WC
      }
    }
  });

  return Array.from(ops.values()).map((op: any) => {
    op.workcenters = Array.from(op.wcs.values()).map((wc: any) => {
      wc.jobs.sort((a: any, b: any) => a.order - b.order);
      return wc;
    });
    // Sort Workcenters by order
    op.workcenters.sort((a: any, b: any) => a.order - b.order);
    return op;
  });
};

const getResourceName = (type: string, id: string) => {
  if (type === "operator")
    return (
      resources.value.operators.find((o) => o.id === id)?.name ||
      t("planning.unknown_operator")
    );
  if (type === "workcenter")
    return (
      resources.value.workcenters.find((w) => w.id === id)?.name ||
      t("planning.unknown_workcenter")
    );
  if (type === "job")
    return (
      resources.value.jobs.find((j) => j.id === id)?.name ||
      t("planning.unknown_job")
    );
  return id;
};

// ACTIONS

const moveJob = async (entryId: string, direction: "up" | "down") => {
  const entry = planning.value.find((e) => e.id === entryId);
  if (!entry) return;

  // Find siblings (Jobs under same Shift+Op+WC)
  const siblings = planning.value
    .filter(
      (e) =>
        e.shift_id === entry.shift_id &&
        e.operator_id === entry.operator_id &&
        e.workcenter_id === entry.workcenter_id &&
        e.job_id
    )
    .sort((a, b) => a.order - b.order);

  const index = siblings.findIndex((e) => e.id === entryId);
  if (index === -1) return;

  if (direction === "up" && index > 0) {
    // Swap with previous
    const temp = siblings[index];
    siblings[index] = siblings[index - 1];
    siblings[index - 1] = temp;
  } else if (direction === "down" && index < siblings.length - 1) {
    // Swap with next
    const temp = siblings[index];
    siblings[index] = siblings[index + 1];
    siblings[index + 1] = temp;
  } else {
    return;
  }

  // Re-assign order based on new array position
  for (let i = 0; i < siblings.length; i++) {
    // We update blindly to ensure consistency, but only if order changed
    if (siblings[i].order !== i + 1) {
      const updated = { ...siblings[i], order: i + 1 };
      await planningStore.updateEntry(siblings[i].id, updated as any);
    }
  }
};

const removeJob = async (entryId: string) => {
  // We don't delete the entry, we just remove the job_id, so the WC/Op remain.
  const entry = planning.value.find((e) => e.id === entryId);
  if (entry) {
    // Create request object from entry, setting JobID to null/empty
    const updated = {
      customer_id: entry.customer_id,
      shopfloor_id: entry.shopfloor_id,
      shift_id: entry.shift_id,
      workcenter_id: entry.workcenter_id,
      operator_id: entry.operator_id,
      job_id: "", // Clear Job
      date: entry.date,
      order: entry.order,
      is_completed: entry.is_completed,
    };
    await planningStore.updateEntry(entryId, updated as any);
  }
};

const onDropOnShift = async (shiftId: string) => {
  if (!draggedItem.value) return;
  const { type, item, origin } = draggedItem.value;

  if (type === "operator") {
    if (origin === "planning") {
      // Move logic if needed
    } else {
      // Prevent duplicate empty entries for same Op on same Shift
      const exists = planning.value.some(
        (e) =>
          e.shift_id === shiftId &&
          e.operator_id === item.id &&
          !e.workcenter_id
      );
      if (exists) {
        toast.add({
          severity: "info",
          summary: t("planning.already_assigned"),
          detail: t("planning.operator_already_in_shift"),
          life: 2000,
        });
        draggedItem.value = null;
        return;
      }

      try {
        await planningStore.createEntry({
          customer_id: authStore.user?.customer_id || "",
          shopfloor_id: selectedShopfloor.value!,
          shift_id: shiftId,
          operator_id: item.id,
          workcenter_id: "",
          job_id: "",
          date: formatDateForApi(selectedDate.value),
          order: 0,
          is_completed: false,
        });
        toast.add({
          severity: "success",
          summary: t("planning.assigned"),
          detail: t("planning.assigned"),
          life: 2000,
        });
      } catch (e) {
        toast.add({
          severity: "error",
          summary: t("common.error"),
          detail: t("planning.cant_assign"),
          life: 3000,
        });
      }
    }
  }
  draggedItem.value = null;
};

const onDropOnOperator = async (
  shiftId: string,
  opId: string,
  currentWcCount: number = 0
) => {
  if (!draggedItem.value) return;
  const { type, item, origin } = draggedItem.value;

  if (type === "workcenter") {
    if (origin === "planning") {
      const original = planning.value.find((e) => e.id === item.entry_id);
      if (original) {
        const updated = {
          ...original,
          operator_id: opId,
          shift_id: shiftId,
          order: currentWcCount + 1,
        };
        await planningStore.updateEntry(item.entry_id, updated as any);
        toast.add({
          severity: "success",
          summary: t("planning.moved"),
          detail: t("planning.workcenter_reordered"),
          life: 2000,
        });
      }
    } else {
      // Check if there is an "Empty" entry for this operator (no WC, no Job) to consume
      const emptyParent = planning.value.find(
        (e) =>
          e.shift_id === shiftId &&
          e.operator_id === opId &&
          !e.workcenter_id &&
          !e.job_id
      );

      if (emptyParent) {
        // Update the existing empty operator entry
        const updated = {
          ...emptyParent,
          workcenter_id: item.id,
          order: currentWcCount + 1,
        };
        // Map to request format
        const req = {
          customer_id: updated.customer_id,
          shopfloor_id: updated.shopfloor_id,
          shift_id: updated.shift_id,
          workcenter_id: item.id,
          operator_id: updated.operator_id,
          job_id: "",
          date: updated.date,
          order: currentWcCount + 1,
          is_completed: updated.is_completed,
        };
        await planningStore.updateEntry(emptyParent.id, req as any);
        toast.add({
          severity: "success",
          summary: t("planning.assigned"),
          detail: t("planning.workcenter_assigned"),
          life: 2000,
        });
      } else {
        // Create New
        try {
          await planningStore.createEntry({
            customer_id: authStore.user?.customer_id || "",
            shopfloor_id: selectedShopfloor.value!,
            shift_id: shiftId,
            operator_id: opId,
            workcenter_id: item.id,
            job_id: "",
            date: formatDateForApi(selectedDate.value),
            order: currentWcCount + 1,
            is_completed: false,
          });
          toast.add({
            severity: "success",
            summary: t("planning.assigned"),
            detail: t("planning.workcenter_assigned"),
            life: 2000,
          });
        } catch (e) {
          toast.add({
            severity: "error",
            summary: t("common.error"),
            detail: t("planning.error_assigning_wc"),
            life: 3000,
          });
        }
      }
    }
  }
  draggedItem.value = null;
};

const onDropOnWorkcenter = async (
  shiftId: string,
  opId: string,
  wcId: string,
  currentJobsCount: number
) => {
  if (!draggedItem.value) return;
  const { type, item, origin } = draggedItem.value;

  if (type === "job") {
    if (origin === "planning") {
      const original = planning.value.find((e) => e.id === item.entry_id);
      if (original) {
        const updated = {
          ...original,
          workcenter_id: wcId,
          operator_id: opId,
          shift_id: shiftId,
          order: currentJobsCount + 1,
        };
        await planningStore.updateEntry(item.entry_id, updated as any);
        toast.add({
          severity: "success",
          summary: "Mogut",
          detail: "Job reordenat",
          life: 2000,
        });
      }
    } else {
      // Check for "Empty" WC entry (has WC but no Job)
      const emptyParent = planning.value.find(
        (e) =>
          e.shift_id === shiftId &&
          e.operator_id === opId &&
          e.workcenter_id === wcId &&
          !e.job_id
      );

      if (emptyParent) {
        // Update existing WC entry to add Job
        const req = {
          customer_id: emptyParent.customer_id,
          shopfloor_id: emptyParent.shopfloor_id,
          shift_id: emptyParent.shift_id,
          workcenter_id: emptyParent.workcenter_id,
          operator_id: emptyParent.operator_id,
          job_id: item.id,
          date: emptyParent.date,
          order: currentJobsCount + 1,
          is_completed: emptyParent.is_completed,
        };
        await planningStore.updateEntry(emptyParent.id, req as any);
        toast.add({
          severity: "success",
          summary: t("planning.assigned"),
          detail: t("planning.job_assigned"),
          life: 2000,
        });
      } else {
        try {
          await planningStore.createEntry({
            customer_id: authStore.user?.customer_id || "",
            shopfloor_id: selectedShopfloor.value!,
            shift_id: shiftId,
            operator_id: opId,
            workcenter_id: wcId,
            job_id: item.id,
            date: formatDateForApi(selectedDate.value),
            order: currentJobsCount + 1,
            is_completed: false,
          });
          toast.add({
            severity: "success",
            summary: t("planning.assigned"),
            detail: t("planning.job_assigned"),
            life: 2000,
          });
        } catch (e) {
          toast.add({
            severity: "error",
            summary: t("common.error"),
            detail: t("planning.error_assigning_job"),
            life: 3000,
          });
        }
      }
    }
  }
};
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-4rem)] bg-surface-50">
    <!-- Top Bar -->

    <div
      class="h-16 bg-white border-b border-surface-200 px-6 flex items-center justify-between shadow-sm z-10"
    >
      <div class="flex items-center gap-4">
        <h1 class="text-xl font-bold text-surface-900">
          {{ t("planning.title") }}
        </h1>
        <div class="w-px h-6 bg-surface-200"></div>

        <Select
          v-model="selectedShopfloor"
          :options="shopfloors"
          optionLabel="name"
          optionValue="id"
          :placeholder="t('shopfloors.detail')"
          class="w-64"
        />

        <DatePicker
          v-model="selectedDate"
          showIcon
          showButtonBar
          class="w-48"
          dateFormat="dd/mm/yy"
        />
      </div>

      <div class="flex items-center gap-2">
        <Button
          v-if="unsavedChanges"
          :label="t('planning.save_changes')"
          icon="pi pi-save"
          severity="success"
          :loading="loading"
          @click="savePlanning"
        />
        <Button
          icon="pi pi-refresh"
          text
          rounded
          @click="
            selectedShopfloor &&
              planningStore.loadPlanning(selectedShopfloor, selectedDate)
          "
        />
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex flex-1 overflow-hidden" v-if="selectedShopfloor">
      <!-- Resources Palette (Left Sidebar) -->
      <div
        class="w-72 bg-white border-r border-surface-200 flex flex-col shadow-sm z-0"
      >
        <div
          class="p-4 border-b border-surface-100 font-semibold text-surface-700 bg-surface-50"
        >
          {{ t("planning.resources") }}
        </div>
        <div class="flex-1 overflow-y-auto p-4 flex flex-col gap-6">
          <!-- Operators -->
          <div>
            <h3 class="text-xs font-bold text-surface-500 uppercase mb-2">
              {{ t("operators.title") }}
            </h3>
            <div class="flex flex-col gap-2">
              <div
                v-for="op in resources.operators"
                :key="op.id"
                draggable="true"
                @dragstart="onDragStart($event, 'operator', op)"
                class="p-2 bg-white border border-surface-200 rounded-md shadow-sm hover:shadow-md cursor-grab active:cursor-grabbing flex items-center gap-2"
              >
                <div
                  class="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-xs font-bold"
                >
                  {{ op.name.charAt(0) }}
                </div>
                <span class="text-sm font-medium">{{ op.name }}</span>
              </div>
            </div>
          </div>

          <!-- Workcenters -->
          <div>
            <h3 class="text-xs font-bold text-surface-500 uppercase mb-2">
              {{ t("workcenters.title") }}
            </h3>
            <div class="flex flex-col gap-2">
              <div
                v-for="wc in resources.workcenters"
                :key="wc.id"
                draggable="true"
                @dragstart="onDragStart($event, 'workcenter', wc)"
                class="p-2 bg-white border border-surface-200 rounded-md shadow-sm hover:shadow-md cursor-grab active:cursor-grabbing flex items-center gap-2"
              >
                <i class="pi pi-cog text-slate-500"></i>
                <span class="text-sm font-medium">{{ wc.name }}</span>
              </div>
            </div>
          </div>

          <!-- Jobs -->
          <div>
            <h3 class="text-xs font-bold text-surface-500 uppercase mb-2">
              {{ t("jobs.title") }}
            </h3>
            <div class="flex flex-col gap-2">
              <div
                v-for="job in resources.jobs"
                :key="job.id"
                draggable="true"
                @dragstart="onDragStart($event, 'job', job)"
                class="p-2 bg-white border border-surface-200 rounded-md shadow-sm hover:shadow-md cursor-grab active:cursor-grabbing flex items-center gap-2"
              >
                <div class="w-2 h-full bg-orange-400 rounded-l-md"></div>
                <div class="flex flex-col overflow-hidden">
                  <span class="text-sm font-bold truncate">{{
                    job.job_code
                  }}</span>
                  <span class="text-xs text-surface-500 truncate">{{
                    job.product_code
                  }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Planning Board (Board) -->
      <div class="flex-1 overflow-x-auto bg-surface-100 p-6">
        <div class="flex gap-6 h-full min-w-max">
          <!-- Shift Columns -->
          <div
            v-for="shift in resources.shifts"
            :key="shift.id"
            class="w-96 bg-surface-50 rounded-xl border border-surface-200 flex flex-col shadow-sm h-full"
            @dragover="onDragOver"
            @drop="onDropOnShift(shift.id)"
          >
            <!-- Shift Header -->
            <div
              class="p-3 border-b border-surface-200 bg-white rounded-t-xl flex justify-between items-center sticky top-0 z-10"
              :style="{ borderTop: '4px solid ' + (shift.color || '#ccc') }"
            >
              <span class="font-bold text-surface-800">{{ shift.name }}</span>
              <span
                class="text-xs text-surface-500 bg-surface-100 px-2 py-1 rounded"
                >{{ formatTime(shift.start_time) }} -
                {{ formatTime(shift.end_time) }}</span
              >
            </div>

            <!-- Shift Content -->
            <div class="flex-1 p-3 flex flex-col gap-4 overflow-y-auto">
              <!-- Operator Groups -->
              <div
                v-for="op in getShiftPlanning(shift.id)"
                :key="op.id"
                class="bg-white rounded-lg border border-blue-200 shadow-sm"
                @dragover.stop="onDragOver"
                @drop.stop="
                  onDropOnOperator(shift.id, op.id, op.workcenters.length)
                "
              >
                <!-- Operator Header -->
                <div
                  class="p-3 bg-blue-50 border-b border-blue-100 rounded-t-lg flex justify-between items-center"
                >
                  <div class="flex items-center gap-2">
                    <div
                      class="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-bold"
                    >
                      {{ op.name.charAt(0) }}
                    </div>
                    <span class="font-bold text-blue-900">{{ op.name }}</span>
                  </div>
                  <Button
                    icon="pi pi-times"
                    text
                    rounded
                    severity="secondary"
                    size="small"
                    @click="planningStore.deleteEntry(op.entry_id)"
                  />
                </div>

                <!-- Workcenters inside Operator -->
                <div class="p-2 flex flex-col gap-2 bg-blue-50/30">
                  <div
                    v-for="wc in op.workcenters"
                    :key="wc.id"
                    class="bg-white rounded border border-surface-200 shadow-sm cursor-grab active:cursor-grabbing"
                    draggable="true"
                    @dragstart.stop="
                      onDragStartFromBoard($event, 'workcenter', wc)
                    "
                    @dragover.stop="onDragOver"
                    @drop.stop="
                      onDropOnWorkcenter(shift.id, op.id, wc.id, wc.jobs.length)
                    "
                  >
                    <!-- WC Header -->
                    <div
                      class="p-2 border-b border-surface-100 flex justify-between items-center bg-surface-50"
                    >
                      <div class="flex items-center gap-2 text-sm">
                        <i class="pi pi-cog text-surface-500"></i>
                        <span class="font-semibold">{{ wc.name }}</span>
                      </div>
                    </div>

                    <!-- Jobs List -->
                    <div class="p-2 flex flex-col gap-1 min-h-[40px]">
                      <div
                        v-for="(job, index) in wc.jobs"
                        :key="job.id"
                        class="p-2 bg-orange-50 border border-orange-200 rounded text-xs text-orange-900 flex justify-between items-start gap-2 cursor-pointer hover:bg-orange-100 transition-colors"
                        draggable="true"
                        @dragstart.stop="
                          onDragStartFromBoard($event, 'job', job)
                        "
                      >
                        <div class="flex flex-col w-full overflow-hidden">
                          <div class="flex justify-between items-center">
                            <span
                              class="font-bold truncate"
                              :title="job.job_code"
                              >{{ job.job_code }}</span
                            >
                          </div>
                          <span
                            class="text-[0.65rem] text-orange-700 truncate"
                            :title="job.product_code"
                            >{{ job.product_code }}</span
                          >
                          <span
                            class="text-[0.65rem] text-orange-600 mt-0.5 flex items-center gap-1"
                          >
                            <i class="pi pi-stopwatch text-[0.6rem]"></i>
                            {{ job.estimated_duration }}'
                          </span>
                        </div>
                        <div class="flex flex-col gap-1 items-center">
                          <Button
                            v-if="index > 0"
                            icon="pi pi-angle-up"
                            class="!w-4 !h-4 !p-0"
                            text
                            rounded
                            size="small"
                            @click.stop="moveJob(job.entry_id, 'up')"
                          />
                          <Button
                            v-if="index < wc.jobs.length - 1"
                            icon="pi pi-angle-down"
                            class="!w-4 !h-4 !p-0"
                            text
                            rounded
                            size="small"
                            @click.stop="moveJob(job.entry_id, 'down')"
                          />
                          <Button
                            icon="pi pi-times"
                            text
                            rounded
                            size="small"
                            class="!w-4 !h-4 !p-0 text-red-400 hover:text-red-700"
                            @click.stop="removeJob(job.entry_id)"
                          />
                        </div>
                      </div>
                      <div
                        v-if="wc.jobs.length === 0"
                        class="text-xs text-surface-400 text-center italic py-2"
                      >
                        {{ t("planning.drag_jobs") }}
                      </div>
                    </div>
                  </div>

                  <div
                    class="h-12 border-2 border-dashed border-blue-200 rounded flex items-center justify-center text-blue-300 text-xs font-medium"
                  >
                    {{ t("planning.add_workcenter") }}
                  </div>
                </div>
              </div>

              <!-- Drop Zone Placeholder if empty -->
              <div
                v-if="getShiftPlanning(shift.id).length === 0"
                class="h-32 border-2 border-dashed border-surface-200 rounded-lg flex items-center justify-center text-surface-400 text-sm italic"
              >
                {{ t("planning.drag_operators") }}
              </div>
            </div>
          </div>

          <!-- Empty State if no shifts -->
          <div
            v-if="resources.shifts.length === 0"
            class="flex items-center justify-center w-full h-full text-surface-400"
          >
            {{ t("planning.no_shifts") }}
          </div>
        </div>
      </div>
    </div>
    <div
      v-else
      class="flex flex-1 items-center justify-center text-surface-500 flex-col gap-4"
    >
      <i class="pi pi-arrow-up text-4xl animate-bounce"></i>
      <p class="text-lg">Selecciona una planta per començar a planificar</p>
    </div>
  </div>
</template>
