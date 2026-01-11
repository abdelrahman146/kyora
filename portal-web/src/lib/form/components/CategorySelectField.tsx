import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'

import { useFieldContext } from '../contexts'
import type { SelectFieldProps } from '../types'
import type { Category } from '@/api/inventory'
import { useCategoriesQuery, useCreateCategoryMutation } from '@/api/inventory'
import { FormSelect } from '@/components/form'
import { queryKeys } from '@/lib/queryKeys'
import { translateErrorAsync } from '@/lib/translateError'

interface CategorySelectFieldProps extends Omit<
  SelectFieldProps<string>,
  'options' | 'multiSelect'
> {
  businessDescriptor: string
}

export function CategorySelectField({
  businessDescriptor,
  label,
  placeholder,
  required,
  disabled,
  hint,
  ...props
}: CategorySelectFieldProps) {
  const field = useFieldContext<string>()
  const { t } = useTranslation(['inventory', 'errors'])
  const queryClient = useQueryClient()
  const [searchValue, setSearchValue] = useState('')

  const { data: categories = [], isLoading } =
    useCategoriesQuery(businessDescriptor)

  const createCategoryMutation = useCreateCategoryMutation(businessDescriptor, {
    onSuccess: (category) => {
      queryClient.setQueryData(
        [...queryKeys.inventory.all, 'categories', businessDescriptor],
        (prev?: Array<Category>) => {
          const existing = prev ?? []
          const exists = existing.some((item) => item.id === category.id)
          if (exists) return existing
          return [...existing, category]
        },
      )
      field.handleChange(category.id)
      setSearchValue('')
      toast.success(
        t('category_created', { ns: 'inventory', name: category.name }),
      )
    },
    onError: async (error) => {
      const message = await translateErrorAsync(error, t)
      toast.error(message)
    },
  })

  const options = useMemo(
    () =>
      categories.map((cat) => ({
        value: cat.id,
        label: cat.name,
      })),
    [categories],
  )

  const error = useMemo(() => {
    const errors = field.state.meta.errors
    if (errors.length === 0) return undefined

    const firstError = errors[0]
    if (typeof firstError === 'string') return t(firstError, { ns: 'errors' })
    if (
      typeof firstError === 'object' &&
      firstError &&
      'message' in firstError
    ) {
      const errorObj = firstError as { message: string; code?: number }
      return t(errorObj.message, { ns: 'errors' })
    }
    return undefined
  }, [field.state.meta.errors, t])

  const showError = field.state.meta.isTouched && error

  const handleCreateCategory = async (name: string) => {
    const trimmed = name.trim()
    if (!trimmed || createCategoryMutation.isPending) return
    await createCategoryMutation.mutateAsync({
      name: trimmed,
      descriptor: businessDescriptor,
    })
  }

  return (
    <FormSelect
      id={field.name}
      label={label}
      placeholder={placeholder}
      required={required}
      helperText={hint}
      options={options}
      value={field.state.value}
      onChange={(value) => field.handleChange(value as string)}
      error={showError}
      disabled={
        disabled ||
        isLoading ||
        createCategoryMutation.isPending ||
        field.state.meta.isValidating
      }
      onClose={field.handleBlur}
      aria-invalid={!field.state.meta.isValid && field.state.meta.isTouched}
      aria-describedby={showError ? `${field.name}-error` : undefined}
      clearable
      searchable
      searchValue={searchValue}
      onSearchChange={setSearchValue}
      onCreateOption={handleCreateCategory}
      isCreatingOption={createCategoryMutation.isPending}
      createOptionLabel={(name) =>
        t('create_category_option', { ns: 'inventory', name })
      }
      mobileTitle={label}
      {...props}
    />
  )
}
