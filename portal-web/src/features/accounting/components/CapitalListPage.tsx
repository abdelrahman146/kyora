/**
 * CapitalListPage Component
 *
 * List view for capital transactions: Investments and Withdrawals.
 * Features separate tabs for each transaction type.
 *
 * Layout:
 * - Tab switcher (Investments | Withdrawals)
 * - Filtered list of transactions
 * - Pagination
 * - Empty states with helpful CTAs
 *
 * Mobile-first with card-based layout.
 */

import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus } from 'lucide-react'

import { TransactionCard } from './TransactionCard'
import { CreateTransactionSheet } from './sheets/CreateTransactionSheet'
import type { Investment, Withdrawal } from '@/api/accounting'
import { useInvestmentsQuery, useWithdrawalsQuery } from '@/api/accounting'
import { Button } from '@/components/atoms/Button'
import { Pagination } from '@/components/molecules/Pagination'

interface CapitalListPageProps {
  businessDescriptor: string
  currency: string
}

type TabType = 'investments' | 'withdrawals'

export function CapitalListPage({
  businessDescriptor,
  currency,
}: CapitalListPageProps) {
  const { t } = useTranslation('accounting')
  const [activeTab, setActiveTab] = useState<TabType>('investments')
  const [investmentsPage, setInvestmentsPage] = useState(1)
  const [withdrawalsPage, setWithdrawalsPage] = useState(1)
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [createTransactionType, setCreateTransactionType] = useState<
    'investment' | 'withdrawal'
  >('investment')

  const pageSize = 20

  // Fetch investments
  const {
    data: investmentsData,
    isLoading: isLoadingInvestments,
    refetch: refetchInvestments,
  } = useInvestmentsQuery(businessDescriptor, {
    page: investmentsPage,
    pageSize,
    orderBy: ['-investedAt'],
  })

  // Fetch withdrawals
  const {
    data: withdrawalsData,
    isLoading: isLoadingWithdrawals,
    refetch: refetchWithdrawals,
  } = useWithdrawalsQuery(businessDescriptor, {
    page: withdrawalsPage,
    pageSize,
    orderBy: ['-withdrawnAt'],
  })

  const investments = investmentsData?.items ?? []
  const withdrawals = withdrawalsData?.items ?? []

  const isInvestmentsTab = activeTab === 'investments'
  const isWithdrawalsTab = activeTab === 'withdrawals'
  const isLoading = isLoadingInvestments || isLoadingWithdrawals

  const handleActionComplete = async () => {
    if (isInvestmentsTab) {
      await refetchInvestments()
    } else {
      await refetchWithdrawals()
    }
  }

  const handleOpenCreate = (type: 'investment' | 'withdrawal') => {
    setCreateTransactionType(type)
    setIsCreateOpen(true)
  }

  const isEmpty =
    (isInvestmentsTab && investments.length === 0) ||
    (isWithdrawalsTab && withdrawals.length === 0)

  return (
    <div className="space-y-4">
      {/* Header with Action Button */}
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-2xl font-bold text-base-content">
          {t('header.capital')}
        </h1>

        {/* Add Transaction Button (Desktop) */}
        <div className="hidden sm:block">
          <Button
            variant="primary"
            size="md"
            onClick={() =>
              handleOpenCreate(isInvestmentsTab ? 'investment' : 'withdrawal')
            }
            className="gap-2"
          >
            <Plus className="h-5 w-5" />
            {t('actions.record_transaction')}
          </Button>
        </div>
      </div>

      {/* Tab Switcher */}
      <div className="tabs tabs-boxed bg-base-200 p-1">
        <button
          type="button"
          onClick={() => setActiveTab('investments')}
          className={`tab flex-1 ${isInvestmentsTab ? 'tab-active' : ''}`}
        >
          {t('tabs.investments')}
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('withdrawals')}
          className={`tab flex-1 ${isWithdrawalsTab ? 'tab-active' : ''}`}
        >
          {t('tabs.withdrawals')}
        </button>
      </div>

      {/* List Content */}
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className="skeleton h-24 w-full rounded-lg bg-base-200"
            />
          ))}
        </div>
      ) : isEmpty ? (
        <div className="card bg-base-100 p-8 text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-base-200">
            <Plus className="h-8 w-8 text-base-content/60" />
          </div>
          <h3 className="mb-2 text-lg font-semibold text-base-content">
            {t(
              isInvestmentsTab
                ? 'empty.no_investments'
                : 'empty.no_withdrawals',
            )}
          </h3>
          <p className="mb-6 text-sm text-base-content/70">
            {t(
              isInvestmentsTab
                ? 'empty.no_investments_message'
                : 'empty.no_withdrawals_message',
            )}
          </p>
          <Button
            variant="primary"
            size="md"
            onClick={() =>
              handleOpenCreate(isInvestmentsTab ? 'investment' : 'withdrawal')
            }
            className="gap-2"
          >
            <Plus className="h-5 w-5" />
            {t(
              isInvestmentsTab
                ? 'actions.add_investment'
                : 'actions.add_withdrawal',
            )}
          </Button>
        </div>
      ) : (
        <>
          {/* Investments List */}
          {isInvestmentsTab && (
            <div className="space-y-3">
              {investments.map((investment: Investment) => (
                <TransactionCard
                  key={investment.id}
                  transaction={investment}
                  type="investment"
                  currency={currency}
                  businessDescriptor={businessDescriptor}
                  onActionComplete={handleActionComplete}
                />
              ))}
            </div>
          )}

          {/* Withdrawals List */}
          {isWithdrawalsTab && (
            <div className="space-y-3">
              {withdrawals.map((withdrawal: Withdrawal) => (
                <TransactionCard
                  key={withdrawal.id}
                  transaction={withdrawal}
                  type="withdrawal"
                  currency={currency}
                  businessDescriptor={businessDescriptor}
                  onActionComplete={handleActionComplete}
                />
              ))}
            </div>
          )}

          {/* Pagination */}
          {isInvestmentsTab && investmentsData && (
            <Pagination
              currentPage={investmentsPage}
              totalPages={Math.ceil(
                (investmentsData.totalCount ?? 0) / pageSize,
              )}
              pageSize={pageSize}
              totalItems={investmentsData.totalCount ?? 0}
              itemsName={t('items.investments')}
              onPageChange={setInvestmentsPage}
            />
          )}
          {isWithdrawalsTab && withdrawalsData && (
            <Pagination
              currentPage={withdrawalsPage}
              totalPages={Math.ceil(
                (withdrawalsData.totalCount ?? 0) / pageSize,
              )}
              pageSize={pageSize}
              totalItems={withdrawalsData.totalCount ?? 0}
              itemsName={t('items.withdrawals')}
              onPageChange={setWithdrawalsPage}
            />
          )}
        </>
      )}

      {/* Mobile FAB */}
      <div className="sm:hidden fixed bottom-6 end-6 z-10">
        <Button
          variant="primary"
          size="lg"
          className="btn-circle border-2 border-base-300"
          onClick={() =>
            handleOpenCreate(isInvestmentsTab ? 'investment' : 'withdrawal')
          }
          aria-label={t('actions.record_transaction')}
        >
          <Plus className="h-6 w-6" />
        </Button>
      </div>

      {/* Create Transaction Sheet */}
      <CreateTransactionSheet
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        businessDescriptor={businessDescriptor}
        defaultType={createTransactionType}
        onCreated={handleActionComplete}
      />
    </div>
  )
}
