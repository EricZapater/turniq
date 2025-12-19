import api from "./http";

export interface Customer {
  id: string;
  name: string;
  email: string;
  vat_number: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  country: string;
  language: string;
  contact_name: string;
  status: string;
  plan: string;
  billing_cycle: string;
  price: number;
  trial_ends_at: string | null;
  internal_notes: string;
  max_operators: number;
  max_workcenters: number;
  max_shop_floors: number;
  max_users: number;
  max_jobs: number;
  created_at: string;
  updated_at: string;
}

export interface CustomerRequest {
  name: string;
  email: string;
  vat_number: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  country: string;
  language: string;
  contact_name: string;
  status: string;
  plan: string;
  billing_cycle: string;
  price: number;
  trial_ends_at: string | null;
  internal_notes: string;
  max_operators: number;
  max_workcenters: number;
  max_shop_floors: number;
  max_users: number;
  max_jobs: number;
}

export interface CustomerListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
}

// Response wrapper based on actual backend response
export interface CustomerListResponse {
  data: Customer[];
  message: string;
}

// Response wrapper for single customer
export interface CustomerResponse {
  data: Customer;
  message: string;
}

export const customersApi = {
  list: async (params: CustomerListParams): Promise<CustomerListResponse> => {
    // Ensuring query params are handled
    const response = await api.get<CustomerListResponse>("/api/customers", {
      params,
    });
    return response.data;
  },

  get: async (id: string): Promise<CustomerResponse> => {
    const response = await api.get<CustomerResponse>(`/api/customers/${id}`);
    return response.data;
  },

  create: async (data: CustomerRequest): Promise<CustomerResponse> => {
    const response = await api.post<CustomerResponse>("/api/customers", data);
    return response.data;
  },

  update: async (
    id: string,
    data: CustomerRequest
  ): Promise<CustomerResponse> => {
    const response = await api.put<CustomerResponse>(
      `/api/customers/${id}`,
      data
    );
    return response.data;
  },
};
