import createClient from "openapi-fetch";
import type { paths } from "~/api/generated/schema.d";

export const apiClient = createClient<paths>({ credentials: "include" });