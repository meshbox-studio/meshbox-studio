<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui";

const route = useRoute();

const open = ref(true);

const links = computed<NavigationMenuItem[][]>(() => [
  [
    {
      label: "Overview",
      icon: "i-lucide-layout-dashboard",
      to: "/",
      active: route.path === "/",
      onSelect: () => {
        open.value = false;
      },
    },
    {
      label: "Projects",
      icon: "i-lucide-folder",
      to: "/projects",
      active: route.path.startsWith("/projects"),
      onSelect: () => {
        open.value = false;
      },
    },
    {
      label: "Trash",
      icon: "i-lucide-trash-2",
      to: "/trash",
      active: route.path.startsWith("/trash"),
      onSelect: () => {
        open.value = false;
      },
    },
  ],
  [
    {
      label: "Documentation",
      icon: "i-lucide-book-open",
      to: "https://meshbox.studio",
      target: "_blank",
    },
    {
      label: "GitHub",
      icon: "i-lucide-github",
      to: "https://github.com/meshbox-studio/meshbox-studio",
      target: "_blank",
    },
  ],
]);
</script>

<template>
  <UDashboardGroup storage="local" storage-key="meshbox-dashboard">
    <UDashboardSidebar
      id="main-sidebar"
      v-model:open="open"
      collapsible
      class="bg-default/90 backdrop-blur"
    >
      <template #header="{ collapsed }">
        <NuxtLink
          to="/"
          class="flex items-center gap-2 overflow-hidden rounded-lg px-1 py-0.5 text-lg font-semibold text-highlighted"
        >
          <span
            class="inline-flex size-7 shrink-0 items-center justify-center rounded-md bg-primary-500/20 text-primary-600 dark:text-primary-400"
          >
            <UIcon name="i-lucide-box" class="size-4" />
          </span>
          <span v-if="!collapsed" class="truncate">Meshbox Studio</span>
        </NuxtLink>
      </template>

      <template #default="{ collapsed }">
        <UButton
          color="neutral"
          variant="soft"
          icon="i-lucide-search"
          :label="collapsed ? undefined : 'Search projects, notes, materials'"
          :square="collapsed"
          block
          class="mb-2"
        />

        <UNavigationMenu
          :items="links[0]"
          :collapsed="collapsed"
          orientation="vertical"
          tooltip
          highlight
        />

        <UNavigationMenu
          :items="links[1]"
          :collapsed="collapsed"
          orientation="vertical"
          tooltip
          class="mt-auto"
        />
      </template>
    </UDashboardSidebar>

    <UDashboardPanel id="main-panel" :ui="{ body: 'flex flex-col gap-4 sm:gap-6 flex-1 overflow-y-auto px-4 py-4 sm:px-6 sm:py-6' }">
      <template #header>
        <UDashboardNavbar :title="route.meta.title as string || 'Workspace'">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>

          <template #right>
            <CreateProjectSlideover />
            <ThemeModeButton />
            <WorkspaceUserMenu />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <slot />
      </template>
    </UDashboardPanel>
  </UDashboardGroup>
</template>
