export interface ScheduleEntry {
  id: string;
  customer_id: string;
  shopfloor_id: string;
  shift_id: string;
  workcenter_id: string; // Required in backend model
  job_id?: string | null;
  operator_id?: string | null;
  date: string; // ISO Date "YYYY-MM-DD"
  order: number;
  start_time?: string | null;
  end_time?: string | null;
  is_completed: boolean;
}

export interface ScheduleEntryRequest {
  customer_id: string;
  shopfloor_id: string;
  shift_id: string;
  workcenter_id: string;
  job_id?: string | null;
  operator_id?: string | null;
  date: string;
  order: number;
  start_time?: string | null;
  end_time?: string | null;
  is_completed: boolean;
}
import api from "./http";

export const scheduleApi = {
  create: async (data: ScheduleEntryRequest): Promise<ScheduleEntry> => {
    // Backend returns wrapped { message, data } usually, but let's check handler
    // Handler: c.JSON(http.StatusOK, gin.H{"message": "...", "data": response})
    const response = await api.post<any>("/api/schedule-entries", data);
    return response.data.data;
  },
  update: async (
    id: string,
    data: ScheduleEntryRequest
  ): Promise<ScheduleEntry> => {
    const response = await api.put<any>(`/api/schedule-entries/${id}`, data);
    return response.data.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/schedule-entries/${id}`);
  },
  getPlanning: async (
    shopfloorId: string,
    date: string
  ): Promise<ScheduleEntry[]> => {
    // Handler: c.JSON(http.StatusOK, gin.H{"data": response})
    const response = await api.get<any>("/api/schedule-entries", {
      params: { shopfloor_id: shopfloorId, date },
    });
    return response.data.data || [];
  },
  sync: async (
    shopfloorId: string,
    date: string,
    entries: ScheduleEntryRequest[]
  ): Promise<void> => {
    await api.post("/api/schedule-entries/sync", {
      shopfloor_id: shopfloorId,
      date,
      entries,
    });
  },
  getOperatorPlanning: async (
    operatorId: string,
    date: string
  ): Promise<ScheduleEntry[]> => {
    const response = await api.get<any>("/api/schedule-entries", {
      params: { operator_id: operatorId, date },
    });
    return response.data.data || [];
  },
  list: async (params?: any): Promise<ScheduleEntry[]> => {
    const response = await api.get<any>("/api/schedule-entries", { params });
    return response.data.data || [];
  },
};
