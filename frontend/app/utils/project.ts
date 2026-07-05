import type { PrintOutcome, ProjectState } from "~/types/meshbox";

export function projectStateColor(
  state: ProjectState,
): "primary" | "success" | "warning" | "neutral" | "error" {
  if (state === "printing") return "warning";
  if (state === "active") return "primary";
  if (state === "draft") return "success";
  return "error";
}

export function projectStateLabel(state: ProjectState): string {
  if (state === "printing") return "Printing";
  if (state === "active") return "Active";
  if (state === "draft") return "Draft";
  return "In Trash";
}

export function outcomeColor(state: PrintOutcome): "success" | "warning" | "error" {
  if (state === "success") return "success";
  if (state === "queued") return "warning";
  return "error";
}

export function outcomeLabel(state: PrintOutcome): string {
  if (state === "success") return "Success";
  if (state === "queued") return "Queued";
  return "Failed";
}

export function fileIconName(fileType: "stl" | "3mf" | "gcode" | "step" | "pdf"): string {
  if (fileType === "gcode") return "i-lucide-file-code-2";
  if (fileType === "3mf") return "i-lucide-package";
  if (fileType === "step") return "i-lucide-cuboid";
  if (fileType === "pdf") return "i-lucide-file-text";
  return "i-lucide-file";
}

export function printTimeLabel(minutes: number): string {
  if (minutes < 60) return `${minutes}m`;
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}