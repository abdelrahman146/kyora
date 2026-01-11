import { createFileRoute } from '@tanstack/react-router'
import { redirectIfAuthenticated } from '@/lib/routeGuards'
import { ForgotPasswordPage } from '@/features/auth/components/ForgotPasswordPage'

export const Route = createFileRoute('/auth/forgot-password')({
  beforeLoad: redirectIfAuthenticated,
  component: ForgotPasswordPage,
})
