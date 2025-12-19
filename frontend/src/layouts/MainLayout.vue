<script setup lang="ts">
import { computed, ref, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth.store";
import { useLayoutStore } from "../stores/layout.store";
import { useCustomersStore } from "../stores/customers.store";
import { useWorkcentersStore } from "../stores/workcenters.store";
import { useJobsStore } from "../stores/jobs.store";
import { useRouter } from "vue-router";
import { useToast } from "primevue/usetoast";
import { operatorsApi } from "../api/operators.api";
import { timeentriesApi } from "../api/timeentries.api";
import { scheduleApi } from "../api/schedule.api";

import { shiftsApi } from "../api/shifts.api";
import { useConfirm } from "primevue/useconfirm";

// PrimeVue components
import Button from "primevue/button";
import Avatar from "primevue/avatar";
import Menu from "primevue/menu";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import ConfirmDialog from "primevue/confirmdialog";
import Select from "primevue/select";

const authStore = useAuthStore();
const layoutStore = useLayoutStore();
const customersStore = useCustomersStore();
const workcentersStore = useWorkcentersStore();
const jobsStore = useJobsStore();
const router = useRouter();
const toast = useToast();

const confirm = useConfirm();
const { t, locale } = useI18n();

// Language State
const currentLang = ref(locale.value);
const langOptions = [
  { label: "Català", value: "ca" },
  { label: "Español", value: "es" },
  { label: "English", value: "en" },
];

watch(currentLang, (newLang) => {
  locale.value = newLang;
  // Here we could persist to localStorage if separate from user profile
  localStorage.setItem("turniq_lang", newLang);
});

// Watch Auth User to set Language
watch(
  () => authStore.user,
  async (user) => {
    // Priority: 1. LocalStorage (User Pref), 2. Customer Config, 3. Default (CA)
    const stored = localStorage.getItem("turniq_lang");
    if (stored) {
      if (locale.value !== stored) {
        locale.value = stored;
        currentLang.value = stored;
      }
      return;
    }

    if (user && user.customer_id) {
      try {
        // Fetch customer to get language
        // We use the store but be careful not to override 'currentCustomer' if used elsewhere in a conflicting way,
        // but fetchCustomer usually sets currentCustomer. MainLayout doesn't display customer details so it's fine.
        // Or better, use a lighter fetch if available, but fetchCustomer is standard.
        await customersStore.fetchCustomer(user.customer_id);
        const cust = customersStore.currentCustomer;
        if (
          cust &&
          cust.language &&
          ["ca", "es", "en"].includes(cust.language)
        ) {
          locale.value = cust.language;
          currentLang.value = cust.language;
        }
      } catch (e) {
        console.warn("Could not load customer for language settings", e);
      }
    }
  },
  { immediate: true }
);

onMounted(() => {
  // Ensure we have metadata for mapping names
  if (!workcentersStore.workcenters.length)
    workcentersStore.fetchWorkcenters({});
  if (!jobsStore.jobs.length) jobsStore.fetchJobs({});
});

// Planning Modal State
const showPlanningDialog = ref(false);
const planningStep = ref<"code" | "list">("code");
const planningCode = ref("");
const planningOperator = ref<any>(null);
const planningEntries = ref<any[]>([]); // Use any to allow enriched fields
const loadingPlanning = ref(false);

// Stop Confirmation specific
const showStopConfirm = ref(false);
const entryToStop = ref<any | null>(null);

const openPlanningDialog = () => {
  planningCode.value = "";
  planningStep.value = "code";
  planningOperator.value = null;
  planningEntries.value = [];
  showPlanningDialog.value = true;
};

const checkPlanningCode = async () => {
  if (!planningCode.value) return;
  loadingPlanning.value = true;
  try {
    const opRes = await operatorsApi.findByCode(planningCode.value);
    const operator = opRes.data;
    if (!operator) throw new Error("Operari no trobat");

    planningOperator.value = operator;

    // Ensure stores are loaded
    if (!workcentersStore.workcenters.length)
      await workcentersStore.fetchWorkcenters({});
    if (!jobsStore.jobs.length) await jobsStore.fetchJobs({});

    // Fetch Schedule for Today
    const today = new Date().toISOString().split("T")[0];
    const entries = await scheduleApi.getOperatorPlanning(operator.id, today);

    // Enrich Data
    planningEntries.value = entries.map((e: any) => {
      const wc = workcentersStore.workcenters.find(
        (w: any) => w.id === e.workcenter_id
      );
      const job = jobsStore.jobs.find((j: any) => j.id === e.job_id);
      return {
        ...e,
        workcenter_name: wc ? wc.name : "Unknown",
        product_code: job ? job.product_code : "-",
        job_code: job ? job.job_code : "-",
        estimated_duration: job ? job.estimated_duration : 0,
        // Add friendly start time
        start_time_fmt: e.start_time
          ? new Date(e.start_time).toLocaleTimeString([], {
              hour: "2-digit",
              minute: "2-digit",
            })
          : "-",
      };
    });

    planningStep.value = "list";
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: e.response?.data?.error || "Error cercant planificació",
      life: 3000,
    });
  } finally {
    loadingPlanning.value = false;
  }
};

const handleStart = async (entry: any) => {
  // 1. Check Shift
  const inShift = await checkShiftAdherence(planningOperator.value);
  if (!inShift) return;

  // 2. Check Clock In
  try {
    const active = await timeentriesApi.findCurrent(planningOperator.value.id);
    if (!active.data || !active.data.id) {
      toast.add({
        severity: "error",
        summary: "No has entrat",
        detail: "Has de fitxar l'entrada primer.",
        life: 4000,
      });
      return;
    }
  } catch (e) {
    toast.add({
      severity: "error",
      summary: "No has entrat",
      detail: "Has de fitxar l'entrada primer.",
      life: 4000,
    });
    return;
  }

  loadingPlanning.value = true;
  try {
    // Start: Update StartTime to Now
    await scheduleApi.update(entry.id, {
      ...entry,
      start_time: new Date().toISOString(),
      // Map required fields from entry
      customer_id: entry.customer_id,
      shopfloor_id: entry.shopfloor_id,
      shift_id: entry.shift_id,
      workcenter_id: entry.workcenter_id,
      job_id: entry.job_id,
      operator_id: entry.operator_id,
      date: entry.date,
      order: entry.order,
      is_completed: entry.is_completed, // preserve
    });

    toast.add({ severity: "success", summary: "Iniciat", life: 2000 });
    showPlanningDialog.value = false;
  } catch (e) {
    toast.add({ severity: "error", summary: "Error iniciant", life: 3000 });
  } finally {
    loadingPlanning.value = false;
  }
};

const handleStopClick = (entry: any) => {
  entryToStop.value = entry;
  showStopConfirm.value = true;
};

const confirmStop = async (completed: boolean) => {
  if (!entryToStop.value) return;
  loadingPlanning.value = true;
  showStopConfirm.value = false;

  try {
    const entry = entryToStop.value;
    await scheduleApi.update(entry.id, {
      customer_id: entry.customer_id,
      shopfloor_id: entry.shopfloor_id,
      shift_id: entry.shift_id,
      workcenter_id: entry.workcenter_id,
      job_id: entry.job_id,
      operator_id: entry.operator_id,
      date: entry.date,
      order: entry.order,
      start_time: entry.start_time,
      end_time: new Date().toISOString(),
      is_completed: completed,
    });

    toast.add({ severity: "success", summary: "Finalitzat", life: 2000 });
    entryToStop.value = null;
    showPlanningDialog.value = false;
  } catch (e) {
    toast.add({ severity: "error", summary: "Error finalitzant", life: 3000 });
  } finally {
    loadingPlanning.value = false;
  }
};

const handleRestart = async (entry: any) => {
  loadingPlanning.value = true;
  try {
    // Create NEW entry based on this one
    await scheduleApi.create({
      customer_id: entry.customer_id,
      shopfloor_id: entry.shopfloor_id,
      shift_id: entry.shift_id,
      workcenter_id: entry.workcenter_id,
      job_id: entry.job_id,
      operator_id: entry.operator_id,
      date: new Date().toISOString().split("T")[0], // Today
      order: entry.order,
      start_time: new Date().toISOString(), // Start Immediately
      end_time: null,
      is_completed: false,
    });

    toast.add({
      severity: "success",
      summary: "Reiniciat (Nova Entrada)",
      life: 2000,
    });
    showPlanningDialog.value = false; // Close modal
  } catch (e) {
    toast.add({ severity: "error", summary: "Error reiniciant", life: 3000 });
  } finally {
    loadingPlanning.value = false;
  }
};

// Auth Modal State (Enter/Exit)
const showAuthDialog = ref(false);
const authMode = ref<"ENTER" | "EXIT">("ENTER");
const authCode = ref("");
const authStep = ref<"code" | "confirm">("code");
const authOperator = ref<any>(null);
const authEntry = ref<any>(null); // For Exit mode mainly
const loadingAuth = ref(false);

// --- VALIDATION HELPERS ---
const checkShiftAdherence = async (operator: any): Promise<boolean> => {
  if (!operator || !operator.shop_floor_id) return true;
  try {
    const shiftsRes = await shiftsApi.listByShopfloor(operator.shop_floor_id);
    const shifts = shiftsRes.data || [];
    if (shifts.length === 0) return true;

    const now = new Date();
    const currentMin = now.getHours() * 60 + now.getMinutes();

    const isInside = shifts.some((s: any) => {
      if (!s.is_active) return false;
      // Parse HH:mm
      const [sh, sm] = s.start_time.split(":").map(Number);
      const [eh, em] = s.end_time.split(":").map(Number);
      const startMin = sh * 60 + sm;
      const endMin = eh * 60 + em;

      if (startMin <= endMin) {
        return currentMin >= startMin && currentMin <= endMin;
      } else {
        return currentMin >= startMin || currentMin <= endMin;
      }
    });

    if (!isInside) {
      return new Promise<boolean>((resolve) => {
        confirm.require({
          message:
            "Atenció: Estàs fora de l'horari del torn assignat. Vols continuar igualment?",
          header: "Fora de Torn",
          icon: "pi pi-exclamation-triangle",
          accept: () => resolve(true),
          reject: () => resolve(false),
          onHide: () => resolve(false), // If cancelled or closed
        });
      });
    }
    return true;
  } catch (e) {
    console.error("Error checking shifts", e);
    return true; // Fail safe
  }
};

const openAuthDialog = (mode: "ENTER" | "EXIT") => {
  authMode.value = mode;
  authCode.value = "";
  authStep.value = "code";
  authOperator.value = null;
  authEntry.value = null;
  showAuthDialog.value = true;
};

const checkAuthCode = async () => {
  if (!authCode.value) return;
  loadingAuth.value = true;
  try {
    const opRes = await operatorsApi.findByCode(authCode.value);
    const operator = opRes.data;
    if (!operator) throw new Error("Operari no trobat");

    authOperator.value = operator;

    // Check current session
    let entry = null;
    try {
      const currentRes = await timeentriesApi.findCurrent(operator.id);
      entry = currentRes.data;
    } catch (e) {
      // Ignore error, assumes no session found if 404/500
    }
    const hasSession = !!(entry && entry.id);

    // CHECK SHIFT before proceeding
    const inShift = await checkShiftAdherence(operator);
    if (!inShift) {
      loadingAuth.value = false;
      return; // Stopped by confirm reject
    }

    if (authMode.value === "ENTER") {
      if (hasSession) {
        toast.add({
          severity: "warn",
          summary: "Ja has entrat",
          detail: "Tens una sessió oberta. Has de sortir primer.",
          life: 3000,
        });
        return;
      }
      await performEnter(operator);
    } else {
      // EXIT
      if (!hasSession) {
        toast.add({
          severity: "warn",
          summary: "Sense registre",
          detail: "No has entrat. Has d'entrar primer.",
          life: 3000,
        });
        return;
      }
      // Has session -> Confirm Step
      authEntry.value = entry;
      authStep.value = "confirm";
    }
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: e.response?.data?.error || "Error cercant operari",
      life: 3000,
    });
  } finally {
    loadingAuth.value = false;
  }
};

const performEnter = async (operator: any) => {
  try {
    await timeentriesApi.create({
      operator_id: operator.id,
      check_in: new Date().toISOString(),
    });
    toast.add({
      severity: "success",
      summary: "Benvingut/da",
      detail: `Hola ${operator.name}, sessió iniciada.`,
      life: 3000,
    });
    showAuthDialog.value = false;
  } catch (e) {
    throw e; // Caught by caller
  }
};

const performExit = async () => {
  if (!authEntry.value) return;
  loadingAuth.value = true;
  try {
    await timeentriesApi.update(authEntry.value.id, {
      operator_id: authOperator.value.id,
      check_out: new Date().toISOString(),
    });

    toast.add({
      severity: "success",
      summary: "Fins aviat!",
      detail: "Sortida registrada correctament.",
      life: 3000,
    });
    showAuthDialog.value = false;
  } catch (e: any) {
    toast.add({
      severity: "error",
      summary: "Error",
      detail: "Error al registrar sortida",
      life: 3000,
    });
  } finally {
    loadingAuth.value = false;
  }
};

// User menu items
const menu = ref();
const userContextMenuItems = [
  {
    label: "Veure perfil",
    icon: "pi pi-user",
    command: () => {
      router.push("/my-customer");
    },
  },
  {
    separator: true,
  },
  {
    label: "Logout",
    icon: "pi pi-sign-out",
    command: () => {
      authStore.logout();
      router.push("/login");
    },
  },
];

const toggleUserMenu = (event: any) => {
  menu.value.toggle(event);
};

// Side menu items placeholder
const sideMenuItems = computed(() => {
  const items = [{ label: t("menu.dashboard"), icon: "pi pi-home", to: "/" }];

  if (authStore.isAdmin) {
    items.push({
      label: t("menu.customers"),
      icon: "pi pi-users",
      to: "/customers",
    });
  } else {
    items.push({
      label: t("menu.users"),
      icon: "pi pi-users",
      to: "/company-users",
    });
  }

  items.push(
    {
      label: t("menu.shopfloors"),
      icon: "pi pi-building",
      to: "/shopfloors",
    },
    { label: t("menu.planning"), icon: "pi pi-calendar", to: "/planning" },
    { label: t("menu.shifts"), icon: "pi pi-clock", to: "/shifts" },
    { label: t("menu.workcenters"), icon: "pi pi-cog", to: "/workcenters" },
    { label: t("menu.operators"), icon: "pi pi-id-card", to: "/operators" },
    { label: t("menu.jobs"), icon: "pi pi-briefcase", to: "/jobs" }
  );

  // Reports
  items.push(
    {
      label: t("menu.schedule_report"),
      icon: "pi pi-list",
      to: "/reports/schedule",
    },
    {
      label: t("menu.hours_report"),
      icon: "pi pi-stopwatch",
      to: "/reports/time-entries",
    }
  );

  if (authStore.isAdmin) {
    items.push({
      label: t("menu.billing_report"),
      icon: "pi pi-money-bill",
      to: "/reports/billing",
    });
  }

  return items;
});

const isOfficeMode = computed(() => layoutStore.mode === "DESPATX");
</script>

<template>
  <!-- Root: Always Surface-50 (Zinc-50) for a clean industrial background -->
  <div
    class="layout-wrapper min-h-screen flex flex-col bg-surface-50 text-surface-900 font-sans"
  >
    <!-- TOP BAR: White, Subtle Border, High Contrast Text -->
    <header
      class="topbar h-14 bg-white border-b border-surface-200 px-6 flex items-center justify-between shadow-sm z-50"
    >
      <!-- LEFT: Branding and Language -->
      <div class="flex items-center gap-4">
        <div class="flex flex-col items-start leading-none mr-4">
          <!-- Industrial Logo Style: Dark Slate, Bold, Tracking Tight -->
          <span class="text-lg font-bold tracking-tight text-primary-800"
            >turniq</span
          >
          <span
            class="text-[0.6rem] text-surface-400 font-mono tracking-widest uppercase"
            >by >_rawcraft</span
          >
        </div>

        <!-- Language Selector -->
        <Select
          v-model="currentLang"
          :options="langOptions"
          optionLabel="label"
          optionValue="value"
          class="w-32 h-8 text-xs"
        />
      </div>

      <!-- CENTER: Mode Toggle (Industrial Switch) -->
      <div
        class="flex items-center gap-3 bg-surface-100 px-1 py-1 rounded-md border border-surface-200"
      >
        <button
          @click="layoutStore.setMode('DESPATX')"
          class="px-3 py-1 text-xs font-semibold rounded transition-all duration-200"
          :class="
            isOfficeMode
              ? 'bg-white text-primary-800 shadow-sm border border-surface-200'
              : 'text-surface-500 hover:text-surface-700'
          "
        >
          DESPATX
        </button>
        <button
          @click="layoutStore.setMode('PLANTA')"
          class="px-3 py-1 text-xs font-semibold rounded transition-all duration-200"
          :class="
            !isOfficeMode
              ? 'bg-white text-primary-800 shadow-sm border border-surface-200'
              : 'text-surface-500 hover:text-surface-700'
          "
        >
          PLANTA
        </button>
      </div>

      <!-- RIGHT: User Profile -->
      <div class="flex items-center gap-3">
        <div
          class="flex items-center gap-2 cursor-pointer hover:bg-surface-50 px-2 py-1.5 rounded transition-colors border border-transparent hover:border-surface-200"
          @click="toggleUserMenu"
          aria-haspopup="true"
          aria-controls="overlay_menu"
        >
          <span class="text-sm font-medium text-surface-700 hidden sm:block">{{
            authStore.user?.username || "Usuari"
          }}</span>
          <!-- Avatar: Subdued Interaction -->
          <Avatar
            :label="authStore.user?.username?.charAt(0).toUpperCase() || 'U'"
            shape="square"
            class="bg-surface-200 text-surface-700 !w-8 !h-8 font-bold rounded-md"
          />
          <i class="pi pi-angle-down text-xs text-surface-400"></i>
        </div>
        <Menu
          ref="menu"
          id="overlay_menu"
          :model="userContextMenuItems"
          :popup="true"
        />
      </div>
    </header>

    <!-- MAIN LAYOUT BODY -->
    <div class="flex flex-1 overflow-hidden relative">
      <!-- SIDE MENU (OFFICE MODE ONLY) -->
      <aside
        v-if="isOfficeMode"
        class="w-64 bg-white border-r border-surface-200 flex flex-col"
      >
        <!-- Menu Items: Clean lists, nice hover states -->
        <div class="p-4 flex-1">
          <div
            class="text-xs font-semibold text-surface-400 uppercase tracking-wider mb-4 px-3"
          >
            Navegació
          </div>
          <ul class="space-y-1">
            <li v-for="item in sideMenuItems" :key="item.label">
              <router-link
                :to="item.to || '#'"
                class="flex items-center gap-3 px-3 py-2 text-surface-600 hover:bg-surface-50 hover:text-primary-700 rounded transition-colors group"
                active-class="bg-surface-100 text-primary-700 font-bold"
              >
                <i
                  :class="item.icon"
                  class="text-surface-400 group-hover:text-primary-600 transition-colors"
                ></i>
                <span class="text-sm">{{ item.label }}</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Footer Info -->
        <div
          class="p-4 border-t border-surface-100 text-[10px] text-center text-surface-300"
        >
          &copy; 2025 Turniq System v1.0
        </div>
      </aside>

      <!-- CONTENT AREA -->
      <main class="flex-1 overflow-auto bg-surface-50 relative flex flex-col">
        <!-- OFFICE CONTENT -->
        <div v-if="isOfficeMode" class="flex-1 p-6">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>

        <!-- PLANT CONTENT -->
        <div
          v-else
          class="flex-1 flex items-center justify-center p-8 bg-surface-100"
        >
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 w-full max-w-5xl">
            <!-- Large Industrial Buttons -->
            <!-- Using plain HTML/Tailwind for maximum control over "Industrial" look vs default PrimeVue styles -->

            <button
              @click="openPlanningDialog"
              class="bg-white border-l-4 border-primary-600 p-8 shadow-sm hover:shadow-md transition-all flex flex-col items-center justify-center gap-4 group h-48 rounded-r-lg md:col-span-2"
            >
              <i
                class="pi pi-calendar text-4xl text-surface-400 group-hover:text-primary-600 transition-colors"
              ></i>
              <span
                class="text-xl font-bold text-surface-700 group-hover:text-primary-800"
                >Consultar Planificació</span
              >
            </button>

            <button
              @click="openAuthDialog('ENTER')"
              class="bg-white border-l-4 border-green-600 p-8 shadow-sm hover:shadow-md transition-all flex flex-col items-center justify-center gap-4 group h-48 rounded-r-lg"
            >
              <i
                class="pi pi-sign-in text-4xl text-surface-400 group-hover:text-green-600 transition-colors"
              ></i>
              <span
                class="text-xl font-bold text-surface-700 group-hover:text-green-800"
                >Entrar</span
              >
            </button>

            <button
              @click="openAuthDialog('EXIT')"
              class="bg-white border-l-4 border-amber-500 p-8 shadow-sm hover:shadow-md transition-all flex flex-col items-center justify-center gap-4 group h-48 rounded-r-lg"
            >
              <i
                class="pi pi-sign-out text-4xl text-surface-400 group-hover:text-amber-600 transition-colors"
              ></i>
              <span
                class="text-xl font-bold text-surface-700 group-hover:text-amber-800"
                >Sortir</span
              >
            </button>
          </div>
        </div>
      </main>
    </div>

    <!-- AUTH MODAL (ENTER/EXIT) -->
    <Dialog
      v-model:visible="showAuthDialog"
      :header="authMode === 'ENTER' ? 'Registrar Entrada' : 'Registrar Sortida'"
      :modal="true"
      class="w-full max-w-md"
      :closable="true"
    >
      <div v-if="authStep === 'code'" class="flex flex-col gap-4 pt-2">
        <label for="authCode" class="font-medium">Codi d'Operari</label>
        <InputText
          id="authCode"
          v-model="authCode"
          placeholder="Introdueix el teu codi..."
          class="w-full text-center text-xl p-3"
          @keydown.enter="checkAuthCode"
          autofocus
        />
        <Button
          label="Continuar"
          class="w-full mt-2"
          :severity="authMode === 'ENTER' ? 'success' : 'warning'"
          @click="checkAuthCode"
          :loading="loadingAuth"
        />
      </div>

      <div
        v-else-if="authStep === 'confirm' && authOperator && authEntry"
        class="flex flex-col gap-4 text-center"
      >
        <div class="text-xl font-bold text-surface-800">
          Hola {{ authOperator.name }}
        </div>
        <div class="text-surface-600">
          Tens una sessió oberta des de:
          <div class="text-lg font-mono text-primary-700 font-bold mt-1">
            {{ new Date(authEntry.check_in).toLocaleString() }}
          </div>
        </div>

        <div class="flex flex-col gap-2 mt-4">
          <Button
            label="CONFIRMAR SORTIDA"
            severity="danger"
            class="w-full font-bold p-3"
            @click="performExit"
            :loading="loadingAuth"
          />
          <Button
            label="Cancel·lar"
            severity="secondary"
            text
            class="w-full"
            @click="showAuthDialog = false"
          />
        </div>
      </div>
    </Dialog>
    <!-- PLANNING MODAL -->
    <Dialog
      v-model:visible="showPlanningDialog"
      header="La teva planificació (Avui)"
      :modal="true"
      class="w-full max-w-4xl"
      :closable="true"
      :maximizable="true"
    >
      <div
        v-if="planningStep === 'code'"
        class="flex flex-col gap-4 pt-2 max-w-md mx-auto"
      >
        <label for="planCode" class="font-medium">Codi d'Operari</label>
        <InputText
          id="planCode"
          v-model="planningCode"
          placeholder="Introdueix el teu codi..."
          class="w-full text-center text-xl p-3"
          @keydown.enter="checkPlanningCode"
          autofocus
        />
        <Button
          label="Veure Planificació"
          class="w-full mt-2"
          @click="checkPlanningCode"
          :loading="loadingPlanning"
        />
      </div>

      <div v-else-if="planningStep === 'list'" class="flex flex-col gap-4">
        <div v-if="planningOperator" class="text-xl font-bold mb-2">
          Operari:
          <span class="text-primary-600">{{ planningOperator.name }}</span>
        </div>

        <DataTable
          :value="planningEntries"
          stripedRows
          responsiveLayout="scroll"
        >
          <Column field="workcenter_name" header="Centre de Treball"></Column>
          <Column field="product_code" header="Codi Producte"></Column>
          <Column field="job_code" header="Codi Feina"></Column>
          <Column field="estimated_duration" header="Durada Est. (min)">
            <template #body="slotProps">
              {{ slotProps.data.estimated_duration }} min
            </template>
          </Column>
          <Column field="start_time_fmt" header="Hora Inici"></Column>
          <Column header="Accions">
            <template #body="slotProps">
              <div class="flex gap-2 items-center">
                <Button
                  v-if="!slotProps.data.start_time"
                  label="Iniciar"
                  icon="pi pi-play"
                  severity="success"
                  @click="handleStart(slotProps.data)"
                  :loading="loadingPlanning"
                />
                <Button
                  v-else-if="!slotProps.data.end_time"
                  label="Parar"
                  icon="pi pi-stop"
                  severity="danger"
                  @click="handleStopClick(slotProps.data)"
                  :loading="loadingPlanning"
                />
                <div v-else class="flex flex-col gap-1 items-center">
                  <span
                    class="text-sm font-bold"
                    :class="
                      slotProps.data.is_completed
                        ? 'text-green-600'
                        : 'text-orange-500'
                    "
                  >
                    {{ slotProps.data.is_completed ? "Completat" : "Aturat" }}
                  </span>
                  <Button
                    label="Reiniciar"
                    icon="pi pi-refresh"
                    severity="info"
                    size="small"
                    @click="handleRestart(slotProps.data)"
                    :loading="loadingPlanning"
                  />
                </div>
              </div>
            </template>
          </Column>
        </DataTable>

        <div class="mt-4 text-center">
          <Button
            label="Tancar"
            severity="secondary"
            @click="showPlanningDialog = false"
          />
        </div>
      </div>
    </Dialog>

    <!-- STOP CONFIRMATION DIALOG -->
    <Dialog
      v-model:visible="showStopConfirm"
      header="Confirmació"
      :modal="true"
      class="w-full max-w-sm"
    >
      <div class="flex flex-col gap-4 text-center">
        <p class="text-lg">Has acabat aquesta feina?</p>
        <div class="flex flex-col gap-2">
          <Button
            label="SÍ, FEINA ACABADA"
            severity="success"
            @click="confirmStop(true)"
          />
          <Button
            label="NO, només aturo temporalment"
            severity="warning"
            @click="confirmStop(false)"
          />
          <Button
            label="Cancel·lar"
            severity="secondary"
            text
            @click="showStopConfirm = false"
          />
        </div>
      </div>
    </Dialog>

    <ConfirmDialog></ConfirmDialog>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
