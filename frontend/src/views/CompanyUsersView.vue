<script setup lang="ts">
import { useAuthStore } from "../stores/auth.store";
import CustomerUsersList from "./customers/CustomerUsersList.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const authStore = useAuthStore();
</script>

<template>
  <div class="card p-6 bg-white rounded-lg shadow-sm border border-surface-200">
    <div class="mb-6 border-b border-surface-100 pb-4">
      <h1 class="text-2xl font-bold text-surface-900 tracking-tight">
        {{ t("users.title") }}
      </h1>
      <p class="text-surface-500 text-sm mt-1">
        {{ t("users.subtitle") }}
      </p>
    </div>

    <CustomerUsersList
      v-if="authStore.user?.customer_id"
      :customerId="authStore.user.customer_id"
    />
    <div v-else class="p-6 text-center text-surface-500">
      {{ t("auth.no_customer_info") }}
    </div>
  </div>
</template>
