import { defineStore } from "pinia";
import { ref } from "vue";
import {
  customersApi,
  type Customer,
  type CustomerRequest,
  type CustomerListParams,
} from "../api/customers.api";

export const useCustomersStore = defineStore("customers", () => {
  const customers = ref<Customer[]>([]);
  const total = ref<number>(0);
  const currentCustomer = ref<Customer | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);
  const isForbidden = ref<boolean>(false);

  async function fetchCustomers(params: CustomerListParams) {
    loading.value = true;
    error.value = null;
    isForbidden.value = false;
    try {
      const response = await customersApi.list(params);
      customers.value = response.data;
      // Backend does not return total count yet, so we use the length of the data array or 0 if empty
      // Ideally backend should return total count for server-side pagination
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching customers", err);
      if (err.response && err.response.status === 403) {
        isForbidden.value = true;
      }
      error.value = err.message || "Error carregant clients";
    } finally {
      loading.value = false;
    }
  }

  async function fetchCustomer(id: string) {
    loading.value = true;
    error.value = null;
    currentCustomer.value = null;
    try {
      const response = await customersApi.get(id);
      currentCustomer.value = response.data;
    } catch (err: any) {
      console.error("Error fetching customer", err);
      error.value = err.message || "Error carregant client";
    } finally {
      loading.value = false;
    }
  }

  async function createCustomer(data: CustomerRequest) {
    loading.value = true;
    error.value = null;
    try {
      await customersApi.create(data);
    } catch (err: any) {
      console.error("Error creating customer", err);
      // Extract backend error message if available
      const msg =
        err.response?.data?.error || err.message || "Error creant client";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateCustomer(id: string, data: CustomerRequest) {
    loading.value = true;
    error.value = null;
    try {
      await customersApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating customer", err);
      error.value = err.message || "Error actualitzant client";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    customers,
    total,
    currentCustomer,
    loading,
    error,
    isForbidden,
    fetchCustomers,
    fetchCustomer,
    createCustomer,
    updateCustomer,
  };
});
