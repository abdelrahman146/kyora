/**
 * useSelectSearch Hook
 *
 * Manages search state and filtering for select components.
 * Provides search query management and filtered options.
 */

import { useMemo, useState } from 'react'
import type { FormSelectOption } from './FormSelect'

interface UseSelectSearchProps<T> {
  options: Array<FormSelectOption<T>>
  searchable: boolean
}

interface UseSelectSearchReturn<T> {
  searchQuery: string
  setSearchQuery: (query: string) => void
  filteredOptions: Array<FormSelectOption<T>>
  clearSearch: () => void
}

export function useSelectSearch<T extends string>({
  options,
  searchable,
}: UseSelectSearchProps<T>): UseSelectSearchReturn<T> {
  const [searchQuery, setSearchQuery] = useState('')

  const filteredOptions = useMemo(() => {
    if (!searchable || !searchQuery.trim()) {
      return options
    }

    const query = searchQuery.toLowerCase()
    return options.filter((option) => {
      // Search in label
      if (option.label.toLowerCase().includes(query)) {
        return true
      }

      // Search in description if available
      if (option.description?.toLowerCase().includes(query)) {
        return true
      }

      return false
    })
  }, [options, searchable, searchQuery])

  const clearSearch = () => {
    setSearchQuery('')
  }

  return {
    searchQuery,
    setSearchQuery,
    filteredOptions,
    clearSearch,
  }
}
