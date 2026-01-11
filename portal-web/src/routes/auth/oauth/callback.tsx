import { createFileRoute } from '@tanstack/react-router'
import { OAuthCallbackPage } from '@/features/auth/components/OAuthCallbackPage'

export const Route = createFileRoute('/auth/oauth/callback')({
  component: OAuthCallbackPage,
  validateSearch: (search: Record<string, unknown>) => ({
    code: (search.code as string) || '',
    state: (search.state as string) || '',
    error: (search.error as string) || '',
    error_description: (search.error_description as string) || '',
  }),
})
