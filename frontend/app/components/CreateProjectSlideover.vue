<script setup lang="ts">
import type { DraftProjectInput } from "~/types/meshbox";

const { createDraftProject, importProjectFile } = useProjects();

const router = useRouter();
const toast = useToast();

const open = ref(false);
const mode = ref<"import" | "manual">("import");

const form = reactive<DraftProjectInput>({
  title: "",
  description: "",
});

const files = ref<File[]>([]);
const importing = ref(false);
const importResults = ref<{ name: string; success: boolean; error?: string }[]>([]);

function resetForm(): void {
  form.title = "";
  form.description = "";
  files.value = [];
  importResults.value = [];
}

function onFileChange(event: Event): void {
  const input = event.target as HTMLInputElement;
  files.value = input.files ? Array.from(input.files) : [];
  importResults.value = [];
}

function removeFile(index: number): void {
  files.value.splice(index, 1);
}

async function importFiles(): Promise<void> {
  if (files.value.length === 0) return;

  importing.value = true;
  importResults.value = [];

  for (const file of files.value) {
    const project = await importProjectFile(file);
    if (project) {
      importResults.value.push({ name: file.name, success: true });
    } else {
      importResults.value.push({
        name: file.name,
        success: false,
        error: "Not a valid Printables export",
      });
    }
  }

  importing.value = false;

  const successful = importResults.value.filter((r) => r.success);
  const failed = importResults.value.filter((r) => !r.success);

  if (successful.length > 0) {
    toast.add({
      title: `${successful.length} project${successful.length > 1 ? "s" : ""} imported`,
      description: successful.map((r) => r.name).join(", "),
      color: "success",
      icon: "i-lucide-check-circle-2",
    });
  }

  if (failed.length > 0) {
    toast.add({
      title: `${failed.length} file${failed.length > 1 ? "s" : ""} rejected`,
      description: "These files don't contain a recognizable Printables PDF.",
      color: "error",
      icon: "i-lucide-alert-triangle",
    });
  }

  if (successful.length > 0) {
    open.value = false;
    resetForm();
    if (successful.length === 1) {
      void router.push("/projects");
    } else {
      void router.push("/projects");
    }
  }
}

async function submitManual(): Promise<void> {
  if (!form.title.trim()) {
    toast.add({
      title: "Project name required",
      description: "Give your project a name.",
      color: "error",
      icon: "i-lucide-alert-triangle",
    });
    return;
  }

  const project = await createDraftProject({
    title: form.title.trim(),
    description: form.description.trim() || "New draft project.",
  });

  if (!project) {
    toast.add({
      title: "Failed to create",
      description: "Could not reach the server.",
      color: "error",
    });
    return;
  }

  open.value = false;

  toast.add({
    title: "Project created",
    description: `${project.title} is ready.`,
    color: "success",
    icon: "i-lucide-check-circle-2",
  });

  resetForm();
  void router.push(`/projects/${project.id}`);
}
</script>

<template>
  <USlideover
    v-model:open="open"
    title="New Project"
    description="Import from Printables or start a blank project."
    side="right"
    :ui="{ footer: 'justify-end' }"
  >
    <UButton
      label="New Project"
      icon="i-lucide-folder-plus"
      color="primary"
      variant="solid"
    />

    <template #body>
      <div class="space-y-4">
        <UTabs
          v-model="mode"
          :items="[
            { label: 'Import from Printables', value: 'import' },
            { label: 'Blank project', value: 'manual' },
          ]"
          class="w-full"
        />

        <div v-if="mode === 'import'" class="space-y-4">
          <UFileUpload
            label="Drop Printables .zip files here"
            accept=".zip"
            multiple
            @change="onFileChange"
          />

          <div v-if="files.length > 0" class="space-y-2">
            <div
              v-for="(file, index) in files"
              :key="file.name"
              class="flex items-center justify-between gap-2 rounded-lg border border-default px-3 py-2"
            >
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-highlighted">{{ file.name }}</p>
                <p class="text-xs text-muted">{{ (file.size / 1024).toFixed(1) }} KB</p>
              </div>
              <UBadge
                v-if="importResults[index]?.success"
                color="success"
                variant="soft"
                label="Imported"
              />
              <UBadge
                v-else-if="importResults[index] && !importResults[index].success"
                color="error"
                variant="soft"
                label="Failed"
              />
              <UButton
                v-else
                color="neutral"
                variant="ghost"
                icon="i-lucide-x"
                size="xs"
                @click="removeFile(index)"
              />
            </div>
          </div>
        </div>

        <form v-else class="space-y-4" @submit.prevent="submitManual">
          <UFormField label="Project name" required>
            <UInput
              v-model="form.title"
              placeholder="Hex Lamp V3"
              icon="i-lucide-pencil-line"
              size="lg"
            />
          </UFormField>

          <UFormField label="Description">
            <UTextarea
              v-model="form.description"
              :rows="4"
              placeholder="What are you building?"
            />
          </UFormField>
        </form>
      </div>
    </template>

    <template #footer>
      <UButton
        label="Cancel"
        color="neutral"
        variant="outline"
        @click="open = false"
      />
      <UButton
        v-if="mode === 'import'"
        label="Import"
        color="primary"
        icon="i-lucide-upload"
        :loading="importing"
        :disabled="files.length === 0"
        @click="importFiles"
      />
      <UButton
        v-else
        label="Create"
        color="primary"
        icon="i-lucide-plus"
        @click="submitManual"
      />
    </template>
  </USlideover>
</template>