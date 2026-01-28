import type { LucideIcon } from 'lucide-react'

/**
 * Navigation Configuration Types
 *
 * Defines the sidebar navigation structure with support for:
 * - Top-level items
 * - Collapsible groups (max 2 levels)
 * - Section separators
 * - Badge indicators
 */

/**
 * Navigation item configuration
 */
export interface NavItemConfig {
  /** Unique identifier for the item (used for i18n key) */
  key: string
  /** Icon component (Lucide icon) */
  icon: LucideIcon
  /** Route path (relative to base path) */
  path: string
  /** Badge indicator (e.g., "NEW") */
  badge?: {
    label: string
    variant?:
      | 'primary'
      | 'secondary'
      | 'success'
      | 'error'
      | 'warning'
      | 'info'
      | 'neutral'
      | 'ghost'
  }
}

/**
 * Navigation group with optional collapsible behavior
 */
export interface NavGroupConfig {
  /** Unique identifier for the group (used for i18n key) */
  key: string
  /** Icon for parent item (if clickable) */
  icon?: LucideIcon
  /** Parent item path (if clickable) */
  path?: string
  /** Child items (max 1 level deep) */
  items: Array<NavItemConfig>
  /** Enable collapse/expand behavior */
  collapsible?: boolean
  /** Default expanded state (if collapsible) */
  defaultExpanded?: boolean
  /** Badge for parent item */
  badge?: {
    label: string
    variant?:
      | 'primary'
      | 'secondary'
      | 'success'
      | 'error'
      | 'warning'
      | 'info'
      | 'neutral'
      | 'ghost'
  }
}

/**
 * Section separator (e.g., "SETTINGS")
 */
export interface NavSectionConfig {
  /** Unique identifier (used for i18n key) */
  key: string
  /** Section type identifier */
  type: 'separator'
}

/**
 * Union type for all navigation config items
 */
export type NavConfigItem = NavItemConfig | NavGroupConfig | NavSectionConfig

/**
 * Full navigation configuration array
 */
export type NavConfig = Array<NavConfigItem>
