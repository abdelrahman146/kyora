import {
  BarChart3,
  Briefcase,
  Building2,
  Calculator,
  CreditCard,
  DollarSign,
  FileBarChart,
  LayoutDashboard,
  LineChart,
  Package,
  ShoppingCart,
  TrendingUp,
  Users,
  Wallet,
} from 'lucide-react'
import type { NavConfig } from '@/types/navigation'

/**
 * Navigation Configuration
 *
 * Defines the complete sidebar navigation structure for business-scoped routes.
 *
 * Structure:
 * - Top-level items (Dashboard, Inventory, Orders, Customers, Accounting)
 * - Financial Reports group (collapsible, clickable parent)
 * - Analytics group (collapsible, container-only, NEW badge)
 * - Settings section separator
 * - Settings items (Business Management, Workspace Settings, Billing)
 */
export const navigationConfig: NavConfig = [
  // Top-level items
  { key: 'dashboard', icon: LayoutDashboard, path: '' },
  { key: 'inventory', icon: Package, path: '/inventory' },
  { key: 'orders', icon: ShoppingCart, path: '/orders' },
  { key: 'customers', icon: Users, path: '/customers' },
  { key: 'accounting', icon: Calculator, path: '/accounting' },

  // Collapsible group: Financial Reports (clickable parent)
  {
    key: 'reports',
    icon: FileBarChart,
    path: '/reports',
    collapsible: true,
    defaultExpanded: false,
    items: [
      { key: 'reports.health', icon: TrendingUp, path: '/reports/health' },
      { key: 'reports.profit', icon: DollarSign, path: '/reports/profit' },
      { key: 'reports.cashflow', icon: Wallet, path: '/reports/cashflow' },
    ],
  },

  // Collapsible group: Analytics (container-only parent, NEW badge)
  {
    key: 'analytics',
    icon: BarChart3,
    collapsible: true,
    defaultExpanded: false,
    badge: { label: 'NEW', variant: 'primary' },
    items: [
      {
        key: 'analytics.customers',
        icon: Users,
        path: '/analytics/customers',
      },
      { key: 'analytics.sales', icon: LineChart, path: '/analytics/sales' },
      {
        key: 'analytics.inventory',
        icon: Package,
        path: '/analytics/inventory',
      },
    ],
  },

  // Section separator
  { key: 'settings_section', type: 'separator' },

  // Settings items (flat, non-collapsible)
  {
    key: 'settings.business',
    icon: Building2,
    path: '/settings/business',
  },
  {
    key: 'settings.workspace',
    icon: Briefcase,
    path: '/settings/workspace',
  },
  {
    key: 'settings.billing',
    icon: CreditCard,
    path: '/settings/billing',
  },
]
