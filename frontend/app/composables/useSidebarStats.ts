import type { SidebarStats } from "~/types/meshbox";
import { apiClient } from "~/api/client";

const fallbackStats: SidebarStats = {
  version: "dev",
  diskUsed: "0 B",
  diskTotal: "0 B",
  diskPercent: 0,
};

function normalizePercent(value: unknown): number {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    return 0;
  }
  return Math.min(100, Math.max(0, Math.round(parsed)));
}

export function useSidebarStats() {
  const stats = useState<SidebarStats>("sidebar-stats", () => fallbackStats);
  const loading = useState<boolean>("sidebar-stats-loading", () => false);

  async function refreshStats(): Promise<void> {
    loading.value = true;

    try {
      const { data } = await apiClient.GET("/api/stats");
      if (data) {
        stats.value = {
          version: data.version ?? fallbackStats.version,
          diskUsed: data.diskUsed ?? fallbackStats.diskUsed,
          diskTotal: data.diskTotal ?? fallbackStats.diskTotal,
          diskPercent: normalizePercent(data.diskPercent),
        };
      }
    } catch {
      stats.value = fallbackStats;
    } finally {
      loading.value = false;
    }
  }

  return {
    stats,
    loading,
    refreshStats,
  };
}