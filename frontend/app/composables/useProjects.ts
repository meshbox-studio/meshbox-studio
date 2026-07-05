import type { Project, ActivityItem, DraftProjectInput, EditProjectInput } from "~/types/meshbox";
import { apiClient } from "~/api/client";
import { toProject, toProjectList } from "~/utils/api-transforms";

export function useProjects() {
  const projects = useState<Project[]>("projects", () => []);
  const recentActivity = useState<ActivityItem[]>("activity", () => []);
  const loaded = useState<boolean>("projects-loaded", () => false);

  const activeProjects = computed(() =>
    projects.value.filter((project) =>
      ["active", "printing", "draft"].includes(project.state),
    ),
  );

  const trashProjects = computed(() =>
    projects.value.filter((project) => project.state === "trash"),
  );

  function findById(id: string): Project | undefined {
    return projects.value.find((project) => project.id === id);
  }

  async function fetchAll(): Promise<void> {
    const [projectsRes, activityRes] = await Promise.all([
      apiClient.GET("/api/projects"),
      apiClient.GET("/api/activity"),
    ]);
    if (projectsRes.data) {
      projects.value = toProjectList(projectsRes.data);
    }
    if (activityRes.data) {
      recentActivity.value = activityRes.data;
    }
    loaded.value = true;
  }

  async function createDraftProject(input: DraftProjectInput): Promise<Project | null> {
    const { data, error } = await apiClient.POST("/api/projects", {
      body: input,
    });
    if (error || !data) {
      return null;
    }
    const project = toProject(data);
    projects.value = [project, ...projects.value];
    recentActivity.value = [
      {
        id: `activity-${data.id}`,
        at: "Now",
        title: "Created draft project",
        description: `${data.title} was added to the workspace.`,
        projectId: data.id,
      },
      ...recentActivity.value,
    ];
    return project;
  }

  async function importProjectFile(file: File): Promise<Project | null> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await fetch("/api/projects/import", {
      method: "POST",
      body: formData,
      credentials: "include",
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    const project = toProject(data);
    projects.value = [project, ...projects.value];
    recentActivity.value = [
      {
        id: `activity-${project.id}`,
        at: "Now",
        title: "Imported from Printables",
        description: `${project.title} was imported with files and metadata.`,
        projectId: project.id,
      },
      ...recentActivity.value,
    ];
    return project;
  }

  async function updateProject(id: string, input: EditProjectInput): Promise<Project | null> {
    const { data, error } = await apiClient.PUT("/api/projects/{id}", {
      params: { path: { id } },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      body: input as any,
    });
    if (error || !data) {
      return null;
    }
    const updated = toProject(data);
    projects.value = projects.value.map((p) => (p.id === id ? updated : p));
    return updated;
  }

  async function restoreProject(id: string): Promise<boolean> {
    const { error } = await apiClient.PATCH("/api/projects/{id}", {
      params: { path: { id } },
      body: { state: "active" },
    });
    if (error) return false;
    projects.value = projects.value.map((project) =>
      project.id === id
        ? { ...project, state: "active" as const, lastActivity: "just now" }
        : project,
    );
    return true;
  }

  async function moveToTrash(id: string): Promise<boolean> {
    const { error } = await apiClient.PATCH("/api/projects/{id}", {
      params: { path: { id } },
      body: { state: "trash" },
    });
    if (error) return false;
    projects.value = projects.value.map((project) =>
      project.id === id
        ? { ...project, state: "trash" as const, lastActivity: "just now" }
        : project,
    );
    return true;
  }

  async function deleteForever(id: string): Promise<boolean> {
    const { error } = await apiClient.DELETE("/api/projects/{id}", {
      params: { path: { id } },
    });
    if (error) return false;
    projects.value = projects.value.filter((project) => project.id !== id);
    return true;
  }

  return {
    projects,
    activeProjects,
    trashProjects,
    recentActivity,
    loaded,
    fetchAll,
    findById,
    createDraftProject,
    importProjectFile,
    updateProject,
    restoreProject,
    moveToTrash,
    deleteForever,
  };
}