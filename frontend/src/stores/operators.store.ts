import { defineStore } from "pinia";
import { ref } from "vue";
import {
  operatorsApi,
  type Operator,
  type OperatorRequest,
  type OperatorListParams,
} from "../api/operators.api";

export const useOperatorsStore = defineStore("operators", () => {
  const operators = ref<Operator[]>([]);
  const total = ref<number>(0);
  const currentOperator = ref<Operator | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchOperators(params: OperatorListParams) {
    loading.value = true;
    error.value = null;
    try {
      const response = await operatorsApi.list(params);
      operators.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching operators", err);
      error.value = err.message || "Error carregant operaris";
    } finally {
      loading.value = false;
    }
  }

  async function fetchOperator(id: string) {
    loading.value = true;
    error.value = null;
    currentOperator.value = null;
    try {
      const response = await operatorsApi.get(id);
      currentOperator.value = response.data;
    } catch (err: any) {
      console.error("Error fetching operator", err);
      error.value = err.message || "Error carregant operari";
    } finally {
      loading.value = false;
    }
  }

  async function createOperator(data: OperatorRequest) {
    loading.value = true;
    error.value = null;
    try {
      await operatorsApi.create(data);
    } catch (err: any) {
      console.error("Error creating operator", err);
      const msg =
        err.response?.data?.error || err.message || "Error creant operari";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateOperator(id: string, data: OperatorRequest) {
    loading.value = true;
    error.value = null;
    try {
      await operatorsApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating operator", err);
      error.value = err.message || "Error actualitzant operari";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deleteOperator(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await operatorsApi.delete(id);
    } catch (err: any) {
      console.error("Error deleting operator", err);
      error.value = err.message || "Error eliminant operari";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    operators,
    total,
    currentOperator,
    loading,
    error,
    fetchOperators,
    fetchOperator,
    createOperator,
    updateOperator,
    deleteOperator,
  };
});
