import { createRootRoute } from '@tanstack/react-router'
import type { RouterContext } from '@/router'
import { RootLayout } from '@/features/app-shell/components/RootLayout'

export const Route = createRootRoute<RouterContext>({
  component: RootLayout,
})
