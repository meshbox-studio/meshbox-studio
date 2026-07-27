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
  <!-- No blurred colour blobs: a glow behind a card is exactly the thing that
       stops an OLED panel from switching its pixels off. The card lifts off the
       black on its own, via surface and hairline. -->
  <main class="flex min-h-svh items-center justify-center px-4 py-8">
    <UPageCard class="w-full max-w-md">
      <template #title>
        <div class="flex items-center justify-center gap-2">
          <span
            class="inline-flex size-8 items-center justify-center rounded-md bg-primary/10 text-primary ring-1 ring-primary/20"
          >
            <UIcon name="i-lucide-box" class="size-4" />
          </span>
          <span class="tracking-tight">Meshbox Studio</span>
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