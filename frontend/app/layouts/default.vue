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
    <!-- The sidebar takes no background of its own: it sits directly on the
         canvas (true black in dark mode) and is separated by a hairline only.
         Panels that paint themselves grey are what make an OLED theme look
         like a grey theme. -->
    <UDashboardSidebar
      id="main-sidebar"
      v-model:open="open"
      collapsible
      :ui="{ root: 'bg-transparent' }"
    >
      <template #header="{ collapsed }">
        <NuxtLink
          to="/"
          class="flex items-center gap-2 overflow-hidden rounded-md px-1 py-0.5 font-semibold tracking-tight text-highlighted"
        >
          <span
            class="inline-flex size-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary ring-1 ring-primary/20"
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
          :label="collapsed ? undefined : 'Search'"
          :square="collapsed"
          block
          class="mb-2 justify-start text-muted"
        >
          <template v-if="!collapsed" #trailing>
            <UKbd value="/" class="ms-auto" />
          </template>
        </UButton>

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

      <template #footer="{ collapsed }">
        <SidebarStatusFooter :collapsed="collapsed" />
      </template>
    </UDashboardSidebar>

    <UDashboardPanel id="main-panel">
      <template #header>
        <UDashboardNavbar
          :title="route.meta.title as string || 'Workspace'"
          :ui="{ title: 'text-base font-semibold tracking-tight' }"
        >
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>

          <template #right>
            <ThemeModeButton />
            <WorkspaceUserMenu />
            <!-- One primary action per view, and it is always the rightmost
                 thing in the navbar. It is the only solid accent on screen. -->
            <CreateProjectSlideover />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <slot />
      </template>
    </UDashboardPanel>
  </UDashboardGroup>
</template>
