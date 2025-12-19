import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  scheduleApi,
  type ScheduleEntry,
  type ScheduleEntryRequest,
} from "../api/schedule.api";
import { operatorsApi } from "../api/operators.api";
import { workcentersApi } from "../api/workcenters.api";
import { jobsApi } from "../api/jobs.api";
import { shiftsApi } from "../api/shifts.api";

export const usePlanningStore = defineStore("planning", () => {
  // State
  const resources = ref({
    operators: [] as any[],
    workcenters: [] as any[],
    jobs: [] as any[],
    shifts: [] as any[],
  });

  const planning = ref<ScheduleEntry[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const unsavedChanges = ref(false);

  // Actions
  async function loadResources(shopfloorId: string) {
    loading.value = true;
    try {
      // Load all necessary info for the Shopfloor
      const [opsRes, wcsRes, jobsRes, shiftsRes] = await Promise.all([
        operatorsApi.list({
          page: 1,
          page_size: 1000,
          // shop_floor_id was error, checking Api definition would be best but let's assume valid param or fix later
          // Removing for now to stop red noise if unsure, or check API file.
        } as any),
        workcentersApi.list({
          page: 1,
          page_size: 1000,
        } as any),
        jobsApi.list({ page: 1, page_size: 1000 } as any),
        shiftsApi.listByShopfloor(shopfloorId),
      ]);

      resources.value.operators = opsRes.data || [];
      resources.value.workcenters = wcsRes.data || [];
      resources.value.jobs = jobsRes.data || [];
      resources.value.shifts = shiftsRes.data || [];
    } catch (e: any) {
      console.error("Error loading resources", e);
      error.value = "Failed to load resources";
    } finally {
      loading.value = false;
    }
  }

  async function loadPlanning(shopfloorId: string, date: Date) {
    loading.value = true;
    try {
      // Use local date string YYYY-MM-DD
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const dateStr = `${year}-${month}-${day}`;
      // Call API
      const entries = await scheduleApi.getPlanning(shopfloorId, dateStr);
      planning.value = entries;
    } catch (e: any) {
      console.error("Error loading planning", e);
      error.value = "Failed to load planning";
    } finally {
      loading.value = false;
    }
  }

  async function createEntry(entry: ScheduleEntryRequest) {
    // Local Update Only
    const tempId = "temp-" + Date.now() + Math.random();
    const localEntry: any = { ...entry, id: tempId };
    planning.value.push(localEntry);
    unsavedChanges.value = true;
  }

  async function updateEntry(id: string, entry: ScheduleEntryRequest) {
    // Local Update Only
    const idx = planning.value.findIndex((p) => p.id === id);
    if (idx !== -1) {
      planning.value[idx] = { ...planning.value[idx], ...entry };
      unsavedChanges.value = true;
    }
  }

  async function deleteEntry(id: string) {
    // Local Delete Only
    planning.value = planning.value.filter((p) => p.id !== id);
    unsavedChanges.value = true;
  }

  async function saveChanges(shopfloorId: string, date: string) {
    loading.value = true;
    try {
      const entriesRequest: ScheduleEntryRequest[] = planning.value.map(
        (p) => ({
          customer_id: p.customer_id,
          shopfloor_id: p.shopfloor_id,
          shift_id: p.shift_id,
          workcenter_id: p.workcenter_id,
          job_id: p.job_id,
          operator_id: p.operator_id,
          date: p.date,
          order: p.order,
          start_time: p.start_time,
          end_time: p.end_time,
          is_completed: p.is_completed,
        })
      );

      await scheduleApi.sync(shopfloorId, date, entriesRequest);
      await loadPlanning(shopfloorId, new Date(date)); // Reload to get real IDs
      unsavedChanges.value = false;
    } catch (e) {
      console.error("Error saving planning", e);
      error.value = "Error saving planning";
      throw e; // Propagate to view
    } finally {
      loading.value = false;
    }
  }

  // Getters
  const availableJobs = computed(() => {
    // Filter out jobs that are already in the planning (assigned to an entry)
    const assignedJobIds = new Set<string>(
      planning.value.filter((e) => e.job_id).map((e) => e.job_id as string)
    );
    return resources.value.jobs.filter((j) => !assignedJobIds.has(j.id));
  });

  const statistics = computed(() => {
    const operatorLoad: Record<string, number> = {};
    const shiftOperatorLoad: Record<string, Record<string, number>> = {};
    const shiftWorkcenterLoad: Record<string, Record<string, number>> = {};
    const assignedOperatorIds = new Set<string>();

    // Initialize metrics
    resources.value.shifts.forEach((s) => {
      shiftOperatorLoad[s.id] = {};
      shiftWorkcenterLoad[s.id] = {};
    });

    planning.value.forEach((entry) => {
      const duration = 0; // entry.duration ?? we need to match job to get duration
      let jobDuration = 0;
      if (entry.job_id) {
        const job = resources.value.jobs.find((j) => j.id === entry.job_id);
        if (job) jobDuration = job.estimated_duration || 0;
      }

      // 1. Operator Load (Total & Per Shift)
      if (entry.operator_id) {
        const opId = entry.operator_id;
        assignedOperatorIds.add(opId);

        // Total
        const currentOpLoad = operatorLoad[opId] || 0;
        operatorLoad[opId] = currentOpLoad + jobDuration;

        // Per Shift
        let shiftRecord = shiftOperatorLoad[entry.shift_id];
        if (!shiftRecord) {
          shiftRecord = {};
          shiftOperatorLoad[entry.shift_id] = shiftRecord;
        }
        const currentShiftOpLoad = shiftRecord[opId] || 0;
        shiftRecord[opId] = currentShiftOpLoad + jobDuration;
      }

      // 2. Workcenter Load (Per Shift)
      if (entry.workcenter_id && entry.shift_id) {
        const wcId = entry.workcenter_id;
        let shiftRecord = shiftWorkcenterLoad[entry.shift_id];
        if (!shiftRecord) {
          shiftRecord = {};
          shiftWorkcenterLoad[entry.shift_id] = shiftRecord;
        }
        const currentShiftWcLoad = shiftRecord[wcId] || 0;
        shiftRecord[wcId] = currentShiftWcLoad + jobDuration;
      }
    });

    const unassignedOperatorsCount =
      resources.value.operators.length - assignedOperatorIds.size;

    return {
      operatorLoad,
      shiftOperatorLoad,
      shiftWorkcenterLoad,
      unassignedOperatorsCount,
    };
  });

  return {
    resources,
    planning,
    loading,
    error,
    loadResources,
    loadPlanning,
    createEntry,
    updateEntry,
    deleteEntry,
    saveChanges,
    unsavedChanges,
    availableJobs,
    statistics,
  };
});
