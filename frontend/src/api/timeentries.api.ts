import api from "./http";

export interface TimeEntry {
  id: string;
  customer_id: string;
  shopfloor_id: string;
  operator_id: string;
  workcenter_id?: string;
  check_in: string; // ISO date
  check_out?: string; // ISO date
  created_at: string;
  updated_at: string;
}

export interface TimeEntryRequest {
  id?: string;
  customer_id?: string;
  shopfloor_id?: string;
  operator_id: string;
  workcenter_id?: string;
  check_in?: string;
  check_out?: string;
}

export interface TimeEntryResponse {
  data: TimeEntry;
  message: string;
}

export interface TimeEntryListResponse {
  data: TimeEntry[];
  message: string;
}

export const timeentriesApi = {
  findCurrent: async (operatorId: string): Promise<TimeEntryResponse> => {
    const response = await api.get<TimeEntryResponse>(
      `/api/time-entries/current/${operatorId}`
    );
    return response.data;
  },
  findByOperator: async (
    operatorId: string
  ): Promise<TimeEntryListResponse> => {
    const response = await api.get<TimeEntryListResponse>(
      `/api/time-entries/operator/${operatorId}`
    );
    return response.data;
  },
  create: async (data: TimeEntryRequest): Promise<TimeEntryResponse> => {
    const response = await api.post<TimeEntryResponse>(
      "/api/time-entries",
      data
    );
    return response.data;
  },
  update: async (
    id: string,
    data: TimeEntryRequest
  ): Promise<TimeEntryResponse> => {
    const response = await api.put<TimeEntryResponse>(
      `/api/time-entries/${id}`,
      data
    );
    return response.data;
  },
  list: async (params?: any): Promise<TimeEntryListResponse> => {
    const response = await api.get<TimeEntryListResponse>("/api/time-entries", {
      params,
    });
    return response.data;
  },
};
