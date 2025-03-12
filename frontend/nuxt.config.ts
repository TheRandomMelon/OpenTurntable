// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  ssr: false,
  router: {
    options: {
      hashMode: true
    }
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: ['@vesp/nuxt-fontawesome'],
  fontawesome: {
    icons: {
      solid: ["file-audio", "book", "circle-play", "bars"],
      regular: [],
			brands: []
    },
    component: 'fa',
  }
})