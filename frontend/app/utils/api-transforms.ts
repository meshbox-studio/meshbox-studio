import type { Project } from "~/types/meshbox";
import type { components } from "~/api/generated/schema.d";

type SchemaProject = components["schemas"]["Project"];

export function toProject(raw: SchemaProject): Project {
  const { $schema: _, ...rest } = raw;
  return {
    ...rest,
    tags: rest.tags ?? [],
    files: rest.files ?? [],
    iterations: rest.iterations ?? [],
    notes: rest.notes ?? [],
    hasThumbnail: rest.hasThumbnail ?? false,
    source: rest.source ?? null,
    license: rest.license ?? null,
    printProfile: rest.printProfile ?? null,
  } as unknown as Project;
}

export function toProjectList(raw: SchemaProject[]): Project[] {
  return raw.map(toProject);
}