import api from "./http";

export interface User {
  id: string;
  customer_id: string;
  username: string;
  email: string;
  is_admin: boolean;
  is_active: boolean;
}

export interface LoginResponse {
  token: string;
  expire: string;
  user: User;
}

export const authApi = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>("/auth/login", {
      email,
      password,
    });
    return response.data;
  },
};
