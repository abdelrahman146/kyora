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
import { useSearchParams, useParams, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Plus, Edit, Trash2, Eye } from "lucide-react";
import { DashboardLayout } from "../../components/templates";
import { useBusinessStore } from "../../stores/businessStore";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import { SearchInput } from "../../components/molecules/SearchInput";
import { CustomerCard } from "../../components/molecules/CustomerCard";
import { InfiniteScroll } from "../../components/molecules/InfiniteScroll";
import { Pagination } from "../../components/molecules/Pagination";
import { FilterButton } from "../../components/organisms/FilterButton";
import { AddCustomerSheet } from "../../components/organisms";
import { EditCustomerSheet } from "../../components/organisms/customers/EditCustomerSheet";
import { Dialog } from "../../components/atoms/Dialog";
import { Table } from "../../components/organisms/Table";
import type { TableColumn } from "../../components/organisms/Table";
import { Avatar } from "../../components/atoms/Avatar";
import { listCustomers, deleteCustomer } from "../../api/customer";
import type { Customer } from "../../api/types/customer";
import toast from "react-hot-toast";
import { translateErrorAsync } from "@/lib/translateError";

export default function CustomersPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { businessDescriptor: urlBusinessDescriptor } = useParams<{ businessDescriptor: string }>();
  const { selectedBusiness, setSelectedBusinessId, businesses } = useBusinessStore();
  const isMobile = useMediaQuery("(max-width: 768px)");
  const [searchParams, setSearchParams] = useSearchParams();

  // Sync URL business descriptor with state on mount
  useEffect(() => {
    if (urlBusinessDescriptor && businesses.length > 0) {
      const business = businesses.find((b) => b.descriptor === urlBusinessDescriptor);
      if (business && business.id !== selectedBusiness?.id) {
        setSelectedBusinessId(business.id);
      }
    }
  }, [urlBusinessDescriptor, businesses, selectedBusiness, setSelectedBusinessId]);

  // Use selected business descriptor from state for operations
  const businessDescriptor = selectedBusiness?.descriptor ?? urlBusinessDescriptor;

  // State
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [totalItems, setTotalItems] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [isAddCustomerOpen, setIsAddCustomerOpen] = useState(false);
  const [isEditCustomerOpen, setIsEditCustomerOpen] = useState(false);
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  // URL params
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("pageSize") ?? (isMobile ? "20" : "20"), 10);
  const search = searchParams.get("search") ?? "";
  const sortBy = searchParams.get("sortBy") ?? "";
  const sortOrder = (searchParams.get("sortOrder") ?? "desc") as "asc" | "desc";

  // Fetch customers
  const fetchCustomers = useCallback(async (append?: boolean) => {
    if (!businessDescriptor) return;

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
        businessDescriptor,
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

      setTotalItems(response.totalCount);
      setTotalPages(response.totalPages);
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    } finally {
      setIsLoading(false);
      setIsLoadingMore(false);
    }
  }, [businessDescriptor, page, pageSize, search, sortBy, sortOrder, searchParams, setSearchParams, t]);

  // Initial load and when params change
  useEffect(() => {
    void fetchCustomers(false);
  }, [businessDescriptor, page, pageSize, search, sortBy, sortOrder, fetchCustomers]);

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
    if (!businessDescriptor) return;
    void navigate(`/businesses/${businessDescriptor}/customers/${customer.id}`);
  };

  const handleEditClick = (customer: Customer, event?: React.MouseEvent) => {
    if (event) {
      event.stopPropagation();
    }
    setSelectedCustomer(customer);
    setIsEditCustomerOpen(true);
  };

  const handleDeleteClick = (customer: Customer, event?: React.MouseEvent) => {
    if (event) {
      event.stopPropagation();
    }
    setSelectedCustomer(customer);
    setIsDeleteDialogOpen(true);
  };

  const handleDelete = async () => {
    if (!businessDescriptor || !selectedCustomer) return;

    try {
      setIsDeleting(true);
      await deleteCustomer(businessDescriptor, selectedCustomer.id);
      toast.success(t("customers.delete_success"));
      await fetchCustomers(false);
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    } finally {
      setIsDeleting(false);
      setIsDeleteDialogOpen(false);
      setSelectedCustomer(null);
    }
  };

  const handleCustomerUpdated = () => {
    void fetchCustomers(false);
  };

  const businessCountryCode = selectedBusiness?.countryCode ?? "AE";

  // Table columns for desktop
  const tableColumns: TableColumn<Customer>[] = [
    {
      key: "name",
      label: t("customers.name"),
      sortable: true,
      render: (customer) => (
        <div className="flex items-center gap-3">
          <Avatar
            src={customer.avatarUrl}
            alt={customer.name}
            fallback={customer.name
              .split(" ")
              .map((w) => w[0])
              .join("")
              .toUpperCase()
              .slice(0, 2)}
            size="sm"
          />
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
      render: (customer) => (
        <div className="badge badge-success badge-sm">
          {customer.ordersCount ?? 0}
        </div>
      ),
    },
    {
      key: "totalSpent",
      label: t("customers.total_spent"),
      sortable: true,
      align: "end",
      render: (customer) => {
        const spent = customer.totalSpent ?? 0;
        const currency = selectedBusiness?.currency ?? "AED";
        return (
          <span className="font-semibold">
            {currency} {spent.toFixed(2)}
          </span>
        );
      },
    },
    {
      key: "actions",
      label: t("common.actions"),
      align: "center",
      width: "120px",
      render: (customer) => (
        <div className="flex items-center justify-center gap-2">
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square"
            onClick={(e) => {
              e.stopPropagation();
              handleCustomerClick(customer);
            }}
            aria-label={t("common.view")}
            title={t("common.view")}
          >
            <Eye size={16} />
          </button>
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square"
            onClick={(e) => {
              handleEditClick(customer, e);
            }}
            aria-label={t("common.edit")}
            title={t("common.edit")}
          >
            <Edit size={16} />
          </button>
          <button
            type="button"
            className="btn btn-ghost btn-sm btn-square text-error"
            onClick={(e) => {
              handleDeleteClick(customer, e);
            }}
            aria-label={t("common.delete")}
            title={t("common.delete")}
          >
            <Trash2 size={16} />
          </button>
        </div>
      ),
    },
  ];

  // Guard: No business descriptor in URL
  if (!businessDescriptor) {
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
          <button
            type="button"
            className="btn btn-primary gap-2"
            onClick={() => {
              setIsAddCustomerOpen(true);
            }}
            disabled={!businessDescriptor}
            aria-disabled={!businessDescriptor}
          >
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
          <FilterButton
            title={t("customers.filters")}
            buttonText={t("common.filter")}
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
          </FilterButton>
        </div>

        {/* Desktop: Table View */}
        {!isMobile && (
          <>
            <div className="overflow-x-auto">
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
            </div>
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
                  <div key={customer.id} className="relative group">
                    <CustomerCard
                      customer={customer}
                      onClick={handleCustomerClick}
                      ordersCount={customer.ordersCount ?? 0}
                      totalSpent={customer.totalSpent ?? 0}
                      currency={selectedBusiness?.currency ?? "AED"}
                    />
                    <div className="absolute top-2 end-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        type="button"
                        className="btn btn-sm btn-circle btn-ghost bg-base-100 shadow-md"
                        onClick={(e) => {
                          handleEditClick(customer, e);
                        }}
                        aria-label={t("common.edit")}
                      >
                        <Edit size={14} />
                      </button>
                      <button
                        type="button"
                        className="btn btn-sm btn-circle btn-ghost bg-base-100 shadow-md text-error"
                        onClick={(e) => {
                          handleDeleteClick(customer, e);
                        }}
                        aria-label={t("common.delete")}
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </InfiniteScroll>
        )}
      </div>

      <AddCustomerSheet
        isOpen={isAddCustomerOpen}
        onClose={() => {
          setIsAddCustomerOpen(false);
        }}
        businessDescriptor={businessDescriptor}
        businessCountryCode={businessCountryCode}
        onCreated={async () => {
          await fetchCustomers(false);
        }}
      />

      {selectedCustomer && (
        <EditCustomerSheet
          isOpen={isEditCustomerOpen}
          onClose={() => {
            setIsEditCustomerOpen(false);
            setSelectedCustomer(null);
          }}
          businessDescriptor={businessDescriptor}
          customer={selectedCustomer}
          onUpdated={handleCustomerUpdated}
        />
      )}

      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false);
          setSelectedCustomer(null);
        }}
        title={t("customers.delete_confirm_title")}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setIsDeleteDialogOpen(false);
                setSelectedCustomer(null);
              }}
              disabled={isDeleting}
            >
              {t("common.cancel")}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => {
                void handleDelete();
              }}
              disabled={isDeleting}
            >
              {isDeleting && <span className="loading loading-spinner loading-sm" />}
              {t("common.delete")}
            </button>
          </div>
        }
      >
        <p>{t("customers.delete_confirm_message", { name: selectedCustomer?.name })}</p>
      </Dialog>
    </DashboardLayout>
  );
}
