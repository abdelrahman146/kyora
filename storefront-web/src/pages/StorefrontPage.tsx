import { useEffect, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { storefrontApi } from '../api/storefront';
import { applyBrandTheme } from '../theme/applyBrandTheme';
import { StickyCartBar } from '../components/molecules';
import { ProductListNew as ProductList, CartDrawer, StorefrontHeader } from '../components/organisms';
import { BrandHeader } from '../components/atoms';
import { cartTotalQuantity, useCartStore } from '../cart/useCartStore';
import { money } from '../utils/money';

/**
 * StorefrontPage - Complete revamp following Kyora Design System
 * - Mobile-first PWA with native app feel
 * - RTL-first design (Arabic as first-class citizen)
 * - Sticky header with scroll shadow
 * - 2-column product grid
 * - Sticky cart bar at bottom
 * - Bottom sheet for cart
 * - Floating WhatsApp button
 * - Proper loading states and skeletons
 * - Full accessibility support
 */
export function StorefrontPage() {
  const { storefrontPublicId } = useParams();
  const { t } = useTranslation();

  const openCart = useCartStore((s) => s.openCart);
  const itemsByVariantId = useCartStore((s) => s.itemsByVariantId);
  const totalQty = useMemo(
    () => cartTotalQuantity({ itemsByVariantId }),
    [itemsByVariantId]
  );

  const catalogQuery = useQuery({
    queryKey: ['catalog', storefrontPublicId],
    queryFn: () => storefrontApi.getCatalog(storefrontPublicId || ''),
    enabled: !!storefrontPublicId,
  });

  const catalog = catalogQuery.data;

  useEffect(() => {
    applyBrandTheme(catalog?.business?.storefrontTheme);
  }, [catalog?.business?.storefrontTheme]);

  useEffect(() => {
    // Set page title dynamically
    if (catalog?.business?.name) {
      document.title = `${catalog.business.name} - Kyora Storefront`;
    }
  }, [catalog?.business?.name]);

  const cartSubtotal = useMemo(() => {
    const items = Object.values(itemsByVariantId);
    return items.reduce(
      (acc, it) => acc.add(money(it.unitPrice).mul(it.quantity)),
      money('0')
    );
  }, [itemsByVariantId]);

  if (!storefrontPublicId) {
    return (
      <div className="min-h-dvh flex items-center justify-center p-6 bg-base-200">
        <div className="text-center">
          <div className="text-lg font-semibold text-base-content">
            {t('store')} ID required
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-dvh flex flex-col bg-gradient-to-b from-base-200 via-base-100 to-base-200 text-base-content">
      <div className="bg-base-100/80 backdrop-blur-sm border-b border-base-300/50 mb-4">
        {/* Minimal Transparent Header */}
      <StorefrontHeader
        cartItemCount={totalQty}
        onCartClick={openCart}
      />
      {/* Brand Header (Centered Logo + Name) */}
      {catalog?.business && <BrandHeader business={catalog.business} />}
      </div>

      {/* Main Content with consistent spacing */}
      <main className="flex-1 pb-32">
        <div className="mx-auto max-w-5xl px-4 md:px-6">
          
          {catalogQuery.isLoading ? (
            <ProductList 
              catalog={{ 
                products: [], 
                categories: [], 
                business: { 
                  id: '', 
                  name: '', 
                  logoUrl: '',
                  currency: 'USD',
                  countryCode: '',
                  descriptor: '',
                  storefrontPublicId: '',
                  storefrontEnabled: false,
                  storefrontTheme: {}
                } 
              }} 
              isLoading 
            />
          ) : catalogQuery.isError ? (
            <div role="alert" className="alert alert-error rounded-lg">
              <div>
                <div className="font-semibold">Failed to load storefront</div>
                <div className="text-sm">Please try again later</div>
              </div>
            </div>
          ) : catalog ? (
            <ProductList catalog={catalog} />
          ) : null}
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-base-300/50 bg-base-100/80 backdrop-blur-sm safe-bottom">
        <div className="mx-auto max-w-5xl px-4 md:px-6 py-6 text-center">
          <div className="text-xs text-base-content/60">{t('poweredBy')}</div>
        </div>
      </footer>

      {/* Cart Drawer/Sheet */}
      {catalog && <CartDrawer storefrontPublicId={storefrontPublicId} catalog={catalog} />}

      {/* Sticky Cart Bar (only when items > 0) */}
      {catalog && totalQty > 0 && (
        <StickyCartBar
          itemCount={totalQty}
          totalAmount={cartSubtotal}
          currency={catalog.business.currency}
          onViewCart={openCart}
        />
      )}


    </div>
  );
}
