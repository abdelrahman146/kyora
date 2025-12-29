/**
 * Table Component
 *
 * A flexible, responsive table component built on DaisyUI.
 * Supports sortable columns, loading states, and empty states.
 *
 * Features:
 * - Type-safe column definitions
 * - Sortable columns with visual indicators
 * - Loading skeleton states
 * - Empty state with custom message
 * - Sticky header on scroll
 * - RTL-compatible
 */

import { ChevronDown, ChevronUp } from "lucide-react";
import { Skeleton } from "../atoms/Skeleton";
import type { ReactNode } from "react";

export interface TableColumn<T> {
  key: string;
  label: string;
  sortable?: boolean;
  render: (item: T) => ReactNode;
  width?: string; // e.g., "w-1/4", "w-32"
  align?: "start" | "center" | "end";
}

export interface TableProps<T> {
  columns: Array<TableColumn<T>>;
  data: Array<T>;
  keyExtractor: (item: T) => string;
  isLoading?: boolean;
  emptyMessage?: string;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  onSort?: (key: string) => void;
  stickyHeader?: boolean;
}

export function Table<T>({
  columns,
  data,
  keyExtractor,
  isLoading = false,
  emptyMessage = "No data available",
  sortBy,
  sortOrder,
  onSort,
  stickyHeader = true,
}: TableProps<T>) {
  const handleSort = (columnKey: string, sortable?: boolean) => {
    if (sortable && onSort) {
      onSort(columnKey);
    }
  };

  const getAlignClass = (align?: "start" | "center" | "end") => {
    switch (align) {
      case "center":
        return "text-center";
      case "end":
        return "text-end";
      default:
        return "text-start";
    }
  };

  return (
    <div className="overflow-x-auto rounded-box border border-base-300">
      <table className="table table-sm md:table-md">
        {/* Head */}
        <thead className={stickyHeader ? "sticky top-0 z-10 bg-base-200" : "bg-base-200"}>
          <tr>
            {columns.map((column) => (
              <th
                key={column.key}
                className={`${column.width ?? ""} ${getAlignClass(column.align)} ${
                  column.sortable ? "cursor-pointer select-none hover:bg-base-300" : ""
                }`}
                onClick={() => {
                  handleSort(column.key, column.sortable);
                }}
              >
                <div className="flex items-center gap-2">
                  <span>{column.label}</span>
                  {column.sortable && (
                    <span className="inline-flex flex-col opacity-50">
                      {sortBy === column.key ? (
                        sortOrder === "asc" ? (
                          <ChevronUp size={16} className="text-primary" />
                        ) : (
                          <ChevronDown size={16} className="text-primary" />
                        )
                      ) : (
                        <>
                          <ChevronUp size={12} className="-mb-1" />
                          <ChevronDown size={12} />
                        </>
                      )}
                    </span>
                  )}
                </div>
              </th>
            ))}
          </tr>
        </thead>

        {/* Body */}
        <tbody>
          {isLoading ? (
            // Loading skeleton rows
            Array.from({ length: 5 }).map((_, index) => (
              <tr key={`skeleton-${String(index)}`}>
                {columns.map((column) => (
                  <td key={column.key} className={getAlignClass(column.align)}>
                    <Skeleton className="h-4 w-full" />
                  </td>
                ))}
              </tr>
            ))
          ) : data.length === 0 ? (
            // Empty state
            <tr>
              <td colSpan={columns.length} className="text-center py-8">
                <div className="text-base-content/60">{emptyMessage}</div>
              </td>
            </tr>
          ) : (
            // Data rows
            data.map((item) => (
              <tr key={keyExtractor(item)} className="hover:bg-base-200">
                {columns.map((column) => (
                  <td key={column.key} className={getAlignClass(column.align)}>
                    {column.render(item)}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
