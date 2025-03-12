// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  ssr: false,
  imports: {
    dirs: [
      'wailsjs/runtime/**',
    ],
  },
  modules: [
    '@vesp/nuxt-fontawesome',
    '@pinia/nuxt'
  ],
	pinia: {
		storesDirs: ['./stores/**'],
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
      solid: ["file-audio", "book", "circle-play", "circle-pause", "bars", "volume-xmark", "volume-off", "volume-low", "volume-high"],
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