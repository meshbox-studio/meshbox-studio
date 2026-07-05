<script setup lang="ts">
definePageMeta({
  title: "Settings",
});

const { user } = useSession();

const profile = reactive({
  name: user.value.name,
  email: user.value.email,
});

const toast = useToast();

function saveProfile(): void {
  user.value = {
    name: profile.name,
    email: profile.email,
  };

  toast.add({
    title: "Profile updated",
    description: "Your workspace identity settings were saved.",
    color: "success",
    icon: "i-lucide-check",
  });
}
</script>

<template>
  <UCard title="Profile" description="Public identity and contact details for this workspace.">
    <form class="space-y-4" @submit.prevent="saveProfile">
      <div class="grid gap-4 md:grid-cols-2">
        <UFormField label="Display name">
          <UInput v-model="profile.name" icon="i-lucide-user-round" />
        </UFormField>

        <UFormField label="Email">
          <UInput v-model="profile.email" type="email" icon="i-lucide-mail" />
        </UFormField>
      </div>

      <div class="flex justify-end">
        <UButton type="submit" label="Save profile" icon="i-lucide-save" />
      </div>
    </form>
  </UCard>
</template>