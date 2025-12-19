<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth.store";
import { useToast } from "primevue/usetoast";
import Card from "primevue/card";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const email = ref("");
const password = ref("");
const loading = ref(false);
const errors = ref({
  email: "",
  password: "",
});

const authStore = useAuthStore();
const router = useRouter();
const toast = useToast();

const validateEmail = (email: string) => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

const handleLogin = async () => {
  errors.value = { email: "", password: "" };

  if (!email.value) {
    errors.value.email = t("auth.email_required");
  } else if (!validateEmail(email.value)) {
    errors.value.email = t("auth.email_invalid");
  }

  if (!password.value) {
    errors.value.password = t("auth.password_required");
  }

  if (errors.value.email || errors.value.password) {
    return;
  }

  loading.value = true;
  try {
    await authStore.login(email.value, password.value);
    router.push("/");
  } catch (error: any) {
    if (
      error.response &&
      (error.response.status === 401 || error.response.status === 400)
    ) {
      toast.add({
        severity: "error",
        summary: t("auth.login_error_title"),
        detail: t("auth.login_error_credentials"),
        life: 3000,
      });
    } else {
      toast.add({
        severity: "error",
        summary: t("auth.login_error_title"),
        detail: t("auth.login_error_server"),
        life: 3000,
      });
    }
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="login-container">
    <Card class="login-card">
      <template #title>
        <div class="text-center">{{ t("auth.login") }}</div>
      </template>
      <template #content>
        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="email">{{ t("auth.email") }}</label>
            <InputText
              id="email"
              v-model="email"
              :class="{ 'p-invalid': errors.email }"
              :placeholder="t('auth.enter_email')"
            />
            <small v-if="errors.email" class="error-text">{{
              errors.email
            }}</small>
          </div>

          <div class="form-group">
            <label for="password">{{ t("auth.password") }}</label>
            <Password
              id="password"
              v-model="password"
              :feedback="false"
              toggleMask
              :class="{ 'p-invalid': errors.password }"
              :placeholder="t('auth.enter_password')"
              inputClass="w-full"
            />
            <small v-if="errors.password" class="error-text">{{
              errors.password
            }}</small>
          </div>

          <Button
            :label="t('auth.login')"
            type="submit"
            :loading="loading"
            class="submit-btn"
          />
        </form>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  font-family: "Inter", sans-serif;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.text-center {
  text-align: center;
  font-weight: 600;
  color: #333;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding-top: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-weight: 500;
  color: #555;
  font-size: 0.9rem;
}

.error-text {
  color: #dc2626;
  font-size: 0.8rem;
}

.submit-btn {
  margin-top: 0.5rem;
  width: 100%;
}

:deep(.p-password) {
  width: 100%;
}
:deep(.p-password-input) {
  width: 100%;
}
</style>
