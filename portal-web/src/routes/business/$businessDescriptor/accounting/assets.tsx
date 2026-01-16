import { createFileRoute } from '@tanstack/react-router'

/**
 * Assets Management Route (Placeholder)
 *
 * Full implementation will be added in Step 6.
 */
export const Route = createFileRoute(
  '/business/$businessDescriptor/accounting/assets',
)({
  staticData: {
    titleKey: 'pages.assets',
  },
  component: AssetsPlaceholder,
})

function AssetsPlaceholder() {
  return (
    <div className="flex items-center justify-center min-h-[50vh]">
      <p className="text-base-content/60">Assets management coming soon...</p>
    </div>
  )
}
