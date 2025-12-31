import { defineConfig } from 'vite'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import viteReact from '@vitejs/plugin-react'
import viteTsConfigPaths from 'vite-tsconfig-paths'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  plugins: [
    // TanStack Router plugin for file-based routing (client-only mode)
    TanStackRouterVite({
      routesDirectory: './src/routes',
      generatedRouteTree: './src/routeTree.gen.ts',
    }),
    // Path aliases from tsconfig
    viteTsConfigPaths({
      projects: ['./tsconfig.json'],
    }),
    // Tailwind CSS v4
    tailwindcss(),
    // React plugin
    viteReact(),
  ],
  build: {
    rollupOptions: {
      // Exclude TanStack Store devtools from production bundle
      external: (id) => {
        if (
          process.env.NODE_ENV === 'production' &&
          id.includes('@tanstack/store-devtools')
        ) {
          return true
        }
        return false
      },
    },
  },
})

export default config
