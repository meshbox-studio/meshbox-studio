<script setup lang="ts">
definePageMeta({
  title: "Trash",
});

const { trashProjects, restoreProject, deleteForever, loaded, fetchAll } = useProjects();
const toast = useToast();

if (!loaded.value) {
  await fetchAll();
}

async function restore(id: string): Promise<void> {
  const ok = await restoreProject(id);
  if (!ok) {
    toast.add({
      title: "Failed to restore",
      description: "Could not restore the project.",
      color: "error",
    });
  }
}

async function remove(id: string): Promise<void> {
  const ok = await deleteForever(id);
  if (!ok) {
    toast.add({
      title: "Failed to delete",
      description: "Could not delete the project.",
      color: "error",
    });
  }
}
</script>

<template>
  <div class="space-y-4">
    <UAlert
      title="Trash holds removable items"
      description="Restore when needed, or delete forever to remove all project metadata."
      icon="i-lucide-trash-2"
      color="warning"
      variant="soft"
    />

    <UCard v-for="project in trashProjects" :key="project.id">
      <template #header>
        <div>
          <h2 class="text-base font-semibold text-highlighted">{{ project.title }}</h2>
          <p class="text-sm text-muted">{{ project.description }}</p>
        </div>
      </template>

      <div class="grid gap-3 md:grid-cols-3">
        <div class="rounded-md bg-muted/60 p-2 text-sm">
          <p class="text-xs text-muted">Files</p>
          <p class="font-medium text-highlighted">{{ project.files.length }}</p>
        </div>
        <div class="rounded-md bg-muted/60 p-2 text-sm">
          <p class="text-xs text-muted">Iterations</p>
          <p class="font-medium text-highlighted">{{ project.iterations.length }}</p>
        </div>
        <div class="rounded-md bg-muted/60 p-2 text-sm">
          <p class="text-xs text-muted">Last activity</p>
          <p class="font-medium text-highlighted">{{ project.lastActivity }}</p>
        </div>
      </div>

      <template #footer>
        <div class="flex flex-wrap items-center justify-end gap-2">
          <UButton
            label="Restore"
            color="success"
            variant="soft"
            icon="i-lucide-rotate-ccw"
            @click="restore(project.id)"
          />
          <UButton
            label="Delete forever"
            color="error"
            variant="outline"
            icon="i-lucide-x"
            @click="remove(project.id)"
          />
        </div>
      </template>
    </UCard>

    <UAlert
      v-if="trashProjects.length === 0"
      title="Trash is empty"
      description="Projects moved to trash will appear here."
      color="success"
      variant="soft"
      icon="i-lucide-check"
    />
  </div>
</template>