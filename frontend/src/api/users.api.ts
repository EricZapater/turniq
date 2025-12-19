import api from "./http";
import type { User } from "./auth.api";

// Reusing User interface from auth.api to avoid duplication
export type { User };

// Response structure for users list
// Assuming backend returns { data: User[], message: string } or similar
export interface UserRequest {
  customer_id: string;
  username: string;
  email: string;
  password?: string; // Optional for updates
  is_active: boolean;
  is_admin: boolean;
}

export interface UsersListResponse {
  data: User[];
  message: string;
}

export interface UserResponse {
  data: User;
  message: string;
}

export const usersApi = {
  list: async (): Promise<UsersListResponse> => {
    const response = await api.get<UsersListResponse>("/api/users");
    return response.data;
  },

  listByCustomer: async (customerId: string): Promise<UsersListResponse> => {
    // Assuming backend supports filtering by query param or specific endpoint
    // Given typical pattern, likely /api/users?customer_id=...
    const response = await api.get<UsersListResponse>(
      `/api/users/customer/${customerId}`
    );
    return response.data;
  },

  create: async (data: UserRequest): Promise<UserResponse> => {
    const response = await api.post<UserResponse>("/api/users", data);
    return response.data;
  },

  update: async (
    id: string,
    data: Partial<UserRequest>
  ): Promise<UserResponse> => {
    const response = await api.put<UserResponse>(`/api/users/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/api/users/${id}`);
  },
};
