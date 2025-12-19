import { createApp } from "vue";
import { createPinia } from "pinia";
import PrimeVue from "primevue/config";
import Aura from "@primevue/themes/aura";
import ToastService from "primevue/toastservice";
import ConfirmationService from "primevue/confirmationservice";
import App from "./App.vue";
import router from "./router";

// Import styles
import "./style.css";
import "primeicons/primeicons.css";

import { useAuthStore } from "./stores/auth.store";
import i18n from "./i18n";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(i18n);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
});
app.use(ToastService);
app.use(ConfirmationService);

const authStore = useAuthStore();
authStore.loadFromStorage();

console.log("App mounting...");
app.mount("#app");
