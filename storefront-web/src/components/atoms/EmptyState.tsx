import { memo, type ReactNode } from 'react';
import { PackageOpen } from 'lucide-react';

interface EmptyStateProps {
  icon?: ReactNode;
  title: string;
  description?: string;
  action?: ReactNode;
}

/**
 * EmptyState Atom - Display when no data is available
 * Memoized to prevent unnecessary re-renders
 * Follows KDS principles with clean, friendly design
 */
export const EmptyState = memo<EmptyStateProps>(function EmptyState({ icon, title, description, action }) {
  return (
    <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
      {/* Icon */}
      <div className="mb-4 text-neutral-400">
        {icon || <PackageOpen className="w-16 h-16" strokeWidth={1.5} />}
      </div>

      {/* Title */}
      <h3 className="text-lg font-semibold text-neutral-900 mb-2">{title}</h3>

      {/* Description */}
      {description && (
        <p className="text-sm text-neutral-500 mb-6 max-w-sm">{description}</p>
      )}

      {/* Action */}
      {action && <div>{action}</div>}
    </div>
  );
});
