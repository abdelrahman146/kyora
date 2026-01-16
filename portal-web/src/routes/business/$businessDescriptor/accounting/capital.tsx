import { createFileRoute } from '@tanstack/react-router'

/**
 * Capital Management Route (Placeholder)
 *
 * Full implementation will be added in Step 5.
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/capital',
)({
  staticData: {
    titleKey: 'pages.capital',
  },
  component: CapitalPlaceholder,
})

function CapitalPlaceholder() {
  return (
    <div className="flex items-center justify-center min-h-[50vh]">
      <p className="text-base-content/60">Capital management coming soon...</p>
    </div>
  )
}
