import { useEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import type { ReactNode } from 'react'

interface TooltipProps {
  content: ReactNode
  children: ReactNode
  placement?: 'top' | 'bottom' | 'start' | 'end'
  className?: string
}

export function Tooltip({
  content,
  children,
  placement = 'top',
  className = '',
}: TooltipProps) {
  const [isVisible, setIsVisible] = useState(false)
  const [actualPlacement, setActualPlacement] = useState(placement)
  const [position, setPosition] = useState({ top: 0, left: 0 })

  const triggerRef = useRef<HTMLSpanElement>(null)
  const tooltipRef = useRef<HTMLDivElement>(null)
  const tooltipId = useRef(`tooltip-${Math.random().toString(36).substr(2, 9)}`)

  useEffect(() => {
    if (!isVisible || !triggerRef.current || !tooltipRef.current) return

    const trigger = triggerRef.current
    const tooltip = tooltipRef.current

    const triggerRect = trigger.getBoundingClientRect()
    const tooltipRect = tooltip.getBoundingClientRect()

    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    const gap = 8 // Gap between trigger and tooltip

    let newPlacement = placement
    let top = 0
    let left = 0

    // Calculate position based on placement and auto-flip if needed
    switch (placement) {
      case 'top':
        top = triggerRect.top - tooltipRect.height - gap
        left = triggerRect.left + triggerRect.width / 2 - tooltipRect.width / 2

        // Flip to bottom if not enough space on top
        if (top < 0) {
          newPlacement = 'bottom'
          top = triggerRect.bottom + gap
        }
        break

      case 'bottom':
        top = triggerRect.bottom + gap
        left = triggerRect.left + triggerRect.width / 2 - tooltipRect.width / 2

        // Flip to top if not enough space on bottom
        if (top + tooltipRect.height > viewportHeight) {
          newPlacement = 'top'
          top = triggerRect.top - tooltipRect.height - gap
        }
        break

      case 'start': {
        // In RTL, start is right; in LTR, start is left
        const isRTL = document.documentElement.dir === 'rtl'

        if (isRTL) {
          left = triggerRect.right + gap
          top =
            triggerRect.top + triggerRect.height / 2 - tooltipRect.height / 2

          // Flip to end if not enough space
          if (left + tooltipRect.width > viewportWidth) {
            newPlacement = 'end'
            left = triggerRect.left - tooltipRect.width - gap
          }
        } else {
          left = triggerRect.left - tooltipRect.width - gap
          top =
            triggerRect.top + triggerRect.height / 2 - tooltipRect.height / 2

          // Flip to end if not enough space
          if (left < 0) {
            newPlacement = 'end'
            left = triggerRect.right + gap
          }
        }
        break
      }

      case 'end': {
        // In RTL, end is left; in LTR, end is right
        const isRTLEnd = document.documentElement.dir === 'rtl'

        if (isRTLEnd) {
          left = triggerRect.left - tooltipRect.width - gap
          top =
            triggerRect.top + triggerRect.height / 2 - tooltipRect.height / 2

          // Flip to start if not enough space
          if (left < 0) {
            newPlacement = 'start'
            left = triggerRect.right + gap
          }
        } else {
          left = triggerRect.right + gap
          top =
            triggerRect.top + triggerRect.height / 2 - tooltipRect.height / 2

          // Flip to start if not enough space
          if (left + tooltipRect.width > viewportWidth) {
            newPlacement = 'start'
            left = triggerRect.left - tooltipRect.width - gap
          }
        }
        break
      }
    }

    // Keep tooltip within horizontal viewport bounds
    if (left < 0) {
      left = gap
    } else if (left + tooltipRect.width > viewportWidth) {
      left = viewportWidth - tooltipRect.width - gap
    }

    // Keep tooltip within vertical viewport bounds
    if (top < 0) {
      top = gap
    } else if (top + tooltipRect.height > viewportHeight) {
      top = viewportHeight - tooltipRect.height - gap
    }

    setActualPlacement(newPlacement)
    setPosition({ top, left })
  }, [isVisible, placement])

  const showTooltip = () => setIsVisible(true)
  const hideTooltip = () => setIsVisible(false)

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      hideTooltip()
    }
  }

  return (
    <>
      <span
        ref={triggerRef}
        onMouseEnter={showTooltip}
        onMouseLeave={hideTooltip}
        onFocus={showTooltip}
        onBlur={hideTooltip}
        onClick={showTooltip}
        onKeyDown={handleKeyDown}
        aria-describedby={isVisible ? tooltipId.current : undefined}
        className="inline-flex items-center"
      >
        {children}
      </span>

      {isVisible &&
        createPortal(
          <div
            ref={tooltipRef}
            id={tooltipId.current}
            role="tooltip"
            className={`fixed z-50 px-3 py-2 text-sm text-base-100 bg-base-content rounded-md transition-opacity duration-200 opacity-100 ${className}`}
            style={{
              top: `${position.top}px`,
              left: `${position.left}px`,
              pointerEvents: 'none',
            }}
          >
            {content}

            {/* Arrow */}
            <div
              className="absolute w-2 h-2 bg-base-content transform rotate-45"
              style={{
                ...(actualPlacement === 'top' && {
                  bottom: '-4px',
                  left: '50%',
                  marginLeft: '-4px',
                }),
                ...(actualPlacement === 'bottom' && {
                  top: '-4px',
                  left: '50%',
                  marginLeft: '-4px',
                }),
                ...(actualPlacement === 'start' && {
                  [document.documentElement.dir === 'rtl' ? 'left' : 'right']:
                    '-4px',
                  top: '50%',
                  marginTop: '-4px',
                }),
                ...(actualPlacement === 'end' && {
                  [document.documentElement.dir === 'rtl' ? 'right' : 'left']:
                    '-4px',
                  top: '50%',
                  marginTop: '-4px',
                }),
              }}
            />
          </div>,
          document.body,
        )}
    </>
  )
}
