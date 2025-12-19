import { defineStore } from "pinia";
import { ref } from "vue";
import {
  paymentsApi,
  type Payment,
  type PaymentRequest,
} from "../api/payments.api";

export const usePaymentsStore = defineStore("payments", () => {
  const payments = ref<Payment[]>([]);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchByCustomer(customerId: string) {
    loading.value = true;
    error.value = null;
    try {
      const response = await paymentsApi.findByCustomer(customerId);
      payments.value = response.data || [];
    } catch (err: any) {
      console.error("Error fetching payments", err);
      error.value = err.message || "Error carregant pagaments";
    } finally {
      loading.value = false;
    }
  }

  async function createPayment(data: PaymentRequest) {
    loading.value = true;
    error.value = null;
    try {
      await paymentsApi.create(data);
      // Refresh list if needed, or caller handles it.
      // Usually fetchByCustomer is called after.
    } catch (err: any) {
      console.error("Error creating payment", err);
      error.value = err.message || "Error creant pagament";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deletePayment(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await paymentsApi.delete(id);
      // Optimistic update
      payments.value = payments.value.filter((p) => p.id !== id);
    } catch (err: any) {
      console.error("Error deleting payment", err);
      error.value = err.message || "Error eliminant pagament";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    payments,
    loading,
    error,
    fetchByCustomer,
    createPayment,
    deletePayment,
  };
});
