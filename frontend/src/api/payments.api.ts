import api from "./http";

export interface Payment {
  id: string;
  customer_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  status: string;
  due_date: string;
  paid_at: string;
  created_at: string;
  updated_at: string;
}

export interface PaymentRequest {
  customer_id: string;
  amount: number;
  currency?: string;
  payment_method?: string;
  status: string;
  due_date?: string;
  paid_at: string;
}

export interface PaymentListResponse {
  data: Payment[];
  message: string;
}

export interface PaymentResponse {
  data: Payment;
  message: string;
}

export const paymentsApi = {
  create: async (data: PaymentRequest): Promise<PaymentResponse> => {
    const response = await api.post<PaymentResponse>("/api/payments", data);
    return response.data;
  },

  delete: async (id: string) => {
    await api.delete(`/api/payments/${id}`);
  },

  findByCustomer: async (customerId: string): Promise<PaymentListResponse> => {
    const response = await api.get<PaymentListResponse>(
      `/api/payments/customer/${customerId}`
    );
    return response.data;
  },

  list: async (params?: any): Promise<PaymentListResponse> => {
    const response = await api.get<PaymentListResponse>("/api/payments", {
      params,
    });
    return response.data;
  },
};
