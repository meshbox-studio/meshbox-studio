<script setup lang="ts">
const { stats, refreshStats } = useSidebarStats();

onMounted(() => {
  void refreshStats();
});
</script>

<template>
  <div class="flex items-center justify-between text-xs text-muted px-3 py-1.5">
    <span>v{{ stats.version }}</span>
  </div>
</template>
