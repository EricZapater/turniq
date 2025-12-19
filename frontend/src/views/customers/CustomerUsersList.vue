<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { usersApi, type User, type UserRequest } from "../../api/users.api";
import { useToast } from "primevue/usetoast";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

// PrimeVue Components
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Button from "primevue/button";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import ToggleSwitch from "primevue/toggleswitch";
import Tag from "primevue/tag";

const props = defineProps<{
  customerId: string;
  readOnly?: boolean;
}>();

const toast = useToast();
const users = ref<User[]>([]);
const loading = ref(false);

// Dialog
const showDialog = ref(false);
const isEdit = ref(false);
const submitting = ref(false);

const form = ref<UserRequest>({
  customer_id: "",
  username: "",
  email: "",
  password: "",
  is_active: true,
  is_admin: false,
});

const currentUserId = ref<string | null>(null);

const fetchUsers = async () => {
  if (!props.customerId) return;
  loading.value = true;
  try {
    const response = await usersApi.listByCustomer(props.customerId);
    users.value = response.data || [];
  } catch (e) {
    toast.add({
      severity: "error",
      summary: t("common.error"),
      detail: t("common.loading") + " " + t("common.error"), // Simplified or specific key
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
};

// Initial load and watch for prop change
onMounted(fetchUsers);
watch(() => props.customerId, fetchUsers);

const openCreateDialog = () => {
  isEdit.value = false;
  currentUserId.value = null;
  form.value = {
    customer_id: props.customerId, // AUTO-ASSOCIATION
    username: "",
    email: "",
    password: "",
    is_active: true,
    is_admin: false,
  };
  showDialog.value = true;
};

const openEditDialog = (user: User) => {
  isEdit.value = true;
  currentUserId.value = user.id;
  // Copy values
  form.value = {
    customer_id: user.customer_id,
    username: user.username,
    email: user.email,
    password: "", // Do not fill password on edit
    is_active: user.is_active,
    is_admin: user.is_admin,
  };
  showDialog.value = true;
};

const saveUser = async () => {
  submitting.value = true;
  try {
    // Ensure customer_id is strictly from the prop context if creating
    if (!isEdit.value) {
      form.value.customer_id = props.customerId;
    }

    const payload = { ...form.value };
    // If password empty on edit, remove it to avoid overwriting
    if (isEdit.value && !payload.password) {
      delete payload.password;
    }

    if (isEdit.value && currentUserId.value) {
      await usersApi.update(currentUserId.value, payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("users.save_success"),
        life: 3000,
      });
    } else {
      await usersApi.create(payload);
      toast.add({
        severity: "success",
        summary: t("common.success"),
        detail: t("users.save_success"),
        life: 3000,
      });
    }
    showDialog.value = false;
    fetchUsers();
  } catch (e: any) {
    const msg =
      e.response?.data?.error || e.message || "No s'ha pogut guardar l'usuari";
    toast.add({
      severity: "error",
      summary: "Error",
      detail: msg,
      life: 3000,
    });
  } finally {
    submitting.value = false;
  }
};

const confirmDelete = async (id: string) => {
  if (!confirm(t("users.delete_confirm"))) return;
  try {
    await usersApi.delete(id);
    toast.add({
      severity: "success",
      summary: t("common.success"),
      detail: "Usuari eliminat", // Could add key 'user_deleted'
      life: 3000,
    });
    fetchUsers();
  } catch (e: any) {
    const msg =
      e.response?.data?.error || e.message || "No s'ha pogut eliminar";
    toast.add({
      severity: "error",
      summary: "Error",
      detail: msg,
      life: 3000,
    });
  }
};
</script>

<template>
  <div class="flex flex-col gap-4">
    <!-- Actions -->
    <div class="flex justify-end" v-if="!readOnly">
      <Button
        :label="t('users.new_user')"
        icon="pi pi-plus"
        size="small"
        @click="openCreateDialog"
      />
    </div>

    <!-- Table -->
    <DataTable :value="users" :loading="loading" stripedRows size="small">
      <template #empty>{{ t("users.no_users") }}</template>
      <Column field="username" :header="t('users.username')" sortable></Column>
      <Column field="email" :header="t('users.email')" sortable></Column>
      <Column field="is_admin" :header="t('users.role')">
        <template #body="slotProps">
          <Tag
            :value="
              slotProps.data.is_admin ? t('users.admin') : t('users.user')
            "
            :severity="slotProps.data.is_admin ? 'info' : 'secondary'"
          />
        </template>
      </Column>
      <Column field="is_active" :header="t('users.active')">
        <template #body="slotProps">
          <Tag
            :value="
              slotProps.data.is_active ? t('users.active') : t('users.inactive')
            "
            :severity="slotProps.data.is_active ? 'success' : 'danger'"
          />
        </template>
      </Column>
      <Column :header="t('users.actions')" class="text-right" v-if="!readOnly">
        <template #body="slotProps">
          <div class="flex justify-end gap-2">
            <Button
              icon="pi pi-pencil"
              text
              rounded
              severity="secondary"
              @click="openEditDialog(slotProps.data)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              severity="danger"
              @click="confirmDelete(slotProps.data.id)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Dialog -->
    <Dialog
      v-model:visible="showDialog"
      :header="isEdit ? t('users.edit_user') : t('users.new_user')"
      modal
      :style="{ width: '30rem' }"
    >
      <div class="flex flex-col gap-4 mb-4 pt-2">
        <div class="flex flex-col gap-2">
          <label for="username" class="font-semibold text-sm">{{
            t("users.username")
          }}</label>
          <InputText id="username" v-model="form.username" autocomplete="off" />
        </div>
        <div class="flex flex-col gap-2">
          <label for="email" class="font-semibold text-sm">{{
            t("users.email")
          }}</label>
          <InputText id="email" v-model="form.email" autocomplete="off" />
        </div>
        <div class="flex flex-col gap-2">
          <label for="password" class="font-semibold text-sm"
            >{{ t("users.password") }}
            {{ isEdit ? t("users.password_hint") : "" }}</label
          >
          <Password
            id="password"
            v-model="form.password"
            :feedback="false"
            toggleMask
            class="w-full"
            inputClass="w-full"
            autocomplete="new-password"
          />
        </div>
        <div
          class="flex items-center justify-between border-t border-surface-100 pt-4 mt-2"
        >
          <div class="flex items-center gap-2">
            <ToggleSwitch v-model="form.is_active" inputId="is_active" />
            <label for="is_active" class="cursor-pointer">{{
              t("users.active")
            }}</label>
          </div>
          <div class="flex items-center gap-2">
            <ToggleSwitch v-model="form.is_admin" inputId="is_admin" />
            <label for="is_admin" class="cursor-pointer">Admin</label>
          </div>
        </div>
      </div>
      <div class="flex justify-end gap-2">
        <Button
          :label="t('common.cancel')"
          severity="secondary"
          text
          @click="showDialog = false"
        />
        <Button
          :label="t('common.save')"
          @click="saveUser"
          :loading="submitting"
        />
      </div>
    </Dialog>
  </div>
</template>
