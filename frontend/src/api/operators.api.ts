import api from "./http";

export interface Operator {
  id: string;
  shop_floor_id: string;
  customer_id: string;
  code: string;
  name: string;
  surname: string;
  vat_number: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface OperatorRequest {
  shop_floor_id: string;
  customer_id?: string;
  code: string;
  name: string;
  surname: string;
  vat_number: string;
  is_active: boolean;
}

export interface OperatorListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  customer_id?: string;
}

// Response wrapper based on actual backend response
export interface OperatorListResponse {
  data: Operator[];
  message: string;
}

// Response wrapper for single operator
export interface OperatorResponse {
  data: Operator;
  message: string;
}

export const operatorsApi = {
  list: async (params?: OperatorListParams): Promise<OperatorListResponse> => {
    const response = await api.get<OperatorListResponse>("/api/operators", {
      params,
    });
    return response.data;
  },
  get: async (id: string): Promise<OperatorResponse> => {
    const response = await api.get<OperatorResponse>(`/api/operators/${id}`);
    return response.data;
  },
  findByCode: async (code: string): Promise<OperatorResponse> => {
    const response = await api.get<OperatorResponse>(
      `/api/operators/code/${code}`
    );
    return response.data;
  },
  create: async (data: OperatorRequest): Promise<OperatorResponse> => {
    const response = await api.post<OperatorResponse>("/api/operators", data);
    return response.data;
  },
  update: async (
    id: string,
    data: OperatorRequest
  ): Promise<OperatorResponse> => {
    const response = await api.put<OperatorResponse>(
      `/api/operators/${id}`,
      data
    );
    return response.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/operators/${id}`);
  },
};
