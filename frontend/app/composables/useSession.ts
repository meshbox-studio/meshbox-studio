import type { SessionUser } from "~/types/meshbox";
import { apiClient } from "~/api/client";

const fallbackUser: SessionUser = {
  name: "Guest",
  email: "",
};

export function useSession() {
  const user = useState<SessionUser>("session-user", () => fallbackUser);
  const isAuthenticated = useState<boolean>("session-authenticated", () => false);

  async function login(email: string): Promise<boolean> {
    try {
      const { data, error } = await apiClient.POST("/api/auth/login", {
        body: { email },
      });
      if (error || !data) {
        return false;
      }
      user.value = data.user;
      isAuthenticated.value = true;
      return true;
    } catch {
      return false;
    }
  }

  async function logout(): Promise<void> {
    try {
      await apiClient.POST("/api/auth/logout");
    } finally {
      user.value = fallbackUser;
      isAuthenticated.value = false;
    }
  }

  async function hydrate(): Promise<void> {
    try {
      const { data, error } = await apiClient.GET("/api/auth/session");
      if (error || !data) {
        isAuthenticated.value = false;
        return;
      }
      user.value = data.user;
      isAuthenticated.value = true;
    } catch {
      isAuthenticated.value = false;
    }
  }

  return {
    user,
    isAuthenticated,
    login,
    logout,
    hydrate,
  };
}