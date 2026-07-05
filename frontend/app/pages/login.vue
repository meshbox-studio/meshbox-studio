<script setup lang="ts">
import * as z from "zod";
import type { AuthFormField, FormSubmitEvent } from "@nuxt/ui";

definePageMeta({
  layout: false,
});

useHead({
  title: "Login - Meshbox Studio",
});

const { login, isAuthenticated } = useSession();
const router = useRouter();
const toast = useToast();
const pending = ref(false);

const fields: AuthFormField[] = [
  {
    name: "email",
    type: "email",
    label: "Email",
    placeholder: "niklas@example.com",
    required: true,
  },
  {
    name: "password",
    type: "password",
    label: "Password",
    placeholder: "Any password works in demo mode",
    required: true,
  },
  {
    name: "remember",
    type: "checkbox",
    label: "Remember this session",
  },
];

const schema = z.object({
  email: z.string().email("Enter a valid email address"),
  password: z.string().min(1, "Password is required"),
  remember: z.boolean().optional(),
});

type Schema = z.output<typeof schema>;

if (isAuthenticated.value) {
  await navigateTo("/");
}

async function onSubmit(event: FormSubmitEvent<Schema>): Promise<void> {
  pending.value = true;
  const ok = await login(event.data.email);
  pending.value = false;

  if (ok) {
    await router.push("/");
  } else {
    toast.add({
      title: "Login failed",
      description: "Could not reach the server. Is the backend running?",
      color: "error",
      icon: "i-lucide-alert-triangle",
    });
  }
}
</script>

<template>
  <main
    class="relative flex min-h-svh items-center justify-center overflow-hidden px-4 py-8"
  >
    <div class="pointer-events-none absolute inset-0 -z-10">
      <div class="absolute -left-12 top-10 h-56 w-56 rounded-full bg-primary-500/15 blur-3xl" />
      <div class="absolute -right-12 bottom-10 h-64 w-64 rounded-full bg-primary-700/10 blur-3xl" />
    </div>

    <UPageCard class="w-full max-w-md border-default/80 bg-default/95 backdrop-blur">
      <template #title>
        <div class="flex items-center justify-center gap-2">
          <span
            class="inline-flex size-8 items-center justify-center rounded-lg bg-primary-500/20 text-primary-600 dark:text-primary-400"
          >
            <UIcon name="i-lucide-box" class="size-4" />
          </span>
          <span>Meshbox Studio</span>
        </div>
      </template>

      <template #description>
        Sign in to your paperless 3D print workspace.
      </template>

      <UAuthForm
        :schema="schema"
        :fields="fields"
        icon="i-lucide-key-round"
        title="Welcome back"
        description="Demo mode enabled: any email/password pair is accepted."
        :submit="{ label: 'Enter Workspace', icon: 'i-lucide-arrow-right', loading: pending }"
        @submit="onSubmit"
      >
        <template #footer>
          Keep projects, print outcomes, and notes tied together over time.
        </template>
      </UAuthForm>
    </UPageCard>
  </main>
</template>