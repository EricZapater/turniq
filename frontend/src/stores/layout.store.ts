import { defineStore } from "pinia";
import { ref } from "vue";

export type LayoutMode = "DESPATX" | "PLANTA";

export const useLayoutStore = defineStore("layout", () => {
  const mode = ref<LayoutMode>("DESPATX");

  function toggleMode() {
    mode.value = mode.value === "DESPATX" ? "PLANTA" : "DESPATX";
  }

  function setMode(newMode: LayoutMode) {
    mode.value = newMode;
  }

  return {
    mode,
    toggleMode,
    setMode,
  };
});
