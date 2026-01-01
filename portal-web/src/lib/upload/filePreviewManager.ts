import { useEffect, useRef } from 'react'

/**
 * Manage object URLs for file previews with automatic cleanup
 */
export function useObjectURLs(files: Array<File>): Map<File, string> {
  const urlsRef = useRef<Map<File, string>>(new Map())

  useEffect(() => {
    const urls = urlsRef.current
    const currentFiles = new Set(files)

    // Create URLs for new files
    files.forEach((file) => {
      if (!urls.has(file)) {
        urls.set(file, URL.createObjectURL(file))
      }
    })

    // Revoke URLs for removed files
    Array.from(urls.keys()).forEach((file) => {
      if (!currentFiles.has(file)) {
        const url = urls.get(file)
        if (url) {
          URL.revokeObjectURL(url)
          urls.delete(file)
        }
      }
    })

    // Cleanup all URLs on unmount
    return () => {
      urls.forEach((url) => URL.revokeObjectURL(url))
      urls.clear()
    }
  }, [files])

  return urlsRef.current
}

/**
 * Create and manage a single object URL
 */
export function useObjectURL(file: File | null | undefined): string | null {
  const urlRef = useRef<string | null>(null)

  useEffect(() => {
    if (file) {
      urlRef.current = URL.createObjectURL(file)
    }

    return () => {
      if (urlRef.current) {
        URL.revokeObjectURL(urlRef.current)
        urlRef.current = null
      }
    }
  }, [file])

  return urlRef.current
}

/**
 * Revoke an object URL immediately
 */
export function revokeObjectURL(url: string | null | undefined): void {
  if (url && url.startsWith('blob:')) {
    URL.revokeObjectURL(url)
  }
}
