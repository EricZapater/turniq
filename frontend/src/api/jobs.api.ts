import api from "./http";

export interface Job {
  id: string;
  customer_id: string;
  shop_floor_id: string;
  workcenter_id: string;
  job_code: string;
  product_code: string;
  description: string;
  estimated_duration: number;
  created_at: string;
  updated_at: string;
}

export interface JobRequest {
  customer_id?: string;
  shop_floor_id: string;
  workcenter_id: string;
  job_code: string;
  product_code: string;
  description: string;
  estimated_duration: number;
}

export interface JobListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_desc?: boolean;
  customer_id?: string;
}

// Response wrapper based on actual backend response
export interface JobListResponse {
  data: Job[];
  message: string;
}

// Response wrapper for single job
export interface JobResponse {
  data: Job;
  message: string;
}

export const jobsApi = {
  list: async (params?: JobListParams): Promise<JobListResponse> => {
    const response = await api.get<JobListResponse>("/api/jobs", { params });
    return response.data;
  },
  get: async (id: string): Promise<JobResponse> => {
    const response = await api.get<JobResponse>(`/api/jobs/${id}`);
    return response.data;
  },
  create: async (data: JobRequest): Promise<JobResponse> => {
    const response = await api.post<JobResponse>("/api/jobs", data);
    return response.data;
  },
  update: async (id: string, data: JobRequest): Promise<JobResponse> => {
    const response = await api.put<JobResponse>(`/api/jobs/${id}`, data);
    return response.data;
  },
  delete: async (id: string) => {
    await api.delete(`/api/jobs/${id}`);
  },
};
