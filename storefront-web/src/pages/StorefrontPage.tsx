import { useEffect, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { storefrontApi } from '../api/storefront';
import { applyBrandTheme } from '../theme/applyBrandTheme';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import { ProductList } from '../components/ProductList';
import { CartDrawer } from '../components/CartDrawer';
import { cartTotalQuantity, useCartStore } from '../cart/useCartStore';
import { formatMoney, money } from '../utils/money';

export function StorefrontPage() {
  const { storefrontPublicId } = useParams();
  const { t, i18n } = useTranslation();

  const openCart = useCartStore((s) => s.openCart);
  const itemsByVariantId = useCartStore((s) => s.itemsByVariantId);
  const totalQty = useMemo(() => cartTotalQuantity({ itemsByVariantId }), [itemsByVariantId]);

  const catalogQuery = useQuery({
    queryKey: ['catalog', storefrontPublicId],
    queryFn: () => storefrontApi.getCatalog(storefrontPublicId || ''),
    enabled: !!storefrontPublicId,
  });

  const catalog = catalogQuery.data;

  useEffect(() => {
    applyBrandTheme(catalog?.business?.storefrontTheme);
  }, [catalog?.business?.storefrontTheme]);

  const cartSubtotal = useMemo(() => {
    const items = Object.values(itemsByVariantId);
    return items.reduce((acc, it) => acc.add(money(it.unitPrice).mul(it.quantity)), money('0'));
  }, [itemsByVariantId]);

  if (!storefrontPublicId) {
    return <div className="p-6">Missing storefront id</div>;
  }

  const businessName = catalog?.business?.name || '…';
  const logoUrl = catalog?.business?.logoUrl || '';
  const avatarLetter = (businessName.trim().charAt(0) || '?').toUpperCase();

  return (
    <div className="min-h-dvh flex flex-col bg-base-100 text-base-content">
      <header className="sticky top-0 z-20 bg-base-100 border-b border-base-200/60 safe-top">
        <div className="mx-auto max-w-5xl px-4">
          <div className="navbar min-h-14 py-3">
            <div className="navbar-start gap-3 min-w-0">
              <div className="avatar">
                <div className="w-11 rounded-full ring ring-base-200/60 ring-offset-base-100 ring-offset-2">
                  {logoUrl ? (
                    <img src={logoUrl} alt={businessName} />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-base-200 text-base-content font-bold">
                      {avatarLetter}
                    </div>
                  )}
                </div>
              </div>

              <div className="min-w-0">
                <div className="font-bold leading-tight truncate text-lg">{businessName}</div>
                {catalog?.business?.address ? (
                  <div className="text-sm opacity-70 truncate">{catalog.business.address}</div>
                ) : null}
              </div>
            </div>

            <div className="navbar-end">
              <LanguageSwitcher />
            </div>
          </div>
        </div>
      </header>

      <main className="flex-1">
        <div className="mx-auto max-w-5xl px-4 py-5">
          {catalogQuery.isLoading ? (
            <div className="space-y-3">
              <div className="skeleton h-10 w-full" />
              <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                <div className="skeleton h-48 w-full" />
                <div className="skeleton h-48 w-full" />
                <div className="skeleton h-48 w-full hidden sm:block" />
              </div>
            </div>
          ) : catalogQuery.isError ? (
            <div role="alert" className="alert alert-error alert-soft">
              Failed to load storefront
            </div>
          ) : catalog ? (
            <ProductList catalog={catalog} />
          ) : null}
        </div>
      </main>

      <footer className="border-t border-base-200/60 bg-base-100 safe-bottom">
        <div className="mx-auto max-w-5xl px-4 py-3 flex items-center justify-between gap-3">
          <div className="text-xs opacity-60">{t('poweredBy')}</div>
          {catalog?.business?.supportEmail ? (
            <a className="link link-hover text-xs" href={`mailto:${catalog.business.supportEmail}`}>
              {catalog.business.supportEmail}
            </a>
          ) : null}
        </div>
      </footer>

      {catalog ? <CartDrawer storefrontPublicId={storefrontPublicId} catalog={catalog} /> : null}

      {catalog && totalQty > 0 ? (
        <div className="fixed bottom-0 left-0 right-0 safe-bottom">
          <div className="mx-auto max-w-5xl px-4 pb-3">
            <button type="button" className="btn btn-primary btn-block h-14" onClick={openCart}>
              <span className="flex w-full items-center justify-between gap-3">
                <span className="truncate">
                  {t('viewCart')} · {totalQty} {t('items')}
                </span>
                <span className="shrink-0 font-semibold">
                  {formatMoney(cartSubtotal, catalog.business.currency, i18n.language)}
                </span>
              </span>
            </button>
          </div>
        </div>
      ) : null}
    </div>
  );
}
