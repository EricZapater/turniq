import api from "./http";

export interface Shopfloor {
  id: string;
  customer_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface ShopfloorRequest {
  customer_id?: string;
  name: string;
}

export interface ShopfloorListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  customer_id?: string;
}

// Response wrapper based on actual backend response
export interface ShopfloorListResponse {
  data: Shopfloor[];
  message: string;
}

// Response wrapper for single shopfloor
export interface ShopfloorResponse {
  data: Shopfloor;
  message: string;
}

export const shopfloorsApi = {
  list: async (
    params?: ShopfloorListParams
  ): Promise<ShopfloorListResponse> => {
    const response = await api.get<ShopfloorListResponse>("/api/shopfloors", {
      params,
    });
    return response.data;
  },
  get: async (id: string): Promise<ShopfloorResponse> => {
    const response = await api.get<ShopfloorResponse>(`/api/shopfloors/${id}`);
    return response.data;
  },
  create: async (data: ShopfloorRequest): Promise<ShopfloorResponse> => {
    const response = await api.post<ShopfloorResponse>("/api/shopfloors", data);
    return response.data;
  },
  update: async (
    id: string,
    data: ShopfloorRequest
  ): Promise<ShopfloorResponse> => {
    const response = await api.put<ShopfloorResponse>(
      `/api/shopfloors/${id}`,
      data
    );
    return response.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/shopfloors/${id}`);
  },
};
