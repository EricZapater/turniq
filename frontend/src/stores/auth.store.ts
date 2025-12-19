import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { authApi, type User } from "../api/auth.api";

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem("token"));

  const isAuthenticated = computed(() => !!token.value);
  const isAdmin = computed(() => !!user.value?.is_admin);

  async function login(email: string, password: string) {
    const response = await authApi.login(email, password);

    token.value = response.token;
    user.value = response.user;

    localStorage.setItem("token", response.token);
    localStorage.setItem("user", JSON.stringify(response.user));
  }

  function logout() {
    token.value = null;
    user.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  }

  function loadFromStorage() {
    const storedToken = localStorage.getItem("token");
    const storedUser = localStorage.getItem("user");

    if (storedToken) {
      token.value = storedToken;
    }
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser);
      } catch (e) {
        console.error("Invalid user data in storage");
        localStorage.removeItem("user");
      }
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    loadFromStorage,
  };
});
