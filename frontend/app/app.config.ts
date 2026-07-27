export default defineAppConfig({
  ui: {
    colors: {
      // Cool accent for everything interactive. The content in this app is
      // photographed plastic — mostly warm — so a cool chrome pushes it
      // forward rather than blending into it.
      primary: "cyan",
      secondary: "sky",

      // Warm hues are reserved for heat and process, never for chrome.
      printing: "amber",
      warning: "orange",
      success: "emerald",
      error: "rose",
      info: "blue",

      neutral: "graphite",
    },
  },
});
