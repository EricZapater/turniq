import api from "./http";

export interface Shift {
  id: string;
  customer_id: string;
  shopfloor_id?: string;
  name: string;
  color: string;
  start_time: string; // "HT:mm:ss" or ISO depending on backend, backend parse "15:04" so expect string HH:mm
  end_time: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ShiftRequest {
  customer_id: string;
  shopfloor_id?: string;
  name: string;
  color: string;
  start_time: string; // HH:mm
  end_time: string; // HH:mm
  is_active: boolean;
}

export interface ShiftListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  shop_floor_id?: string; // Filter
  customer_id?: string;
}

// Response wrapper based on actual backend response
export interface ShiftListResponse {
  data: Shift[];
  message: string;
}

// Response wrapper for single shift
export interface ShiftResponse {
  data: Shift;
  message: string;
}

export const shiftsApi = {
  list: async (params?: ShiftListParams): Promise<ShiftListResponse> => {
    const response = await api.get<ShiftListResponse>("/api/shifts", {
      params,
    });
    return response.data;
  },
  listByShopfloor: async (
    shopfloorId: string,
    params?: ShiftListParams
  ): Promise<ShiftListResponse> => {
    const response = await api.get<ShiftListResponse>(
      `/api/shifts/shopfloor/${shopfloorId}`,
      {
        params,
      }
    );
    return response.data;
  },
  get: async (id: string): Promise<ShiftResponse> => {
    const response = await api.get<ShiftResponse>(`/api/shifts/${id}`);
    return response.data;
  },
  create: async (data: ShiftRequest): Promise<ShiftResponse> => {
    const response = await api.post<ShiftResponse>("/api/shifts", data);
    return response.data;
  },
  update: async (id: string, data: ShiftRequest): Promise<ShiftResponse> => {
    const response = await api.put<ShiftResponse>(`/api/shifts/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/shifts/${id}`);
  },
};
