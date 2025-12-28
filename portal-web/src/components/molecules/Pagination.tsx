/**
 * Pagination Component
 *
 * A pagination component for desktop table views.
 *
 * Features:
 * - Page number navigation
 * - Previous/Next buttons
 * - Jump to first/last page
 * - Page size selector
 * - Current page indicator
 * - RTL-compatible (chevrons flip direction automatically)
 * - Responsive design
 */

import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-react";
import { useTranslation } from "react-i18next";

export interface PaginationProps {
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
  onPageChange: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
  pageSizeOptions?: number[];
  showPageSizeSelector?: boolean;
  itemsName?: string; // e.g., "customers", "orders"
}

export function Pagination({
  currentPage,
  totalPages,
  pageSize,
  totalItems,
  onPageChange,
  onPageSizeChange,
  pageSizeOptions = [10, 20, 50, 100],
  showPageSizeSelector = true,
  itemsName = "items",
}: PaginationProps) {
  const { t, i18n } = useTranslation();
  const isRTL = i18n.dir() === "rtl";
  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

  const getPageNumbers = (): (number | string)[] => {
    const pages: (number | string)[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (currentPage > 3) {
        pages.push("...");
      }

      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push("...");
      }

      pages.push(totalPages);
    }

    return pages;
  };

  if (totalPages === 0) {
    return null;
  }

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 px-4 py-3 bg-base-100 border-t border-base-300">
      {/* Info & Page Size Selector */}
      <div className="flex items-center gap-4 text-sm text-base-content/70">
        <span>
          {t('pagination.showing', {
            start: startItem,
            end: endItem,
            total: totalItems,
            items: itemsName,
          })}
        </span>

        {showPageSizeSelector && onPageSizeChange && (
          <div className="flex items-center gap-2">
            <span>{t('pagination.show')}</span>
            <select
              value={pageSize}
              onChange={(e) => {
                onPageSizeChange(Number(e.target.value));
              }}
              className="select select-sm select-bordered"
              aria-label={t('pagination.itemsPerPage')}
            >
              {pageSizeOptions.map((size) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {/* Pagination Controls */}
      <div className="join">
        {/* First Page */}
        <button
          type="button"
          onClick={() => {
            onPageChange(1);
          }}
          disabled={currentPage === 1}
          className="btn btn-sm join-item"
          aria-label={t('pagination.firstPage')}
        >
          {isRTL ? <ChevronsRight size={16} /> : <ChevronsLeft size={16} />}
        </button>

        {/* Previous Page */}
        <button
          type="button"
          onClick={() => {
            onPageChange(currentPage - 1);
          }}
          disabled={currentPage === 1}
          className="btn btn-sm join-item"
          aria-label={t('pagination.previousPage')}
        >
          {isRTL ? <ChevronRight size={16} /> : <ChevronLeft size={16} />}
        </button>

        {/* Page Numbers */}
        {getPageNumbers().map((page, index) => {
          if (page === "...") {
            return (
              <button
                key={`ellipsis-${String(index)}`}
                type="button"
                className="btn btn-sm join-item btn-disabled"
                disabled
              >
                ...
              </button>
            );
          }

          return (
            <button
              key={page}
              type="button"
              onClick={() => {
                onPageChange(page as number);
              }}
              className={`btn btn-sm join-item ${
                currentPage === page ? "btn-active" : ""
              }`}
              aria-label={t('pagination.page', { number: page })}
              aria-current={currentPage === page ? "page" : undefined}
            >
              {page}
            </button>
          );
        })}

        {/* Next Page */}
        <button
          type="button"
          onClick={() => {
            onPageChange(currentPage + 1);
          }}
          disabled={currentPage === totalPages}
          className="btn btn-sm join-item"
          aria-label={t('pagination.nextPage')}
        >
          {isRTL ? <ChevronLeft size={16} /> : <ChevronRight size={16} />}
        </button>

        {/* Last Page */}
        <button
          type="button"
          onClick={() => {
            onPageChange(totalPages);
          }}
          disabled={currentPage === totalPages}
          className="btn btn-sm join-item"
          aria-label={t('pagination.lastPage')}
        >
          {isRTL ? <ChevronsLeft size={16} /> : <ChevronsRight size={16} />}
        </button>
      </div>
    </div>
  );
}
