import { createFileRoute } from '@tanstack/react-router'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { ResetPasswordPage } from '@/features/auth/components/ResetPasswordPage'

export const Route = createFileRoute('/auth/reset-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ResetPasswordPage,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
  }),
})
