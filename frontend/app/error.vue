<script setup lang="ts">
import type { NuxtError } from "#app";

const props = defineProps<{
  error: NuxtError;
}>();

const statusCode = computed(() => props.error.statusCode || 500);
const title = computed(() => {
  if (statusCode.value === 404) {
    return "Page not found";
  }
  return "Something went wrong";
});

const description = computed(() => {
  if (statusCode.value === 404) {
    return "This route does not exist in your workspace.";
  }
  return props.error.statusMessage || "Unexpected application error.";
});

function clear(): void {
  clearError({ redirect: "/" });
}
</script>

<template>
  <main class="flex min-h-svh items-center justify-center px-4">
    <UPageCard :title="title" :description="description" class="w-full max-w-md">
      <template #title>
        <span class="inline-flex items-center gap-2 tracking-tight">
          <!-- Errors are warning-hued, not accent-hued. The accent means
               "you can act here", which is the opposite of what this is. -->
          <UIcon name="i-lucide-triangle-alert" class="size-5 text-warning" />
          {{ title }}
        </span>
      </template>

      <div class="space-y-4 text-sm text-muted">
        <p>
          Status code
          <span class="mb-measure text-toned">{{ statusCode }}</span>
        </p>
        <UButton label="Back to overview" icon="i-lucide-home" @click="clear" />
      </div>
    </UPageCard>
  </main>
</template>
