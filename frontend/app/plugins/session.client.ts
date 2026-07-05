export default defineNuxtPlugin(async () => {
  const { hydrate } = useSession();
  await hydrate();
});