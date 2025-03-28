// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  css: ['~/assets/css/main.css'],
  ssr: false,
  imports: {
    dirs: [
      'wailsjs/runtime/**',
    ],
  },
  modules: ['@vesp/nuxt-fontawesome', '@pinia/nuxt', '@nuxtjs/i18n'],
  pinia: {
      storesDirs: ['./stores/**'],
  },
  i18n: {
    vueI18n: './i18n.config.ts',
    locales: [
      {
        code: 'en',
        file: 'en_US.json'
      }
    ],
    defaultLocale: 'en',
    langDir: 'lang/',
    lazy: true
  },
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

  fontawesome: {
    icons: {
      solid: ["file-audio", "book", "circle-play", "circle-pause", "bars", "volume-xmark", "volume-off", "volume-low", "volume-high", "caret-up", "caret-down", "backward-step", "forward-step", "shuffle", "repeat", "gauge", "square-plus", "folder-plus"],
      regular: [],
            brands: []
    },
    component: 'fa',
  },

  app: {
    head: {
      charset: "utf-8",
      viewport: "width=device-width, initial-scale=1",
    }
  }
})