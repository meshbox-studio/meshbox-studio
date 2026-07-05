export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated } = useSession();

  if (to.path === "/login") {
    if (isAuthenticated.value) {
      return navigateTo("/");
    }
    return;
  }

  if (to.path === "/logout") {
    return;
  }

  if (!isAuthenticated.value) {
    return navigateTo("/login");
  }
});