<script setup lang="ts">
import { projectStateColor, projectStateLabel } from "~/utils/project";

definePageMeta({
  title: "Projects",
});

const route = useRoute();
const router = useRouter();
const { activeProjects, moveToTrash, loaded, fetchAll } = useProjects();

if (!loaded.value) {
  await fetchAll();
}

const q = computed({
  get: () => String(route.query.q ?? ""),
  set: (value: string) => {
    void router.replace({
      query: {
        ...route.query,
        q: value || undefined,
      },
    });
  },
});

const filtered = computed(() => {
  const query = q.value.toLowerCase().trim();
  if (!query) {
    return activeProjects.value;
  }

  return activeProjects.value.filter((project) => {
    const haystack = [
      project.title,
      project.description,
      project.tags.join(" "),
      project.source?.designer ?? "",
      project.source?.category ?? "",
      project.printProfile?.material ?? "",
      project.printProfile?.printer ?? "",
    ]
      .join(" ")
      .toLowerCase();

    return haystack.includes(query);
  });
});

function trashProject(id: string): void {
  moveToTrash(id);
}

function thumbnailUrl(id: string): string {
  return `/api/projects/${id}/thumbnail`;
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <UInput
        v-model="q"
        icon="i-lucide-search"
        placeholder="Search projects, tags..."
        size="lg"
        class="w-full md:max-w-md"
      />
      <UBadge color="neutral" variant="subtle" :label="`${filtered.length} visible`" />
    </div>

    <div class="grid gap-4 lg:grid-cols-2">
      <UCard v-for="project in filtered" :key="project.id">
        <template #header>
          <div class="flex items-start justify-between gap-3">
            <div class="flex gap-3">
              <img
                v-if="project.hasThumbnail"
                :src="thumbnailUrl(project.id)"
                class="size-16 shrink-0 rounded-md border border-default object-cover"
                alt=""
              >
              <div>
                <NuxtLink
                  :to="`/projects/${project.id}`"
                  class="text-base font-semibold text-highlighted transition hover:text-primary-600 dark:hover:text-primary-400"
                >
                  {{ project.title }}
                </NuxtLink>
                <p class="mt-1 text-sm text-muted">{{ project.description }}</p>
              </div>
            </div>
            <UBadge
              :color="projectStateColor(project.state)"
              variant="soft"
              :label="projectStateLabel(project.state)"
            />
          </div>
        </template>

        <div class="space-y-3">
          <div v-if="project.source" class="grid grid-cols-2 gap-2 text-sm">
            <div v-if="project.source.designer" class="rounded-md bg-muted/60 px-2 py-1.5">
              <p class="text-xs text-muted">Designer</p>
              <p class="font-medium text-highlighted">{{ project.source.designer }}</p>
            </div>
            <div v-if="project.source.category" class="rounded-md bg-muted/60 px-2 py-1.5">
              <p class="text-xs text-muted">Category</p>
              <p class="font-medium text-highlighted">{{ project.source.category }}</p>
            </div>
          </div>

          <div v-if="project.license" class="text-xs text-muted">
            <UIcon name="i-lucide-scale" class="inline size-3" />
            {{ project.license.name || 'Licensed' }}
          </div>

          <div class="flex flex-wrap gap-1.5">
            <UBadge
              v-for="tag in project.tags"
              :key="tag"
              color="neutral"
              variant="outline"
              :label="`#${tag}`"
            />
          </div>

          <div class="flex items-center justify-between text-xs text-muted">
            <span>{{ project.files.length }} files</span>
            <span>{{ project.iterations.length }} iterations</span>
            <span>{{ project.lastActivity }}</span>
          </div>
        </div>

        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <UButton
              :to="`/projects/${project.id}`"
              label="Open"
              color="neutral"
              variant="outline"
              icon="i-lucide-arrow-up-right"
            />
            <UButton
              label="Move to trash"
              color="error"
              variant="ghost"
              icon="i-lucide-trash-2"
              @click="trashProject(project.id)"
            />
          </div>
        </template>
      </UCard>
    </div>

    <UAlert
      v-if="filtered.length === 0"
      title="No matching projects"
      description="Try a different search term or import a project from Printables."
      color="neutral"
      variant="soft"
      icon="i-lucide-search-x"
    />
  </div>
</template>