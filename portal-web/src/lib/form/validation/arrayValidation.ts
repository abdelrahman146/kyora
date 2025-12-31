/**
 * Array Validation Utilities for TanStack Form
 *
 * Provides reusable validation functions for array fields, including:
 * - Min/max items validation
 * - Unique values validation
 * - Cross-item validation
 * - Conditional array requirements
 *
 * These utilities work with both string arrays and object arrays.
 *
 * @example
 * ```tsx
 * <form.Field
 *   name="phoneNumbers"
 *   validators={{
 *     onChange: (info) => {
 *       // Validate min/max items
 *       const minMaxError = validateArrayLength(info.value, { min: 1, max: 5 })
 *       if (minMaxError) return minMaxError
 *
 *       // Validate uniqueness
 *       return validateUniqueValues(info.value, { errorKey: 'form.duplicatePhones' })
 *     },
 *   }}
 * />
 * ```
 */

/**
 * Validates array has minimum and/or maximum number of items
 *
 * @param value - Array to validate
 * @param options - Min/max constraints and error keys
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * validators: {
 *   onChange: ({ value }) => validateArrayLength(value, {
 *     min: 1,
 *     max: 10,
 *     minErrorKey: 'form.atLeastOneItem',
 *     maxErrorKey: 'form.tooManyItems',
 *   })
 * }
 * ```
 */
export function validateArrayLength<T>(
  value: Array<T> | undefined | null,
  options: {
    min?: number
    max?: number
    minErrorKey?: string
    maxErrorKey?: string
  },
): string | undefined {
  const array = value ?? []

  // Validate minimum items
  if (options.min !== undefined && array.length < options.min) {
    return options.minErrorKey ?? 'form.minItemsRequired'
  }

  // Validate maximum items
  if (options.max !== undefined && array.length > options.max) {
    return options.maxErrorKey ?? 'form.maxItemsExceeded'
  }

  return undefined
}

/**
 * Validates all array items are unique
 *
 * @param value - Array to validate
 * @param options - Extractor function and error key
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * // Simple string array
 * validators: {
 *   onChange: ({ value }) => validateUniqueValues(value, {
 *     errorKey: 'form.duplicateEmails',
 *   })
 * }
 *
 * // Object array with extractor
 * validators: {
 *   onChange: ({ value }) => validateUniqueValues(value, {
 *     extractor: (item) => item.email,
 *     errorKey: 'form.duplicateEmails',
 *   })
 * }
 * ```
 */
export function validateUniqueValues<T>(
  value: Array<T> | undefined | null,
  options: {
    extractor?: (item: T) => any
    errorKey?: string
  } = {},
): string | undefined {
  const array = value ?? []

  if (array.length === 0) return undefined

  // Extract values using extractor function or use items directly
  const values = options.extractor ? array.map(options.extractor) : array

  // Check for duplicates using Set
  const uniqueValues = new Set(values)

  if (uniqueValues.size !== values.length) {
    return options.errorKey ?? 'form.duplicateValues'
  }

  return undefined
}

/**
 * Validates each array item using a custom validator function
 *
 * @param value - Array to validate
 * @param validator - Function that validates each item
 * @param options - Configuration options
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * validators: {
 *   onChange: ({ value }) => validateArrayItems(value, (item, index) => {
 *     if (!item.name) return 'form.nameRequired'
 *     if (item.price < 0) return 'form.priceNegative'
 *     return undefined
 *   }, {
 *     errorKey: 'form.invalidItem',
 *   })
 * }
 * ```
 */
export function validateArrayItems<T>(
  value: Array<T> | undefined | null,
  validator: (item: T, index: number) => string | undefined,
  options: {
    errorKey?: string
  } = {},
): string | undefined {
  const array = value ?? []

  if (array.length === 0) return undefined

  // Validate each item
  for (let i = 0; i < array.length; i++) {
    const error = validator(array[i], i)
    if (error) {
      return options.errorKey ?? error
    }
  }

  return undefined
}

/**
 * Validates array items don't overlap on a specific property (e.g., time ranges)
 *
 * @param value - Array to validate
 * @param options - Range extractor and error key
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * // Validate time ranges don't overlap
 * validators: {
 *   onChange: ({ value }) => validateNoOverlap(value, {
 *     extractor: (item) => ({ start: item.startTime, end: item.endTime }),
 *     errorKey: 'form.overlappingTimeRanges',
 *   })
 * }
 * ```
 */
export function validateNoOverlap<T>(
  value: Array<T> | undefined | null,
  options: {
    extractor: (item: T) => {
      start: number | Date | string
      end: number | Date | string
    }
    errorKey?: string
  },
): string | undefined {
  const array = value ?? []

  if (array.length < 2) return undefined

  // Extract ranges
  const ranges = array.map(options.extractor)

  // Convert to comparable values
  const comparableRanges = ranges.map((range) => ({
    start:
      range.start instanceof Date
        ? range.start.getTime()
        : typeof range.start === 'string'
          ? new Date(range.start).getTime()
          : range.start,
    end:
      range.end instanceof Date
        ? range.end.getTime()
        : typeof range.end === 'string'
          ? new Date(range.end).getTime()
          : range.end,
  }))

  // Check for overlaps (O(n²) but acceptable for reasonable array sizes)
  for (let i = 0; i < comparableRanges.length; i++) {
    for (let j = i + 1; j < comparableRanges.length; j++) {
      const a = comparableRanges[i]
      const b = comparableRanges[j]

      // Check if ranges overlap
      if (a.start < b.end && b.start < a.end) {
        return options.errorKey ?? 'form.overlappingRanges'
      }
    }
  }

  return undefined
}

/**
 * Validates array conditionally based on another field value
 *
 * @param value - Array to validate
 * @param condition - Condition function that returns true if validation should run
 * @param validator - Validation function to run when condition is true
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * <form.Field
 *   name="phoneNumbers"
 *   validators={{
 *     onChange: ({ value, fieldApi }) => {
 *       // Only require phone numbers if hasPhoneSupport is true
 *       return validateArrayConditionally(
 *         value,
 *         () => fieldApi.form.getFieldValue('hasPhoneSupport'),
 *         (array) => validateArrayLength(array, {
 *           min: 1,
 *           minErrorKey: 'form.atLeastOnePhone',
 *         })
 *       )
 *     },
 *   }}
 * />
 * ```
 */
export function validateArrayConditionally<T>(
  value: Array<T> | undefined | null,
  condition: () => boolean,
  validator: (array: Array<T>) => string | undefined,
): string | undefined {
  if (!condition()) {
    return undefined
  }

  return validator(value ?? [])
}

/**
 * Combines multiple array validators with OR logic
 * Returns error only if ALL validators fail
 *
 * @param value - Array to validate
 * @param validators - Array of validator functions
 * @returns Error key or undefined if any validator passes
 *
 * @example
 * ```tsx
 * validators: {
 *   onChange: ({ value }) => validateArrayOr(value, [
 *     (array) => validateArrayLength(array, { min: 1 }),
 *     (array) => array.some(item => item.isPrimary) ? undefined : 'form.noPrimary',
 *   ])
 * }
 * ```
 */
export function validateArrayOr<T>(
  value: Array<T> | undefined | null,
  validators: Array<(array: Array<T>) => string | undefined>,
): string | undefined {
  const array = value ?? []
  const errors: Array<string> = []

  for (const validator of validators) {
    const error = validator(array)
    if (!error) {
      // At least one validator passed
      return undefined
    }
    errors.push(error)
  }

  // All validators failed, return first error
  return errors[0]
}

/**
 * Combines multiple array validators with AND logic
 * Returns first error encountered
 *
 * @param value - Array to validate
 * @param validators - Array of validator functions
 * @returns Error key or undefined if all validators pass
 *
 * @example
 * ```tsx
 * validators: {
 *   onChange: ({ value }) => validateArrayAnd(value, [
 *     (array) => validateArrayLength(array, { min: 1, max: 10 }),
 *     (array) => validateUniqueValues(array),
 *     (array) => validateArrayItems(array, (item) =>
 *       item.email?.includes('@') ? undefined : 'form.invalidEmail'
 *     ),
 *   ])
 * }
 * ```
 */
export function validateArrayAnd<T>(
  value: Array<T> | undefined | null,
  validators: Array<(array: Array<T>) => string | undefined>,
): string | undefined {
  const array = value ?? []

  for (const validator of validators) {
    const error = validator(array)
    if (error) {
      // Return first error
      return error
    }
  }

  // All validators passed
  return undefined
}

/**
 * Validates array items match a specific count of a property value
 *
 * @param value - Array to validate
 * @param options - Extractor, expected count, and error key
 * @returns Error key or undefined if valid
 *
 * @example
 * ```tsx
 * // Ensure exactly one primary contact
 * validators: {
 *   onChange: ({ value }) => validateArrayCount(value, {
 *     extractor: (item) => item.isPrimary,
 *     matchValue: true,
 *     exactCount: 1,
 *     errorKey: 'form.exactlyOnePrimary',
 *   })
 * }
 *
 * // Ensure at least 2 admin users
 * validators: {
 *   onChange: ({ value }) => validateArrayCount(value, {
 *     extractor: (item) => item.role,
 *     matchValue: 'admin',
 *     minCount: 2,
 *     errorKey: 'form.atLeastTwoAdmins',
 *   })
 * }
 * ```
 */
export function validateArrayCount<T>(
  value: Array<T> | undefined | null,
  options: {
    extractor: (item: T) => any
    matchValue: any
    exactCount?: number
    minCount?: number
    maxCount?: number
    errorKey?: string
  },
): string | undefined {
  const array = value ?? []

  // Count items matching the value
  const count = array.filter(
    (item) => options.extractor(item) === options.matchValue,
  ).length

  // Validate exact count
  if (options.exactCount !== undefined && count !== options.exactCount) {
    return options.errorKey ?? 'form.invalidCount'
  }

  // Validate min count
  if (options.minCount !== undefined && count < options.minCount) {
    return options.errorKey ?? 'form.countTooLow'
  }

  // Validate max count
  if (options.maxCount !== undefined && count > options.maxCount) {
    return options.errorKey ?? 'form.countTooHigh'
  }

  return undefined
}
