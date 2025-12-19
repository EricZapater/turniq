import { defineStore } from "pinia";
import { ref } from "vue";
import {
  workcentersApi,
  type Workcenter,
  type WorkcenterRequest,
  type WorkcenterListParams,
} from "../api/workcenters.api";

export const useWorkcentersStore = defineStore("workcenters", () => {
  const workcenters = ref<Workcenter[]>([]);
  const total = ref<number>(0);
  const currentWorkcenter = ref<Workcenter | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchWorkcenters(params: WorkcenterListParams) {
    loading.value = true;
    error.value = null;
    try {
      const response = await workcentersApi.list(params);
      workcenters.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching workcenters", err);
      error.value = err.message || "Error carregant centres de treball";
    } finally {
      loading.value = false;
    }
  }

  async function fetchWorkcenter(id: string) {
    loading.value = true;
    error.value = null;
    currentWorkcenter.value = null;
    try {
      const response = await workcentersApi.get(id);
      currentWorkcenter.value = response.data;
    } catch (err: any) {
      console.error("Error fetching workcenter", err);
      error.value = err.message || "Error carregant centre de treball";
    } finally {
      loading.value = false;
    }
  }

  async function createWorkcenter(data: WorkcenterRequest) {
    loading.value = true;
    error.value = null;
    try {
      await workcentersApi.create(data);
    } catch (err: any) {
      console.error("Error creating workcenter", err);
      const msg =
        err.response?.data?.error ||
        err.message ||
        "Error creant centre de treball";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateWorkcenter(id: string, data: WorkcenterRequest) {
    loading.value = true;
    error.value = null;
    try {
      await workcentersApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating workcenter", err);
      error.value = err.message || "Error actualitzant centre de treball";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deleteWorkcenter(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await workcentersApi.delete(id);
    } catch (err: any) {
      console.error("Error deleting workcenter", err);
      error.value = err.message || "Error eliminant centre de treball";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    workcenters,
    total,
    currentWorkcenter,
    loading,
    error,
    fetchWorkcenters,
    fetchWorkcenter,
    createWorkcenter,
    updateWorkcenter,
    deleteWorkcenter,
  };
});
