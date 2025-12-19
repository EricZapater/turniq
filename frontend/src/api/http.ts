import axios from "axios";

const api = axios.create({
  // Use runtime configuration (window.env) or fallback to build-time env
  baseURL: (window as any).env?.API_URL || import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Clear auth data
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // Force reload/redirect to login
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
