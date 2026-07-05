<script setup lang="ts">
import { fileIconName, outcomeColor, outcomeLabel, printTimeLabel, projectStateColor, projectStateLabel } from "~/utils/project";
import type { EditProjectInput } from "~/types/meshbox";

definePageMeta({
  title: "Project",
});

const route = useRoute();
const toast = useToast();
const { findById, loaded, fetchAll, updateProject } = useProjects();

if (!loaded.value) {
  await fetchAll();
}

const project = computed(() => findById(String(route.params.id)));

if (!project.value) {
  throw createError({
    statusCode: 404,
    statusMessage: "Project not found",
  });
}

useHead(() => ({
  title: `${project.value?.title ?? "Project"} - Meshbox Studio`,
}));

const editing = ref(false);
const editForm = ref<EditProjectInput>({
  title: "",
  description: "",
  tags: [],
  source: null,
  license: null,
});
const saving = ref(false);

function startEdit(): void {
  if (!project.value) return;
  editForm.value = {
    title: project.value.title,
    description: project.value.description,
    tags: [...project.value.tags],
    source: project.value.source ? { ...project.value.source } : null,
    license: project.value.license ? { ...project.value.license } : null,
  };
  editing.value = true;
}

async function saveEdit(): Promise<void> {
  if (!project.value) return;
  saving.value = true;

  const updated = await updateProject(project.value.id, editForm.value);
  saving.value = false;

  if (updated) {
    editing.value = false;
    toast.add({
      title: "Project updated",
      color: "success",
      icon: "i-lucide-check-circle-2",
    });
  } else {
    toast.add({
      title: "Failed to update",
      color: "error",
      icon: "i-lucide-alert-triangle",
    });
  }
}

function cancelEdit(): void {
  editing.value = false;
}

function updateTagInput(val: string): void {
  editForm.value.tags = val.split(/[, ]+/).filter(Boolean);
}

function thumbnailUrl(id: string): string {
  return `/api/projects/${id}/thumbnail`;
}

function fileDownloadUrl(id: string, fileId: string): string {
  return `/api/projects/${id}/files/${fileId}`;
}
</script>

<template>
  <div v-if="project" class="space-y-4">
    <UCard>
      <template #header>
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="flex items-start gap-4">
            <img
              v-if="project.hasThumbnail"
              :src="thumbnailUrl(project.id)"
              class="size-24 rounded-lg border border-default object-cover"
              alt="Thumbnail"
            >
            <div>
              <h1 class="text-xl font-semibold text-highlighted">{{ project.title }}</h1>
              <p class="mt-1 text-sm text-muted">{{ project.description }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <UBadge
              :color="projectStateColor(project.state)"
              variant="soft"
              :label="projectStateLabel(project.state)"
            />
            <UButton
              icon="i-lucide-pencil"
              color="neutral"
              variant="ghost"
              size="xs"
              @click="startEdit"
            />
          </div>
        </div>
      </template>

      <div class="space-y-3">
        <div v-if="project.source" class="grid gap-3 md:grid-cols-3">
          <div v-if="project.source.designer" class="rounded-lg bg-muted/60 p-3">
            <p class="text-xs text-muted">Designer</p>
            <a
              v-if="project.source.designerUrl"
              :href="project.source.designerUrl"
              target="_blank"
              class="text-sm font-medium text-primary hover:underline"
            >
              {{ project.source.designer }}
            </a>
            <p v-else class="text-sm font-medium text-highlighted">{{ project.source.designer }}</p>
          </div>
          <div v-if="project.source.url" class="rounded-lg bg-muted/60 p-3">
            <p class="text-xs text-muted">Source</p>
            <a
              :href="project.source.url"
              target="_blank"
              class="text-sm font-medium text-primary hover:underline"
            >
              View on Printables
            </a>
          </div>
          <div v-if="project.license" class="rounded-lg bg-muted/60 p-3">
            <p class="text-xs text-muted">License</p>
            <a
              v-if="project.license.url"
              :href="project.license.url"
              target="_blank"
              class="text-sm font-medium text-primary hover:underline"
            >
              {{ project.license.name || project.license.url }}
            </a>
            <p v-else class="text-sm font-medium text-highlighted">{{ project.license.name }}</p>
          </div>
        </div>

        <div v-if="project.printProfile" class="grid gap-2 grid-cols-4 md:grid-cols-7">
          <div v-if="project.printProfile.printTimeMinutes" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Print time</p>
            <p class="text-sm font-medium text-highlighted">{{ printTimeLabel(project.printProfile.printTimeMinutes) }}</p>
          </div>
          <div v-if="project.printProfile.weightGrams" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Weight</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.weightGrams }}g</p>
          </div>
          <div v-if="project.printProfile.quantity" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Qty</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.quantity }}</p>
          </div>
          <div v-if="project.printProfile.nozzleMm" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Nozzle</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.nozzleMm }}mm</p>
          </div>
          <div v-if="project.printProfile.layerHeightMm" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Layer</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.layerHeightMm }}mm</p>
          </div>
          <div v-if="project.printProfile.material" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Material</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.material }}</p>
          </div>
          <div v-if="project.printProfile.printer" class="rounded-md bg-muted/60 p-2 text-center">
            <p class="text-xs text-muted">Printer</p>
            <p class="text-sm font-medium text-highlighted">{{ project.printProfile.printer }}</p>
          </div>
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
      </div>
    </UCard>

    <UCard v-if="editing">
      <template #header>
        <h2 class="text-base font-semibold text-highlighted">Edit Project</h2>
      </template>
      <div class="space-y-3">
        <UFormField label="Title">
          <UInput v-model="editForm.title" />
        </UFormField>
        <UFormField label="Description">
          <UTextarea v-model="editForm.description" :rows="3" />
        </UFormField>
        <UFormField label="Tags (comma or space separated)">
          <UInput
            :model-value="editForm.tags.join(' ')"
            placeholder="lighting functional v2"
            @update:model-value="updateTagInput"
          />
        </UFormField>
        <div class="grid gap-3 md:grid-cols-2">
          <UFormField label="Source URL">
            <UInput v-model="editForm.source!.url" />
          </UFormField>
          <UFormField label="Designer">
            <UInput v-model="editForm.source!.designer" />
          </UFormField>
          <UFormField label="License name">
            <UInput v-model="editForm.license!.name" />
          </UFormField>
          <UFormField label="License URL">
            <UInput v-model="editForm.license!.url" />
          </UFormField>
        </div>
        <div class="flex justify-end gap-2">
          <UButton label="Cancel" color="neutral" variant="outline" @click="cancelEdit" />
          <UButton label="Save" color="primary" :loading="saving" @click="saveEdit" />
        </div>
      </div>
    </UCard>

    <div class="grid gap-4 xl:grid-cols-2">
      <UCard>
        <template #header>
          <h2 class="text-base font-semibold text-highlighted">Linked files</h2>
        </template>
        <ul class="space-y-2">
          <li
            v-for="file in project.files"
            :key="file.id"
            class="flex items-center justify-between gap-3 rounded-lg border border-default px-3 py-2"
          >
            <div class="min-w-0">
              <a
                :href="fileDownloadUrl(project.id, file.id)"
                class="truncate text-sm font-medium text-highlighted hover:text-primary"
              >
                <UIcon :name="fileIconName(file.type)" class="mr-1 inline size-4" />
                {{ file.name }}
              </a>
              <p class="text-xs text-muted">Updated {{ file.updatedAt }}</p>
            </div>
            <UBadge color="neutral" variant="soft" :label="file.size" />
          </li>
        </ul>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="text-base font-semibold text-highlighted">Project notes</h2>
        </template>

        <ul class="space-y-2">
          <li
            v-for="(note, index) in project.notes"
            :key="`${project.id}-note-${index}`"
            class="rounded-lg border border-default px-3 py-2 text-sm text-muted"
          >
            {{ note }}
          </li>
        </ul>
      </UCard>
    </div>

    <UCard>
      <template #header>
        <h2 class="text-base font-semibold text-highlighted">Print iterations</h2>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-default text-sm">
          <thead>
            <tr class="text-left text-xs uppercase tracking-wide text-muted">
              <th class="px-2 py-2">Run</th>
              <th class="px-2 py-2">Printer</th>
              <th class="px-2 py-2">Material</th>
              <th class="px-2 py-2">Duration</th>
              <th class="px-2 py-2">Outcome</th>
              <th class="px-2 py-2">Notes</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-default">
            <tr v-for="iteration in project.iterations" :key="iteration.id" class="align-top">
              <td class="px-2 py-2 font-medium text-highlighted">{{ iteration.label }}</td>
              <td class="px-2 py-2 text-muted">{{ iteration.printer }}</td>
              <td class="px-2 py-2 text-muted">{{ iteration.material }}</td>
              <td class="px-2 py-2 text-muted">{{ iteration.duration }}</td>
              <td class="px-2 py-2">
                <UBadge
                  :color="outcomeColor(iteration.outcome)"
                  variant="soft"
                  :label="outcomeLabel(iteration.outcome)"
                />
              </td>
              <td class="px-2 py-2 text-muted">{{ iteration.notes }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>