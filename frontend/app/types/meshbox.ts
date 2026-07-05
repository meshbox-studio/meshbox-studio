export type ProjectState = "active" | "printing" | "trash" | "draft";

export type PrintOutcome = "success" | "failed" | "queued";

export type FileType = "stl" | "3mf" | "gcode" | "step" | "pdf";

export interface ProjectFile {
  id: string;
  name: string;
  type: FileType;
  size: string;
  updatedAt: string;
}

export interface PrintIteration {
  id: string;
  label: string;
  printer: string;
  material: string;
  startedAt: string;
  duration: string;
  outcome: PrintOutcome;
  notes: string;
}

export interface Source {
  platform: string;
  url: string;
  designer: string;
  designerUrl: string;
  category: string;
  publishedAt: string;
  updatedAt: string;
}

export interface License {
  name: string;
  url: string;
}

export interface PrintProfile {
  printTimeMinutes: number;
  weightGrams: number;
  quantity: number;
  nozzleMm: number;
  layerHeightMm: number;
  material: string;
  printer: string;
}

export interface Project {
  id: string;
  title: string;
  description: string;
  state: ProjectState;
  owner: string;
  lastActivity: string;
  createdAt: string;
  tags: string[];
  files: ProjectFile[];
  iterations: PrintIteration[];
  notes: string[];
  source?: Source | null;
  license?: License | null;
  printProfile?: PrintProfile | null;
  hasThumbnail: boolean;
}

export interface ActivityItem {
  id: string;
  at: string;
  title: string;
  description: string;
  projectId?: string;
}

export interface SidebarStats {
  version: string;
  diskUsed: string;
  diskTotal: string;
  diskPercent: number;
}

export interface SessionUser {
  name: string;
  email: string;
}

export interface DraftProjectInput {
  title: string;
  description: string;
}

export interface EditProjectInput {
  title: string;
  description: string;
  tags: string[];
  source?: Source | null;
  license?: License | null;
}