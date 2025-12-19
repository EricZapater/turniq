import { defineStore } from "pinia";
import { ref } from "vue";
import {
  shopfloorsApi,
  type Shopfloor,
  type ShopfloorRequest,
  type ShopfloorListParams,
} from "../api/shopfloors.api";

export const useShopfloorsStore = defineStore("shopfloors", () => {
  const shopfloors = ref<Shopfloor[]>([]);
  const total = ref<number>(0);
  const currentShopfloor = ref<Shopfloor | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  async function fetchShopfloors(params: ShopfloorListParams) {
    loading.value = true;
    error.value = null;
    try {
      const response = await shopfloorsApi.list(params);
      shopfloors.value = response.data;
      total.value = response.data ? response.data.length : 0;
    } catch (err: any) {
      console.error("Error fetching shopfloors", err);
      error.value = err.message || "Error carregant plantes";
    } finally {
      loading.value = false;
    }
  }

  async function fetchShopfloor(id: string) {
    loading.value = true;
    error.value = null;
    currentShopfloor.value = null;
    try {
      const response = await shopfloorsApi.get(id);
      currentShopfloor.value = response.data;
    } catch (err: any) {
      console.error("Error fetching shopfloor", err);
      error.value = err.message || "Error carregant planta";
    } finally {
      loading.value = false;
    }
  }

  async function createShopfloor(data: ShopfloorRequest) {
    loading.value = true;
    error.value = null;
    try {
      await shopfloorsApi.create(data);
    } catch (err: any) {
      console.error("Error creating shopfloor", err);
      const msg =
        err.response?.data?.error || err.message || "Error creant planta";
      error.value = msg;
      throw new Error(msg);
    } finally {
      loading.value = false;
    }
  }

  async function updateShopfloor(id: string, data: ShopfloorRequest) {
    loading.value = true;
    error.value = null;
    try {
      await shopfloorsApi.update(id, data);
    } catch (err: any) {
      console.error("Error updating shopfloor", err);
      error.value = err.message || "Error actualitzant planta";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deleteShopfloor(id: string) {
    loading.value = true;
    error.value = null;
    try {
      await shopfloorsApi.delete(id);
    } catch (err: any) {
      console.error("Error deleting shopfloor", err);
      error.value = err.message || "Error eliminant planta";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    shopfloors,
    total,
    currentShopfloor,
    loading,
    error,
    fetchShopfloors,
    fetchShopfloor,
    createShopfloor,
    updateShopfloor,
    deleteShopfloor,
  };
});
