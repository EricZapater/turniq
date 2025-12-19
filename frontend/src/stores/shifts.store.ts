import { defineStore } from "pinia";
import { ref } from "vue";
import {
  shiftsApi,
  type Shift,
  type ShiftRequest,
  type ShiftListParams,
} from "../api/shifts.api";

export const useShiftsStore = defineStore("shifts", () => {
  const shifts = ref<Shift[]>([]);
  const total = ref<number>(0);
  const currentShift = ref<Shift | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchShifts(params: ShiftListParams) {
    loading.value = true;
    error.value = null;
    try {
      const response = await shiftsApi.list(params);
      shifts.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching shifts", err);
      const msg =
        err.response?.data?.error || err.message || "Error carregant torns";
      error.value = msg;
    } finally {
      loading.value = false;
    }
  }

  async function fetchShiftsByShopfloor(
    shopfloorId: string,
    params?: ShiftListParams
  ) {
    loading.value = true;
    error.value = null;
    shifts.value = []; // Clear previous
    try {
      const response = await shiftsApi.listByShopfloor(shopfloorId, params);
      shifts.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching shifts", err);
      const msg =
        err.response?.data?.error || err.message || "Error carregant torns";
      error.value = msg;
    } finally {
      loading.value = false;
    }
  }

  async function fetchShift(id: string) {
    loading.value = true;
    error.value = null;
    currentShift.value = null;
    try {
      const response = await shiftsApi.get(id);
      currentShift.value = response.data;
    } catch (err: any) {
      console.error("Error fetching shift", err);
      const msg =
        err.response?.data?.error || err.message || "Error carregant torn";
      error.value = msg;
    } finally {
      loading.value = false;
    }
  }

  async function createShift(data: ShiftRequest) {
    loading.value = true;
    error.value = null;
    try {
      await shiftsApi.create(data);
    } catch (err: any) {
      console.error("Error creating shift", err);
      const msg =
        err.response?.data?.error || err.message || "Error creant torn";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateShift(id: string, data: ShiftRequest) {
    loading.value = true;
    error.value = null;
    try {
      await shiftsApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating shift", err);
      const msg =
        err.response?.data?.error || err.message || "Error actualitzant torn";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function deleteShift(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await shiftsApi.delete(id);
    } catch (err: any) {
      console.error("Error deleting shift", err);
      const msg =
        err.response?.data?.error || err.message || "Error eliminant torn";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  return {
    shifts,
    total,
    currentShift,
    loading,
    error,
    fetchShifts,
    fetchShiftsByShopfloor,
    fetchShift,
    createShift,
    updateShift,
    deleteShift,
  };
});
