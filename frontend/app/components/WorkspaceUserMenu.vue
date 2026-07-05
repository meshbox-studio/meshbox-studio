<script setup lang="ts">
const { user } = useSession();

const items = computed(() => [
  [
    {
      label: "Settings",
      icon: "i-lucide-settings",
      to: "/settings",
    },
  ],
  [
    {
      label: "Logout",
      icon: "i-lucide-log-out",
      to: "/logout",
    },
  ],
]);
</script>

<template>
  <UDropdownMenu :items="items" :content="{ align: 'end' }" :ui="{ content: 'w-40' }">
    <UButton
      color="neutral"
      variant="ghost"
      class="rounded-full"
      :avatar="{
        alt: user.name,
      }"
      :label="user.name"
      trailing-icon="i-lucide-chevrons-up-down"
      :ui="{ label: 'hidden xl:inline', trailingIcon: 'hidden xl:inline' }"
    />
  </UDropdownMenu>
</template>
