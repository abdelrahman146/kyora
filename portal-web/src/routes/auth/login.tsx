import { createFileRoute } from '@tanstack/react-router'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { LoginPage } from '@/features/auth/components/LoginPage'

export const Route = createFileRoute('/auth/login')({
  beforeLoad: redirectIfAuthenticated,
  component: LoginPage,
  validateSearch: (search: Record<string, unknown>) => ({
    redirect: (search.redirect as string) || '/',
  }),
})
