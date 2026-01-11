import { createFileRoute } from '@tanstack/react-router'
import { OnboardingRoot } from '@/features/onboarding/components/OnboardingRoot'

export const Route = createFileRoute('/onboarding/')({
  component: OnboardingRoot,
})
