import { useTranslation } from 'react-i18next'
import { Link, useLocation } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { ChevronRight } from 'lucide-react'
import { Badge } from '../atoms/Badge'
import type {
  NavConfigItem,
  NavGroupConfig,
  NavItemConfig,
} from '@/types/navigation'
import { businessStore, toggleNavGroup } from '@/stores/businessStore'
import { cn } from '@/lib/utils'
import { useLanguage } from '@/hooks/useLanguage'

export interface NavGroupProps {
  /** Navigation config item (item, group, or separator) */
  config: NavConfigItem
  /** Base path for route construction (e.g., /business/:id) */
  basePath: string
  /** Whether sidebar is collapsed (icon-only mode) */
  sidebarCollapsed: boolean
  /** Callback when item is clicked (for mobile drawer close) */
  onItemClick?: () => void
}

/**
 * NavGroup Component
 *
 * Renders navigation items, groups with collapsible behavior, or section separators.
 *
 * Features:
 * - Top-level navigation items
 * - Collapsible groups with chevron rotation
 * - Section separators with divider
 * - Active route highlighting
 * - Badge support ("NEW" indicators)
 * - RTL support
 * - Desktop collapsed mode (icon-only)
 */
export function NavGroup({
  config,
  basePath,
  sidebarCollapsed,
  onItemClick,
}: NavGroupProps) {
  // Handle section separator
  if ('type' in config && config.type === 'separator') {
    return <NavSeparator config={config} sidebarCollapsed={sidebarCollapsed} />
  }

  // Handle collapsible group - ensure items is defined and is an array
  if (
    'items' in config &&
    Array.isArray(config.items) &&
    config.items.length > 0
  ) {
    return (
      <NavGroupCollapsible
        config={config}
        basePath={basePath}
        sidebarCollapsed={sidebarCollapsed}
        onItemClick={onItemClick}
      />
    )
  }

  // Handle simple nav item
  return (
    <NavItem
      config={config as NavItemConfig}
      basePath={basePath}
      sidebarCollapsed={sidebarCollapsed}
      onItemClick={onItemClick}
    />
  )
}

/**
 * NavItem - Individual navigation item
 */
function NavItem({
  config,
  basePath,
  sidebarCollapsed,
  onItemClick,
  isNested = false,
}: {
  config: NavItemConfig
  basePath: string
  sidebarCollapsed: boolean
  onItemClick?: () => void
  isNested?: boolean
}) {
  const { t } = useTranslation('navigation')
  const location = useLocation()

  const itemPath = `${basePath}${config.path}`
  const isActive = isActiveRoute(config.path, location.pathname, basePath)
  const Icon = config.icon

  return (
    <Link
      to={itemPath}
      onClick={onItemClick}
      className={cn(
        'flex items-center gap-3 rounded-lg transition-colors text-sm',
        'hover:bg-base-200 active:scale-98',
        // Active state
        isActive && 'bg-primary/10 text-primary font-semibold',
        !isActive && 'text-base-content',
        // Sizing
        isNested ? 'px-4 py-2.5 text-sm h-10' : 'px-4 py-2 ',
        // Indentation for nested items
        isNested && 'ps-8',
        // Desktop collapsed: center icon
        sidebarCollapsed && !isNested && 'justify-center px-0',
      )}
      title={sidebarCollapsed && !isNested ? t(config.key) : undefined}
    >
      <Icon size={20} className="shrink-0" />
      {/* Hide text when collapsed on desktop (or always show for nested) */}
      {(!sidebarCollapsed || isNested) && (
        <span className="truncate flex-1">{t(config.key)}</span>
      )}
      {/* Badge */}
      {config.badge && (!sidebarCollapsed || isNested) && (
        <Badge variant={config.badge.variant ?? 'primary'} size="sm">
          {config.badge.label}
        </Badge>
      )}
    </Link>
  )
}

/**
 * NavGroupCollapsible - Collapsible navigation group
 */
function NavGroupCollapsible({
  config,
  basePath,
  sidebarCollapsed,
  onItemClick,
}: {
  config: NavGroupConfig
  basePath: string
  sidebarCollapsed: boolean
  onItemClick?: () => void
}) {
  const { t } = useTranslation('navigation')
  const { isRTL } = useLanguage()
  const location = useLocation()

  const navGroupsExpanded = useStore(businessStore, (s) => s.navGroupsExpanded)
  const isExpanded =
    navGroupsExpanded?.[config.key] ?? config.defaultExpanded ?? false

  // Check if any child is active
  const hasActiveChild =
    config.items?.some((item) =>
      isActiveRoute(item.path, location.pathname, basePath),
    ) ?? false

  const Icon = config.icon

  // If sidebar is collapsed on desktop, only show parent icon (no group functionality)
  if (sidebarCollapsed) {
    // If parent is clickable, show as simple item
    if (config.path && Icon) {
      const itemPath = `${basePath}${config.path}`
      const isActive = isActiveRoute(config.path, location.pathname, basePath)

      return (
        <Link
          to={itemPath}
          onClick={onItemClick}
          className={cn(
            'flex items-center justify-center px-0 py-3 rounded-lg transition-colors  text-sm',
            'hover:bg-base-200 active:scale-98',
            isActive && 'bg-primary/10 text-primary font-semibold',
            !isActive && hasActiveChild && 'text-primary font-medium',
            !isActive && !hasActiveChild && 'text-base-content',
          )}
          title={t(config.key)}
        >
          <Icon size={20} className="shrink-0" />
        </Link>
      )
    }

    // Container-only group: show first child's icon as representative
    if (config.items.length > 0) {
      const FirstIcon = config.items[0].icon
      return (
        <div
          className={cn(
            'flex items-center justify-center px-0 py-3 rounded-lg  text-sm',
            hasActiveChild && 'text-primary',
            !hasActiveChild && 'text-base-content',
          )}
          title={t(config.key)}
        >
          <FirstIcon size={20} className="shrink-0" />
        </div>
      )
    }

    return null
  }

  // Full expanded sidebar: show collapsible group
  const handleToggle = () => {
    toggleNavGroup(config.key)
  }

  const handleParentClick = () => {
    if (config.collapsible) {
      handleToggle()
    }
    if (onItemClick) {
      onItemClick()
    }
  }

  return (
    <div>
      {/* Parent Item */}
      {config.path && Icon ? (
        // Clickable parent
        <Link
          to={`${basePath}${config.path}`}
          onClick={onItemClick}
          className={cn(
            'flex items-center gap-3 px-4 py-2 rounded-lg transition-colors  text-sm',
            'hover:bg-base-200 active:scale-98',
            hasActiveChild && 'text-primary font-medium',
            !hasActiveChild && 'text-base-content',
          )}
        >
          <Icon size={20} className="shrink-0" />
          <span className="truncate flex-1">{t(config.key)}</span>
          {config.badge && (
            <Badge variant={config.badge.variant ?? 'primary'} size="sm">
              {config.badge.label}
            </Badge>
          )}
          {config.collapsible && (
            <button
              type="button"
              onClick={(e) => {
                e.preventDefault()
                e.stopPropagation()
                handleToggle()
              }}
              className="p-1 -me-1 hover:bg-base-300 rounded transition-colors"
              aria-expanded={isExpanded}
              aria-controls={`nav-group-${config.key}`}
            >
              <ChevronRight
                size={16}
                className={cn(
                  'transition-transform duration-200 ease-in-out',
                  isExpanded && 'rotate-90',
                  isRTL && !isExpanded && 'rotate-180',
                )}
              />
            </button>
          )}
        </Link>
      ) : (
        // Container-only parent (not clickable)
        <button
          type="button"
          onClick={handleParentClick}
          className={cn(
            'flex items-center gap-3 px-4 py-2 rounded-lg transition-colors w-full  text-sm',
            'hover:bg-base-200 active:scale-98 text-start',
            hasActiveChild && 'text-primary font-medium',
            !hasActiveChild && 'text-base-content',
          )}
          aria-expanded={isExpanded}
          aria-controls={`nav-group-${config.key}`}
        >
          {Icon && <Icon size={20} className="shrink-0" />}
          <span className="truncate flex-1">{t(config.key)}</span>
          {config.badge && (
            <Badge variant={config.badge.variant ?? 'primary'} size="sm">
              {config.badge.label}
            </Badge>
          )}
          {config.collapsible && (
            <ChevronRight
              size={16}
              className={cn(
                'transition-transform duration-200 ease-in-out shrink-0',
                isExpanded && 'rotate-90',
                isRTL && !isExpanded && 'rotate-180',
              )}
            />
          )}
        </button>
      )}

      {/* Child Items */}
      {config.collapsible && (
        <div
          id={`nav-group-${config.key}`}
          role="group"
          aria-labelledby={`nav-group-label-${config.key}`}
          className={cn(
            'overflow-hidden transition-all duration-200 ease-in-out',
            isExpanded ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0',
          )}
        >
          <div className="space-y-1 pt-1">
            {config.items.map((item) => (
              <NavItem
                key={item.key}
                config={item}
                basePath={basePath}
                sidebarCollapsed={false}
                onItemClick={onItemClick}
                isNested
              />
            ))}
          </div>
        </div>
      )}

      {/* Non-collapsible group: always show children */}
      {!config.collapsible && (
        <div className="space-y-1 pt-1">
          {config.items.map((item) => (
            <NavItem
              key={item.key}
              config={item}
              basePath={basePath}
              sidebarCollapsed={false}
              onItemClick={onItemClick}
              isNested
            />
          ))}
        </div>
      )}
    </div>
  )
}

/**
 * NavSeparator - Section separator with label
 */
function NavSeparator({
  config,
  sidebarCollapsed,
}: {
  config: { key: string; type: 'separator' }
  sidebarCollapsed: boolean
}) {
  const { t } = useTranslation('navigation')

  // Hide separator when sidebar is collapsed
  if (sidebarCollapsed) return null

  return (
    <div
      role="separator"
      aria-label={t(config.key)}
      className="pt-6 mt-4 border-t border-base-300"
    >
      <div className="px-4 py-2">
        <span className="text-xs font-semibold uppercase text-base-content/50">
          {t(config.key)}
        </span>
      </div>
    </div>
  )
}

/**
 * Check if a route is active
 */
function isActiveRoute(
  itemPath: string,
  currentPath: string,
  basePath: string,
): boolean {
  const fullPath = `${basePath}${itemPath}`

  // Special case: Dashboard (exact match)
  if (itemPath === '') {
    return currentPath === basePath || currentPath === `${basePath}/`
  }

  // Standard case: prefix match
  return currentPath.startsWith(fullPath)
}
