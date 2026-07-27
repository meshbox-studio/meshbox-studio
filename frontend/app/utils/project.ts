import type { PrintOutcome, ProjectState } from "~/types/meshbox";

/**
 * Hue carries attention, form carries completeness.
 *
 * Amber is the heat signature: it means something is happening right now, and
 * nothing else in the interface is allowed to use it. "Active" is the resting
 * state of most projects, so it stays neutral — a badge that appears on every
 * card teaches the eye to ignore badges.
 *
 * See docs/design-language.md.
 */
export function projectStateColor(
  state: ProjectState,
): "printing" | "neutral" | "error" {
  if (state === "printing") return "printing";
  if (state === "trash") return "error";
  return "neutral";
}

/**
 * Outline reads as unfinished, subtle reads as settled. This is what keeps
 * "Draft" and "Active" apart now that they share a hue.
 */
export function projectStateVariant(
  state: ProjectState,
): "outline" | "subtle" {
  return state === "draft" ? "outline" : "subtle";
}

/** Whether the state warrants a breathing indicator. */
export function projectStateIsLive(state: ProjectState): boolean {
  return state === "printing";
}

export function projectStateLabel(state: ProjectState): string {
  if (state === "printing") return "Printing";
  if (state === "active") return "Active";
  if (state === "draft") return "Draft";
  return "In Trash";
}

/** Queued is information, not a warning — nothing has gone wrong yet. */
export function outcomeColor(state: PrintOutcome): "success" | "info" | "error" {
  if (state === "success") return "success";
  if (state === "queued") return "info";
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