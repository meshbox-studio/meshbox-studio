<script setup lang="ts">
definePageMeta({
  title: "Overview",
});

const { activeProjects, recentActivity, fetchAll } = useProjects();
const { stats, refreshStats } = useSidebarStats();

const activeCount = computed(
  () => activeProjects.value.filter((project) => project.state !== "draft").length,
);

const printingCount = computed(
  () => activeProjects.value.filter((project) => project.state === "printing").length,
);

const draftCount = computed(
  () => activeProjects.value.filter((project) => project.state === "draft").length,
);

onMounted(() => {
  void fetchAll();
  void refreshStats();
});
</script>

<template>
  <div class="space-y-6">
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between text-sm text-muted">
            Active projects
            <UIcon name="i-lucide-folder-open" class="size-4" />
          </div>
        </template>
        <p class="text-3xl font-semibold text-highlighted">{{ activeCount }}</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between text-sm text-muted">
            Printing now
            <UIcon name="i-lucide-printer" class="size-4" />
          </div>
        </template>
        <p class="text-3xl font-semibold text-highlighted">{{ printingCount }}</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between text-sm text-muted">
            Draft projects
            <UIcon name="i-lucide-file-pen-line" class="size-4" />
          </div>
        </template>
        <p class="text-3xl font-semibold text-highlighted">{{ draftCount }}</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between text-sm text-muted">
            Disk usage
            <UIcon name="i-lucide-hard-drive" class="size-4" />
          </div>
        </template>
        <p class="text-3xl font-semibold text-highlighted">{{ stats.diskPercent }}%</p>
        <p class="mt-1 text-sm text-muted">{{ stats.diskUsed }} / {{ stats.diskTotal }}</p>
      </UCard>
    </section>

    <section class="grid gap-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-base font-semibold text-highlighted">Recent activity</h2>
              <p class="text-sm text-muted">Quick timeline across projects, prints, and status changes.</p>
            </div>
            <UBadge color="neutral" variant="subtle" :label="`${recentActivity.length} events`" />
          </div>
        </template>

        <ul class="space-y-3">
          <li
            v-for="event in recentActivity"
            :key="event.id"
            class="rounded-lg border border-default bg-default px-3 py-2"
          >
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="font-medium text-highlighted">{{ event.title }}</p>
                <p class="text-sm text-muted">{{ event.description }}</p>
              </div>
              <UBadge color="neutral" variant="soft" size="sm" :label="event.at" />
            </div>
          </li>
        </ul>
      </UCard>
    </section>
  </div>
</template>
