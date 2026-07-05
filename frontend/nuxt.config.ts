export default defineNuxtConfig({
  compatibilityDate: "2026-07-04",
  ssr: false,
  devServer: {
    port: 5173,
  },
  modules: ["@nuxt/ui", "@nuxt/eslint"],
  css: ["~/assets/css/main.css"],
  ui: {},
  nitro: {
    output: {
      publicDir: "../backend/internal/webui/dist/app",
    },
  },
  vite: {
    server: {
      proxy: {
        "/api": {
          target: "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },
  },
});
