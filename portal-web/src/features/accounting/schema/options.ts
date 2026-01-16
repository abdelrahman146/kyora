import {
  Banknote,
  Box,
  Boxes,
  Briefcase,
  Building,
  Calculator,
  CreditCard,
  FileText,
  GraduationCap,
  HardDrive,
  Laptop,
  Lightbulb,
  Plane,
  Receipt,
  Scale,
  Search,
  Shield,
  Truck,
  Users,
  Wrench,
} from 'lucide-react'
import type { ExpenseCategory } from '@/api/accounting'
import type { RecurringExpenseFrequency } from '@/api/types/accounting'

/**
 * Category options for expense forms
 * Must match backend ExpenseCategory enum in:
 * backend/internal/domain/accounting/model.go
 */
export const CATEGORY_OPTIONS: Array<{
  value: ExpenseCategory
  labelKey: string
}> = [
  { value: 'office', labelKey: 'category.office' },
  { value: 'travel', labelKey: 'category.travel' },
  { value: 'supplies', labelKey: 'category.supplies' },
  { value: 'utilities', labelKey: 'category.utilities' },
  { value: 'payroll', labelKey: 'category.payroll' },
  { value: 'marketing', labelKey: 'category.marketing' },
  { value: 'rent', labelKey: 'category.rent' },
  { value: 'software', labelKey: 'category.software' },
  { value: 'maintenance', labelKey: 'category.maintenance' },
  { value: 'insurance', labelKey: 'category.insurance' },
  { value: 'taxes', labelKey: 'category.taxes' },
  { value: 'training', labelKey: 'category.training' },
  { value: 'consulting', labelKey: 'category.consulting' },
  { value: 'miscellaneous', labelKey: 'category.miscellaneous' },
  { value: 'legal', labelKey: 'category.legal' },
  { value: 'research', labelKey: 'category.research' },
  { value: 'equipment', labelKey: 'category.equipment' },
  { value: 'shipping', labelKey: 'category.shipping' },
  { value: 'transaction_fee', labelKey: 'category.transaction_fee' },
  { value: 'other', labelKey: 'category.other' },
]

// Category-specific background colors
export const categoryColors: Record<ExpenseCategory, string> = {
  office: 'bg-blue-100 text-blue-700',
  travel: 'bg-cyan-100 text-cyan-700',
  supplies: 'bg-teal-100 text-teal-700',
  utilities: 'bg-yellow-100 text-yellow-700',
  payroll: 'bg-green-100 text-green-700',
  marketing: 'bg-purple-100 text-purple-700',
  rent: 'bg-amber-100 text-amber-700',
  software: 'bg-indigo-100 text-indigo-700',
  maintenance: 'bg-orange-100 text-orange-700',
  insurance: 'bg-rose-100 text-rose-700',
  taxes: 'bg-red-100 text-red-700',
  training: 'bg-violet-100 text-violet-700',
  consulting: 'bg-fuchsia-100 text-fuchsia-700',
  miscellaneous: 'bg-gray-100 text-gray-700',
  legal: 'bg-slate-100 text-slate-700',
  research: 'bg-sky-100 text-sky-700',
  equipment: 'bg-emerald-100 text-emerald-700',
  shipping: 'bg-lime-100 text-lime-700',
  transaction_fee: 'bg-stone-100 text-stone-700',
  other: 'bg-base-200 text-base-content/70',
}

// Category icon mapping for expenses
export const categoryIcons: Record<ExpenseCategory, typeof Receipt> = {
  office: Building,
  travel: Plane,
  supplies: Boxes,
  utilities: Lightbulb,
  payroll: Users,
  marketing: CreditCard,
  rent: Banknote,
  software: Laptop,
  maintenance: Wrench,
  insurance: Shield,
  taxes: Calculator,
  training: GraduationCap,
  consulting: Briefcase,
  miscellaneous: Box,
  legal: Scale,
  research: Search,
  equipment: HardDrive,
  shipping: Truck,
  transaction_fee: Receipt,
  other: FileText,
}

export const FREQUENCY_OPTIONS: Array<{
  value: RecurringExpenseFrequency
  labelKey: string
}> = [
  { value: 'daily', labelKey: 'frequency.daily' },
  { value: 'weekly', labelKey: 'frequency.weekly' },
  { value: 'monthly', labelKey: 'frequency.monthly' },
  { value: 'yearly', labelKey: 'frequency.yearly' },
]
