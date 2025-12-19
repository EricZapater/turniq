import api from "./http";

export interface Workcenter {
  id: string;
  customer_id: string;
  shop_floor_id: string | null;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface WorkcenterRequest {
  customer_id?: string;
  shop_floor_id?: string | null;
  name: string;
  is_active: boolean;
}

export interface WorkcenterListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  customer_id?: string;
}

// Response wrapper based on actual backend response
export interface WorkcenterListResponse {
  data: Workcenter[];
  message: string;
}

// Response wrapper for single workcenter
export interface WorkcenterResponse {
  data: Workcenter;
  message: string;
}

export const workcentersApi = {
  list: async (
    params?: WorkcenterListParams
  ): Promise<WorkcenterListResponse> => {
    const response = await api.get<WorkcenterListResponse>("/api/workcenters", {
      params,
    });
    return response.data;
  },
  get: async (id: string): Promise<WorkcenterResponse> => {
    const response = await api.get<WorkcenterResponse>(
      `/api/workcenters/${id}`
    );
    return response.data;
  },
  create: async (data: WorkcenterRequest): Promise<WorkcenterResponse> => {
    const response = await api.post<WorkcenterResponse>(
      "/api/workcenters",
      data
    );
    return response.data;
  },
  update: async (
    id: string,
    data: WorkcenterRequest
  ): Promise<WorkcenterResponse> => {
    const response = await api.put<WorkcenterResponse>(
      `/api/workcenters/${id}`,
      data
    );
    return response.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/workcenters/${id}`);
  },
};
