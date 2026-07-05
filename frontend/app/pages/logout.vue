<script setup lang="ts">
definePageMeta({
  layout: false,
});

const { logout } = useSession();
const router = useRouter();

async function doLogout(): Promise<void> {
  await logout();
  void router.push("/login");
}
</script>

<template>
  <main class="flex min-h-svh items-center justify-center px-4 py-8">
    <UPageCard
      title="Sign out"
      description="This clears your local session and returns you to login."
      class="w-full max-w-md"
    >
      <div class="space-y-4">
        <UAlert
          color="warning"
          icon="i-lucide-log-out"
          title="Ready to leave the workspace?"
          description="You can sign back in immediately from the same machine."
        />

        <div class="flex items-center justify-end gap-2">
          <UButton
            label="Cancel"
            color="neutral"
            variant="outline"
            to="/"
          />
          <UButton
            label="Logout"
            color="error"
            icon="i-lucide-log-out"
            @click="doLogout"
          />
        </div>
      </div>
    </UPageCard>
  </main>
</template>
