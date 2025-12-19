<script setup lang="ts">
import { useAuthStore } from "../stores/auth.store";
import CustomerDetailView from "./customers/CustomerDetailView.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const authStore = useAuthStore();
</script>

<template>
  <div class="h-full">
    <CustomerDetailView
      v-if="authStore.user?.customer_id"
      :customerIdProp="authStore.user.customer_id"
      :forceViewMode="true"
    />
    <div v-else class="p-6 text-center text-surface-500">
      {{ t("auth.no_customer_info") }}
    </div>
  </div>
</template>
