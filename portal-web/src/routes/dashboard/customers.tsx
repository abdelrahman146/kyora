/**
 * Customers Page
 *
 * Displays a list of all customers with search, filter, and pagination.
 *
 * Features:
 * - Desktop: Table view with pagination
 * - Mobile: Card view with infinite scroll
 * - Search: Debounced search input
 * - Filter: Drawer with filter options
 * - Sorting: Sortable columns
 * - RTL-compatible
 */

import { useState, useEffect, useCallback } from "react";
import { useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Plus, Filter } from "lucide-react";
import { DashboardLayout } from "../../components/templates";
import { useBusinessStore } from "../../stores/businessStore";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import { SearchInput } from "../../components/molecules/SearchInput";
import { CustomerCard } from "../../components/molecules/CustomerCard";
import { InfiniteScroll } from "../../components/molecules/InfiniteScroll";
import { Pagination } from "../../components/molecules/Pagination";
import { FilterDrawer } from "../../components/organisms/FilterDrawer";
import { Table } from "../../components/organisms/Table";
import type { TableColumn } from "../../components/organisms/Table";
import { listCustomers } from "../../api/customer";
import type { Customer } from "../../api/types/customer";

export default function CustomersPage() {
  const { t } = useTranslation();
  const { selectedBusiness } = useBusinessStore();
  const isMobile = useMediaQuery("(max-width: 768px)");
  const [searchParams, setSearchParams] = useSearchParams();

  // State
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [totalItems, setTotalItems] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [isFilterOpen, setIsFilterOpen] = useState(false);

  // URL params
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("pageSize") ?? (isMobile ? "20" : "20"), 10);
  const search = searchParams.get("search") ?? "";
  const sortBy = searchParams.get("sortBy") ?? "";
  const sortOrder = (searchParams.get("sortOrder") ?? "desc") as "asc" | "desc";

  // Fetch customers
  const fetchCustomers = useCallback(async (append?: boolean) => {
    if (!selectedBusiness) return;

    try {
      if (append) {
        setIsLoadingMore(true);
      } else {
        setIsLoading(true);
      }

      const orderBy = sortBy
        ? [`${sortOrder === "desc" ? "-" : ""}${sortBy}`]
        : ["-createdAt"];

      const response = await listCustomers({
        businessDescriptor: selectedBusiness.id,
        page: append ? page + 1 : page,
        pageSize,
        search: search || undefined,
        orderBy,
      });

      if (append) {
        setCustomers((prev) => [...prev, ...response.items]);
        setSearchParams({
          ...Object.fromEntries(searchParams),
          page: (page + 1).toString(),
        });
      } else {
        setCustomers(response.items);
      }

      setTotalItems(response.total);
      setTotalPages(response.totalPages);
    } catch (error) {
      console.error("Failed to fetch customers:", error);
    } finally {
      setIsLoading(false);
      setIsLoadingMore(false);
    }
  }, [selectedBusiness, page, pageSize, search, sortBy, sortOrder, searchParams, setSearchParams]);

  // Initial load and when params change
  useEffect(() => {
    void fetchCustomers(false);
  }, [selectedBusiness, page, pageSize, search, sortBy, sortOrder, fetchCustomers]);

  // Handlers
  const handleSearch = (value: string) => {
    setSearchParams({
      ...Object.fromEntries(searchParams),
      search: value,
      page: "1",
    });
  };

  const handlePageChange = (newPage: number) => {
    setSearchParams({
      ...Object.fromEntries(searchParams),
      page: newPage.toString(),
    });
  };

  const handlePageSizeChange = (newPageSize: number) => {
    setSearchParams({
      ...Object.fromEntries(searchParams),
      pageSize: newPageSize.toString(),
      page: "1",
    });
  };

  const handleSort = (key: string) => {
    const newSortOrder =
      sortBy === key && sortOrder === "asc" ? "desc" : "asc";

    setSearchParams({
      ...Object.fromEntries(searchParams),
      sortBy: key,
      sortOrder: newSortOrder,
      page: "1",
    });
  };

  const handleLoadMore = () => {
    void fetchCustomers(true);
  };

  const handleCustomerClick = (customer: Customer) => {
    console.log("Customer clicked:", customer.id);
  };

  // Table columns for desktop
  const tableColumns: TableColumn<Customer>[] = [
    {
      key: "name",
      label: t("customers.name"),
      sortable: true,
      render: (customer) => (
        <div className="flex items-center gap-2">
          <div className="avatar placeholder">
            <div className="bg-primary text-primary-content rounded-full w-10 h-10">
              <span className="text-xs">
                {customer.name
                  .split(" ")
                  .map((w) => w[0])
                  .join("")
                  .toUpperCase()
                  .slice(0, 2)}
              </span>
            </div>
          </div>
          <span className="font-medium">{customer.name}</span>
        </div>
      ),
    },
    {
      key: "phone",
      label: t("customers.phone"),
      render: (customer) => {
        if (customer.phoneCode && customer.phoneNumber) {
          return `${customer.phoneCode} ${customer.phoneNumber}`;
        }
        return <span className="text-base-content/40">—</span>;
      },
    },
    {
      key: "ordersCount",
      label: t("customers.orders_count"),
      sortable: true,
      align: "center",
      render: () => (
        <div className="badge badge-success badge-sm">0</div>
      ),
    },
    {
      key: "totalSpent",
      label: t("customers.total_spent"),
      sortable: true,
      align: "end",
      render: () => (
        <span className="font-semibold">AED 0.00</span>
      ),
    },
  ];

  // Guard: No business selected
  if (!selectedBusiness) {
    return (
      <DashboardLayout title={t("customers.title")}>
        <div className="alert alert-warning">
          <span>{t("dashboard.select_business_to_start")}</span>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout title={t("customers.title")}>
      <div className="space-y-4">
        {/* Header */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold">{t("customers.title")}</h1>
            <p className="text-sm text-base-content/60 mt-1">
              {t("customers.subtitle")}
            </p>
          </div>
          <button className="btn btn-primary gap-2">
            <Plus size={20} />
            {t("customers.add_customer")}
          </button>
        </div>

        {/* Toolbar */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="flex-1">
            <SearchInput
              value={search}
              onChange={handleSearch}
              placeholder={t("customers.search_placeholder")}
            />
          </div>
          <button
            className="btn btn-outline gap-2"
            onClick={() => {
              setIsFilterOpen(true);
            }}
          >
            <Filter size={18} />
            {t("common.filter")}
          </button>
        </div>

        {/* Desktop: Table View */}
        {!isMobile && (
          <>
            <Table
              columns={tableColumns}
              data={customers}
              keyExtractor={(customer) => customer.id}
              isLoading={isLoading}
              emptyMessage={t("customers.no_customers")}
              sortBy={sortBy}
              sortOrder={sortOrder}
              onSort={handleSort}
            />
            <Pagination
              currentPage={page}
              totalPages={totalPages}
              pageSize={pageSize}
              totalItems={totalItems}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
              itemsName={t("customers.customers").toLowerCase()}
            />
          </>
        )}

        {/* Mobile: Card View with Infinite Scroll */}
        {isMobile && (
          <InfiniteScroll
            hasMore={page < totalPages}
            isLoading={isLoadingMore}
            onLoadMore={handleLoadMore}
            loadingMessage={t("common.loading_more")}
            endMessage={t("customers.no_more_customers")}
          >
            <div className="space-y-3">
              {isLoading && customers.length === 0 ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <div key={i} className="skeleton h-32 rounded-box" />
                ))
              ) : customers.length === 0 ? (
                <div className="text-center py-12 text-base-content/60">
                  {t("customers.no_customers")}
                </div>
              ) : (
                customers.map((customer) => (
                  <CustomerCard
                    key={customer.id}
                    customer={customer}
                    onClick={handleCustomerClick}
                    ordersCount={0}
                    totalSpent={0}
                    currency="AED"
                  />
                ))
              )}
            </div>
          </InfiniteScroll>
        )}
      </div>

      {/* Filter Drawer */}
      <FilterDrawer
        isOpen={isFilterOpen}
        onClose={() => {
          setIsFilterOpen(false);
        }}
        title={t("customers.filters")}
        applyLabel={t("common.apply")}
        resetLabel={t("common.reset")}
        onApply={() => {
          // Filter application logic will be implemented here
        }}
        onReset={() => {
          // Filter reset logic will be implemented here
        }}
      >
        <div className="space-y-4">
          <p className="text-sm text-base-content/60">
            {t("customers.filters_coming_soon")}
          </p>
        </div>
      </FilterDrawer>
    </DashboardLayout>
  );
}
