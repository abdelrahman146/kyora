import {
  Banknote,
  Box,
  Boxes,
  Calculator,
  CreditCard,
  Plane,
  Receipt,
  Truck,
  Wallet,
} from 'lucide-react'
import type { ExpenseCategory } from '@/api/accounting'
import type { RecurringExpenseFrequency } from '@/api/types/accounting'

export const CATEGORY_OPTIONS: Array<{
  value: ExpenseCategory
  labelKey: string
}> = [
  { value: 'rent', labelKey: 'category.rent' },
  { value: 'marketing', labelKey: 'category.marketing' },
  { value: 'salaries', labelKey: 'category.salaries' },
  { value: 'packaging', labelKey: 'category.packaging' },
  { value: 'software', labelKey: 'category.software' },
  { value: 'logistics', labelKey: 'category.logistics' },
  { value: 'transaction_fee', labelKey: 'category.transaction_fee' },
  { value: 'travel', labelKey: 'category.travel' },
  { value: 'supplies', labelKey: 'category.supplies' },
  { value: 'other', labelKey: 'category.other' },
]

// Category-specific background colors
export const categoryColors: Record<ExpenseCategory, string> = {
  rent: 'bg-amber-100 text-amber-700',
  marketing: 'bg-purple-100 text-purple-700',
  salaries: 'bg-green-100 text-green-700',
  packaging: 'bg-blue-100 text-blue-700',
  software: 'bg-indigo-100 text-indigo-700',
  logistics: 'bg-orange-100 text-orange-700',
  transaction_fee: 'bg-slate-100 text-slate-700',
  travel: 'bg-cyan-100 text-cyan-700',
  supplies: 'bg-teal-100 text-teal-700',
  other: 'bg-base-200 text-base-content/70',
}

// Category icon mapping for expenses
export const categoryIcons: Record<ExpenseCategory, typeof Receipt> = {
  rent: Banknote,
  marketing: CreditCard,
  salaries: Wallet,
  packaging: Box,
  software: Calculator,
  logistics: Truck,
  transaction_fee: Receipt,
  travel: Plane,
  supplies: Boxes,
  other: Receipt,
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
