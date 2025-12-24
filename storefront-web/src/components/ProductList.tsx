import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Bars3BottomLeftIcon,
  Squares2X2Icon,
} from '@heroicons/react/24/outline';
import type { CatalogResponse, PublicCategory, PublicProduct, PublicVariant } from '../api/types';
import { useCartStore } from '../cart/useCartStore';
import { formatMoney, money } from '../utils/money';
import { ImageTile } from './ImageTile';

type ViewMode = 'grid' | 'list';

function firstPhotoUrl(product: PublicProduct, variant: PublicVariant): string | undefined {
  return variant.photos?.[0] || product.photos?.[0];
}

function pickDefaultVariant(product: PublicProduct): PublicVariant | undefined {
  return product.variants?.[0];
}

export function ProductList({ catalog }: { catalog: CatalogResponse }) {
  const { t, i18n } = useTranslation();
  const addItem = useCartStore((s) => s.addItem);
  const itemsByVariantId = useCartStore((s) => s.itemsByVariantId);

  const [viewMode, setViewMode] = useState<ViewMode>('grid');
  const [activeCategoryId, setActiveCategoryId] = useState<string>('');
  const [selectedVariantByProductId, setSelectedVariantByProductId] = useState<Record<string, string>>({});

  const categoriesWithAll = useMemo(() => {
    const all: PublicCategory = { id: '', name: t('all'), descriptor: 'all' };
    return [all, ...catalog.categories];
  }, [catalog.categories, t]);

  const products = useMemo(() => {
    const list = catalog.products || [];
    if (!activeCategoryId) return list;
    return list.filter((p) => p.categoryId === activeCategoryId);
  }, [catalog.products, activeCategoryId]);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between gap-3">
        <div className="flex-1 min-w-0">
          <div className="overflow-x-auto">
            <div className="tabs tabs-box w-max whitespace-nowrap border border-base-200/60 bg-base-100">
              {categoriesWithAll.map((cat) => (
                <button
                  key={cat.id || 'all'}
                  role="tab"
                  className={`tab ${activeCategoryId === cat.id ? 'tab-active' : ''}`}
                  onClick={() => setActiveCategoryId(cat.id)}
                  type="button"
                >
                  {cat.name}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div className="join border border-base-200/60">
          <button
            type="button"
            className={`btn btn-sm btn-square join-item btn-ghost ${viewMode === 'grid' ? 'btn-active' : ''}`}
            onClick={() => setViewMode('grid')}
            aria-label={t('grid')}
          >
            <Squares2X2Icon className="h-5 w-5" />
          </button>
          <button
            type="button"
            className={`btn btn-sm btn-square join-item btn-ghost ${viewMode === 'list' ? 'btn-active' : ''}`}
            onClick={() => setViewMode('list')}
            aria-label={t('list')}
          >
            <Bars3BottomLeftIcon className="h-5 w-5" />
          </button>
        </div>
      </div>

      <div className={viewMode === 'grid' ? 'grid grid-cols-2 gap-3 sm:grid-cols-3' : 'space-y-3'}>
        {products.map((product) => {
          const variantId = selectedVariantByProductId[product.id] || pickDefaultVariant(product)?.id;
          const variant = product.variants.find((v) => v.id === variantId) || pickDefaultVariant(product);
          if (!variant) return null;

          const qty = itemsByVariantId[variant.id]?.quantity || 0;
          const photo = firstPhotoUrl(product, variant);

          return (
            <div key={product.id}>
              {viewMode === 'grid' ? (
                <div className="card card-compact bg-base-100 border border-base-200/60">
                  <figure className="p-3">
                    <ImageTile src={photo} alt={product.name} aspectClassName="aspect-square" />
                  </figure>

                  <div className="card-body gap-2">
                    <div className="flex items-start justify-between gap-2">
                      <h3 className="card-title text-base leading-tight">{product.name}</h3>
                      <div className="text-sm font-semibold">
                        {formatMoney(money(variant.salePrice), variant.currency, i18n.language)}
                      </div>
                    </div>

                    {product.variants.length > 1 ? (
                      <select
                        className="select select-sm w-full border-base-200"
                        value={variant.id}
                        onChange={(e) =>
                          setSelectedVariantByProductId((prev) => ({ ...prev, [product.id]: e.target.value }))
                        }
                      >
                        {product.variants.map((v) => (
                          <option key={v.id} value={v.id}>
                            {v.name}
                          </option>
                        ))}
                      </select>
                    ) : null}

                    <div className="card-actions justify-end">
                      {qty <= 0 ? (
                        <button
                          type="button"
                          className="btn btn-primary btn-sm"
                          onClick={() =>
                            addItem(
                              {
                                variantId: variant.id,
                                productId: product.id,
                                productName: product.name,
                                variantName: variant.name,
                                unitPrice: variant.salePrice,
                                currency: variant.currency,
                                photoUrl: photo,
                              },
                              1,
                            )
                          }
                        >
                          {t('add')}
                        </button>
                      ) : (
                        <div className="join">
                          <button
                            type="button"
                            className="btn btn-sm join-item"
                            onClick={() => useCartStore.getState().removeItem(variant.id, 1)}
                          >
                            −
                          </button>
                          <div className="btn btn-sm join-item btn-ghost pointer-events-none">{qty}</div>
                          <button
                            type="button"
                            className="btn btn-sm join-item"
                            onClick={() =>
                              addItem(
                                {
                                  variantId: variant.id,
                                  productId: product.id,
                                  productName: product.name,
                                  variantName: variant.name,
                                  unitPrice: variant.salePrice,
                                  currency: variant.currency,
                                  photoUrl: photo,
                                },
                                1,
                              )
                            }
                          >
                            +
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ) : (
                <div className="card card-compact bg-base-100 border border-base-200/60">
                  <div className="card-body p-3">
                  <div className="flex items-start gap-3">
                    <div className="w-20 shrink-0">
                      <ImageTile src={photo} alt={product.name} aspectClassName="aspect-square" />
                    </div>

                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-2">
                        <div className="min-w-0">
                          <div className="font-semibold truncate">{product.name}</div>
                          <div className="text-sm opacity-70 truncate">{variant.name}</div>
                        </div>
                        <div className="text-sm font-semibold whitespace-nowrap">
                          {formatMoney(money(variant.salePrice), variant.currency, i18n.language)}
                        </div>
                      </div>

                      {product.variants.length > 1 ? (
                        <select
                          className="select select-sm w-full mt-2 border-base-200"
                          value={variant.id}
                          onChange={(e) =>
                            setSelectedVariantByProductId((prev) => ({ ...prev, [product.id]: e.target.value }))
                          }
                        >
                          {product.variants.map((v) => (
                            <option key={v.id} value={v.id}>
                              {v.name}
                            </option>
                          ))}
                        </select>
                      ) : null}

                      <div className="mt-3 flex items-center justify-end">
                        {qty <= 0 ? (
                          <button
                            type="button"
                            className="btn btn-primary btn-sm"
                            onClick={() =>
                              addItem(
                                {
                                  variantId: variant.id,
                                  productId: product.id,
                                  productName: product.name,
                                  variantName: variant.name,
                                  unitPrice: variant.salePrice,
                                  currency: variant.currency,
                                  photoUrl: photo,
                                },
                                1,
                              )
                            }
                          >
                            {t('add')}
                          </button>
                        ) : (
                          <div className="join">
                            <button
                              type="button"
                              className="btn btn-sm join-item"
                              onClick={() => useCartStore.getState().removeItem(variant.id, 1)}
                            >
                              −
                            </button>
                            <div className="btn btn-sm join-item btn-ghost pointer-events-none">{qty}</div>
                            <button
                              type="button"
                              className="btn btn-sm join-item"
                              onClick={() =>
                                addItem(
                                  {
                                    variantId: variant.id,
                                    productId: product.id,
                                    productName: product.name,
                                    variantName: variant.name,
                                    unitPrice: variant.salePrice,
                                    currency: variant.currency,
                                    photoUrl: photo,
                                  },
                                  1,
                                )
                              }
                            >
                              +
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
