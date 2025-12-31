export function getErrorText(error: unknown): string | undefined {
  if (!error) return undefined

  if (typeof error === 'string') return error
  if (typeof error === 'number' || typeof error === 'boolean')
    return String(error)

  if (Array.isArray(error)) {
    return getErrorText(error[0])
  }

  if (typeof error === 'object') {
    const record = error as Record<string, unknown>

    const message = record.message
    if (typeof message === 'string' && message.trim()) {
      return message
    }

    const key = record.key
    if (typeof key === 'string' && key.trim()) {
      return key
    }

    const fallback = record.fallback
    if (typeof fallback === 'string' && fallback.trim()) {
      return fallback
    }

    const errorField = record.error
    if (typeof errorField === 'string' && errorField.trim()) {
      return errorField
    }

    try {
      const json = JSON.stringify(error)
      if (json && json !== '{}' && json !== '[]') {
        return json
      }
    } catch {
      // ignore
    }

    return undefined
  }

  return undefined
}
