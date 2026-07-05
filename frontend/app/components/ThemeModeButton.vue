<script setup lang="ts">
const colorMode = useColorMode();

const icons: Record<string, string> = {
  dark: "i-lucide-moon-star",
  light: "i-lucide-sun",
  system: "i-lucide-monitor",
};

const labels: Record<string, string> = {
  dark: "Switch to light mode",
  light: "Switch to system mode",
  system: "Switch to dark mode",
};

const icon = computed(() => icons[colorMode.preference] || icons.system);

function toggleTheme(): void {
  const current = colorMode.preference || "system";
  if (current === "dark") colorMode.preference = "light";
  else if (current === "light") colorMode.preference = "system";
  else colorMode.preference = "dark";
}
</script>

<template>
  <UButton
    :icon="icon"
    color="neutral"
    variant="ghost"
    square
    :aria-label="labels[colorMode.preference] || labels.system"
    @click="toggleTheme"
  />
</template>
