export default defineNuxtConfig({
  compatibilityDate: "2026-07-04",
  ssr: false,
  devServer: {
    port: 5173,
  },
  modules: ["@nuxt/ui", "@nuxt/eslint"],
  css: ["~/assets/css/main.css"],
  ui: {
    theme: {
      // Nuxt UI's defaults plus "printing" — a first-class semantic color for
      // work that is happening right now. See docs/design-language.md.
      colors: [
        "primary",
        "secondary",
        "success",
        "info",
        "warning",
        "error",
        "printing",
      ],
    },
  },
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
