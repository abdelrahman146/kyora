import { memo, useMemo, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { PackageOpen, Grid3x3, List } from 'lucide-react';
import type { CatalogResponse, PublicCategory, PublicProduct, PublicVariant } from '../../api/types';
import { useCartStore } from '../../cart/useCartStore';
import { ProductCard, ProductCardSkeleton, ProductCardList } from '../molecules';
import { EmptyState } from '../atoms';

function firstPhotoUrl(product: PublicProduct, variant: PublicVariant): string | undefined {
  return variant.photos?.[0] || product.photos?.[0];
}

function pickDefaultVariant(product: PublicProduct): PublicVariant | undefined {
  return product.variants?.[0];
}

interface ProductListProps {
  catalog: CatalogResponse;
  isLoading?: boolean;
}

/**
 * ProductListNew Organism - Product grid with category filtering
 * Memoized to prevent unnecessary re-renders
 * Optimized with useMemo and useCallback for performance
 * Follows KDS 5.2 with 2-column mobile grid and horizontal category tabs
 */
export const ProductListNew = memo<ProductListProps>(function ProductListNew({ catalog, isLoading }) {
  const { t } = useTranslation();
  const addItem = useCartStore((s) => s.addItem);
  const removeItem = useCartStore((s) => s.removeItem);
  const itemsByVariantId = useCartStore((s) => s.itemsByVariantId);

  const [activeCategoryId, setActiveCategoryId] = useState<string>('');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [selectedVariantByProductId, setSelectedVariantByProductId] = useState<
    Record<string, string>
  >({});

  const categoriesWithAll = useMemo(() => {
    const all: PublicCategory = { id: '', name: t('all'), descriptor: 'all' };
    return [all, ...catalog.categories];
  }, [catalog.categories, t]);

  const products = useMemo(() => {
    const list = catalog.products || [];
    if (!activeCategoryId) return list;
    return list.filter((p) => p.categoryId === activeCategoryId);
  }, [catalog.products, activeCategoryId]);

  const handleAddToCart = (product: PublicProduct, variant: PublicVariant) => {
    addItem(
      {
        variantId: variant.id,
        productId: product.id,
        productName: product.name,
        variantName: variant.name,
        unitPrice: variant.salePrice,
        currency: variant.currency,
        photoUrl: firstPhotoUrl(product, variant),
      },
      1
    );
  };

  // Memoize event handlers to prevent re-creating functions
  const handleIncrement = useCallback((product: PublicProduct, variant: PublicVariant) => {
    addItem(
      {
        variantId: variant.id,
        productId: product.id,
        productName: product.name,
        variantName: variant.name,
        unitPrice: variant.salePrice,
        currency: variant.currency,
        photoUrl: firstPhotoUrl(product, variant),
      },
      1
    );
  }, [addItem]);

  const handleDecrement = useCallback((variantId: string) => {
    removeItem(variantId, 1);
  }, [removeItem]);

  const handleVariantChange = useCallback((productId: string, variantId: string) => {
    setSelectedVariantByProductId((prev) => ({ ...prev, [productId]: variantId }));
  }, []);

  return (
    <div className="space-y-6">
      {/* Category Tabs and View Toggle */}
      <div className="space-y-4">
        {/* Category Tabs - Horizontal Scroll */}
        <div className="overflow-x-auto scrollbar-hide -mx-4 px-4 md:-mx-6 md:px-6">
          <div className="flex gap-2 min-w-max pb-2">
          {categoriesWithAll.map((cat) => (
            <button
              key={cat.id || 'all'}
              type="button"
              onClick={() => setActiveCategoryId(cat.id)}
              className={`px-4 py-2 rounded-xl font-semibold text-xs whitespace-nowrap transition-all active-scale focus-ring ${
                activeCategoryId === cat.id
                  ? 'bg-primary text-primary-content shadow-md'
                  : 'bg-base-100 text-base-content border border-base-300 hover:border-primary hover:text-primary hover:shadow-sm'
              }`}
            >
              {cat.name}
            </button>
          ))}
          </div>
        </div>

        {/* View Toggle */}
        <div className="flex items-center justify-end gap-2">
          <button
            type="button"
            onClick={() => setViewMode('grid')}
            className={`btn btn-sm btn-square shadow-sm ${
              viewMode === 'grid'
                ? 'btn-primary'
                : 'bg-base-100 border border-base-300 hover:border-primary hover:bg-base-100'
            }`}
            aria-label={t('gridView')}
          >
            <Grid3x3 className="w-4 h-4" strokeWidth={2} />
          </button>
          <button
            type="button"
            onClick={() => setViewMode('list')}
            className={`btn btn-sm btn-square shadow-sm ${
              viewMode === 'list'
                ? 'btn-primary'
                : 'bg-base-100 border border-base-300 hover:border-primary hover:bg-base-100'
            }`}
            aria-label={t('listView')}
          >
            <List className="w-4 h-4" strokeWidth={2} />
          </button>
        </div>
      </div>

      {/* Products Display */}
      {isLoading ? (
        <div className={viewMode === 'grid' ? 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4' : 'space-y-4'}>
          {Array.from({ length: 6 }).map((_, i) => (
            <ProductCardSkeleton key={i} />
          ))}
        </div>
      ) : products.length === 0 ? (
        <EmptyState
          icon={<PackageOpen className="w-16 h-16" strokeWidth={1.5} />}
          title={t('emptyCart')}
          description={activeCategoryId ? 'No products in this category' : 'No products available'}
        />
      ) : (
        <div className={viewMode === 'grid' ? 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4' : 'space-y-4'}>
          {products.map((product) => {
            const variantId =
              selectedVariantByProductId[product.id] || pickDefaultVariant(product)?.id;
            const variant =
              product.variants.find((v) => v.id === variantId) || pickDefaultVariant(product);
            if (!variant) return null;

            const quantity = itemsByVariantId[variant.id]?.quantity || 0;

            const CardComponent = viewMode === 'grid' ? ProductCard : ProductCardList;

            return (
              <CardComponent
                key={product.id}
                product={product}
                variant={variant}
                quantity={quantity}
                onAddToCart={() => handleAddToCart(product, variant)}
                onIncrement={() => handleIncrement(product, variant)}
                onDecrement={() => handleDecrement(variant.id)}
                onVariantChange={
                  product.variants.length > 1
                    ? (vid) => handleVariantChange(product.id, vid)
                    : undefined
                }
              />
            );
          })}
        </div>
      )}
    </div>
  );
});
