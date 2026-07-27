<script setup lang="ts">
import type { ProjectState } from "~/types/meshbox";
import {
  projectStateColor,
  projectStateIsLive,
  projectStateLabel,
  projectStateVariant,
} from "~/utils/project";

const props = defineProps<{
  state: ProjectState;
  size?: "sm" | "md";
}>();

const live = computed(() => projectStateIsLive(props.state));
</script>

<template>
  <UBadge
    :color="projectStateColor(state)"
    :variant="projectStateVariant(state)"
    :size="size ?? 'sm'"
    class="gap-1.5"
  >
    <span
      v-if="live"
      class="mb-live size-1.5 shrink-0 rounded-full bg-current"
      aria-hidden="true"
    />
    {{ projectStateLabel(state) }}
  </UBadge>
</template>
