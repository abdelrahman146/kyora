/**
 * Reports Feature Module
 *
 * Exports all public components and utilities for the Financial Reports feature.
 *
 * This module provides financial reporting capabilities:
 * - Reports Hub (overview page with Safe to Draw + report cards)
 * - Business Health (Financial Position / Balance Sheet)
 * - Profit & Earnings (Profit and Loss Statement)
 * - Cash Movement (Cash Flow Statement)
 */

// Components
export { ReportCard, ReportCardSkeleton } from './components/ReportCard'
export type {
  ReportCardProps,
  ReportCardSkeletonProps,
} from './components/ReportCard'
