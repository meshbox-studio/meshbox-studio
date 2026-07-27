<script setup lang="ts">
defineProps<{
  collapsed?: boolean;
}>();

const { stats, refreshStats } = useSidebarStats();

onMounted(() => {
  void refreshStats();
});

/* Storage is the one number that becomes a problem if you ignore it, so it
   earns a hue only once it is worth acting on. Below 75% it stays neutral. */
const meterColor = computed(() => {
  if (stats.value.diskPercent >= 90) return "bg-error";
  if (stats.value.diskPercent >= 75) return "bg-warning";
  return "bg-inverted/40";
});
</script>

<template>
  <div class="w-full min-w-0 space-y-1.5 py-1">
    <div v-if="!collapsed" class="flex items-baseline justify-between gap-2">
      <span class="text-xs text-dimmed">Storage</span>
      <span class="mb-measure text-xs text-muted">{{ stats.diskPercent }}%</span>
    </div>

    <div
      class="h-1 w-full overflow-hidden rounded-full bg-accented"
      role="progressbar"
      :aria-valuenow="stats.diskPercent"
      aria-valuemin="0"
      aria-valuemax="100"
      :aria-label="`Storage used: ${stats.diskUsed} of ${stats.diskTotal}`"
    >
      <div
        class="h-full rounded-full transition-[width] duration-500"
        :class="meterColor"
        :style="{ width: `${stats.diskPercent}%` }"
      />
    </div>

    <div
      v-if="!collapsed"
      class="mb-measure flex items-baseline justify-between gap-2 text-[0.6875rem] text-dimmed"
    >
      <span class="truncate">{{ stats.diskUsed }} / {{ stats.diskTotal }}</span>
      <span class="shrink-0">v{{ stats.version }}</span>
    </div>
  </div>
</template>
