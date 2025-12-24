import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import {
  MinusIcon,
  PlusIcon,
  TrashIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';

import type {
  CatalogResponse,
  CreateOrderRequest,
  PublicBusiness,
  PublicShippingZone,
} from '../api/types';
import { storefrontApi } from '../api/storefront';
import {
  cartItemsArray,
  cartTotalQuantity,
  useCartStore,
} from '../cart/useCartStore';
import { formatMoney, money } from '../utils/money';
import { ImageTile } from './ImageTile';
import { WhatsAppButton } from './WhatsAppButton';

function getDefaultWhatsappNumber(business: PublicBusiness): string {
  return business.whatsappNumber || business.phoneNumber || '';
}

function firstNonEmptyString(values: Array<string | undefined | null>): string {
  for (const v of values) {
    if (typeof v === 'string' && v.trim()) return v;
  }
  return '';
}

export function CartDrawer(props: {
  storefrontPublicId: string;
  catalog: CatalogResponse;
}) {
  const { t, i18n } = useTranslation();

  const isOpen = useCartStore((s) => s.isCartOpen);
  const closeCart = useCartStore((s) => s.closeCart);
  const clearCart = useCartStore((s) => s.clearCart);
  const itemsByVariantId = useCartStore((s) => s.itemsByVariantId);

  const items = useMemo(
    () => cartItemsArray({ itemsByVariantId }),
    [itemsByVariantId],
  );
  const qty = useMemo(
    () => cartTotalQuantity({ itemsByVariantId }),
    [itemsByVariantId],
  );

  const zonesQuery = useQuery({
    queryKey: ['shipping-zones', props.storefrontPublicId],
    queryFn: () => storefrontApi.listShippingZones(props.storefrontPublicId),
    enabled: isOpen,
  });

  const zones: PublicShippingZone[] = useMemo(
    () => zonesQuery.data || [],
    [zonesQuery.data],
  );

  const [shippingZoneId, setShippingZoneId] = useState<string>('');
  const selectedZone = useMemo(() => {
    if (!zones.length) return undefined;
    return zones.find((z) => z.id === shippingZoneId) || zones[0];
  }, [zones, shippingZoneId]);

  const availableCountries = useMemo(
    () => selectedZone?.countries || [],
    [selectedZone],
  );

  const [countryCode, setCountryCode] = useState<string>(
    firstNonEmptyString([
      availableCountries[0],
      props.catalog.business.countryCode,
    ]),
  );

  const effectiveCountryCode = useMemo(() => {
    if (availableCountries.length === 0) return countryCode;
    if (availableCountries.includes(countryCode)) return countryCode;
    return availableCountries[0];
  }, [availableCountries, countryCode]);

  const [customerName, setCustomerName] = useState('');
  const [customerEmail, setCustomerEmail] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [instagramUsername, setInstagramUsername] = useState('');

  const [stateField, setStateField] = useState('');
  const [city, setCity] = useState('');
  const [street, setStreet] = useState('');
  const [zipCode, setZipCode] = useState('');
  const [phoneCode, setPhoneCode] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');

  const currency = props.catalog.business.currency;

  const subtotal = useMemo(() => {
    return items.reduce(
      (acc, x) => acc.add(money(x.unitPrice).mul(x.quantity)),
      money('0'),
    );
  }, [items]);

  const shipping = useMemo(() => {
    if (!selectedZone) return money('0');
    const cost = money(selectedZone.shippingCost || '0');
    const threshold = money(selectedZone.freeShippingThreshold || '0');
    if (
      threshold.greaterThan(0) &&
      subtotal.greaterThanOrEqualTo(threshold)
    ) {
      return money('0');
    }
    return cost;
  }, [selectedZone, subtotal]);

  const total = useMemo(() => subtotal.add(shipping), [subtotal, shipping]);

  const request: CreateOrderRequest = useMemo(
    () => ({
      customer: {
        email: customerEmail.trim(),
        name: customerName.trim(),
        phoneNumber: customerPhone.trim() || undefined,
        instagramUsername: instagramUsername.trim() || undefined,
      },
      shippingAddress: {
        countryCode: (effectiveCountryCode || '').trim(),
        state: stateField.trim(),
        city: city.trim(),
        street: street.trim() || undefined,
        zipCode: zipCode.trim() || undefined,
        phoneCode: phoneCode.trim(),
        phoneNumber: phoneNumber.trim(),
      },
      items: items.map((it) => ({
        variantId: it.variantId,
        quantity: it.quantity,
        specialRequest: (it.note || '').trim() || undefined,
      })),
    }),
    [
      items,
      customerEmail,
      customerName,
      customerPhone,
      instagramUsername,
      effectiveCountryCode,
      stateField,
      city,
      street,
      zipCode,
      phoneCode,
      phoneNumber,
    ],
  );

  const whatsappNumber = getDefaultWhatsappNumber(props.catalog.business);

  return (
    <dialog
      className={`modal modal-bottom sm:modal-middle ${
        isOpen ? 'modal-open' : ''
      }`}
    >
      <div className="modal-box p-0 overflow-hidden rounded-t-box sm:rounded-box flex flex-col max-h-[85vh]">
        <div className="p-4 border-b border-base-200/60 flex items-center justify-between gap-2">
          <button
            type="button"
            className="btn btn-ghost btn-square"
            onClick={closeCart}
            aria-label={t('close')}
          >
            <XMarkIcon className="h-6 w-6" />
          </button>

          <div className="min-w-0 flex-1 text-center">
            <div className="font-bold text-lg">{t('cart')}</div>
            <div className="text-sm opacity-70">
              {qty} {t('items')}
            </div>
          </div>

          {items.length > 0 ? (
            <button
              type="button"
              className="btn btn-ghost btn-square"
              onClick={clearCart}
              aria-label={t('clearCart')}
            >
              <TrashIcon className="h-6 w-6" />
            </button>
          ) : (
            <div className="w-12" />
          )}
        </div>

        <div className="p-4 space-y-4 overflow-y-auto flex-1">
          {items.length === 0 ? (
            <div className="text-center opacity-70 py-10">{t('emptyCart')}</div>
          ) : (
            <>
              <div className="space-y-3">
                {items.map((it) => (
                  <div
                    key={it.variantId}
                    className="card card-compact bg-base-100 border border-base-200/60"
                  >
                    <div className="card-body gap-3">
                      <div className="flex items-start gap-3">
                        <div className="w-14 shrink-0">
                          <ImageTile
                            src={it.photoUrl}
                            alt={it.productName}
                            aspectClassName="aspect-square"
                          />
                        </div>

                        <div className="min-w-0 flex-1">
                          <div className="flex items-start justify-between gap-3">
                            <div className="min-w-0">
                              <div className="font-semibold truncate">{it.productName}</div>
                              <div className="text-sm opacity-70 truncate">{it.variantName}</div>
                            </div>
                            <div className="text-sm font-semibold whitespace-nowrap">
                              {formatMoney(
                                money(it.unitPrice).mul(it.quantity),
                                it.currency,
                                i18n.language,
                              )}
                            </div>
                          </div>

                          <div className="mt-3 flex items-center justify-between gap-2">
                            <div className="join">
                              <button
                                type="button"
                                className="btn btn-sm btn-square join-item"
                                onClick={() => useCartStore.getState().removeItem(it.variantId, 1)}
                                aria-label={t('decrease')}
                              >
                                <MinusIcon className="h-5 w-5" />
                              </button>
                              <div className="btn btn-sm join-item btn-ghost pointer-events-none">
                                {it.quantity}
                              </div>
                              <button
                                type="button"
                                className="btn btn-sm btn-square join-item"
                                onClick={() =>
                                  useCartStore.getState().addItem(
                                    {
                                      variantId: it.variantId,
                                      productId: it.productId,
                                      productName: it.productName,
                                      variantName: it.variantName,
                                      unitPrice: it.unitPrice,
                                      currency: it.currency,
                                      photoUrl: it.photoUrl,
                                    },
                                    1,
                                  )
                                }
                                aria-label={t('increase')}
                              >
                                <PlusIcon className="h-5 w-5" />
                              </button>
                            </div>

                            <button
                              type="button"
                              className="btn btn-ghost btn-sm btn-square"
                              onClick={() => useCartStore.getState().setQuantity(it.variantId, 0)}
                              aria-label={t('remove')}
                            >
                              <TrashIcon className="h-5 w-5" />
                            </button>
                          </div>

                          <input
                            className="input input-bordered border-base-200 input-sm w-full mt-3"
                            value={it.note || ''}
                            onChange={(e) =>
                              useCartStore.getState().updateNote(it.variantId, e.target.value)
                            }
                            placeholder={`${t('specialRequest')} (${t('optional')})`}
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <fieldset className="fieldset border border-base-200/60 rounded-box p-4">
                <legend className="fieldset-legend">{t('shippingZone')}</legend>

                <select
                  className="select select-bordered border-base-200 w-full"
                  value={selectedZone?.id || ''}
                  onChange={(e) => {
                    const id = e.target.value;
                    setShippingZoneId(id);

                    const zone = zones.find((z) => z.id === id);
                    if (!zone || !zone.countries?.length) return;
                    if (zone.countries.includes(countryCode)) return;
                    setCountryCode(zone.countries[0]);
                  }}
                  disabled={zonesQuery.isLoading || zones.length === 0}
                >
                  {zones.length === 0 ? (
                    <option value="">{t('select')}</option>
                  ) : (
                    zones.map((z) => (
                      <option key={z.id} value={z.id}>
                        {z.name}
                      </option>
                    ))
                  )}
                </select>

                {availableCountries.length > 0 ? (
                  <select
                    className="select select-bordered border-base-200 w-full"
                    value={effectiveCountryCode}
                    onChange={(e) => setCountryCode(e.target.value)}
                  >
                    {availableCountries.map((cc) => (
                      <option key={cc} value={cc}>
                        {cc}
                      </option>
                    ))}
                  </select>
                ) : (
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={effectiveCountryCode}
                    onChange={(e) => setCountryCode(e.target.value)}
                    placeholder={t('countryCode')}
                  />
                )}

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={stateField}
                    onChange={(e) => setStateField(e.target.value)}
                    placeholder={t('state')}
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={city}
                    onChange={(e) => setCity(e.target.value)}
                    placeholder={t('city')}
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={street}
                    onChange={(e) => setStreet(e.target.value)}
                    placeholder={t('street')}
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={zipCode}
                    onChange={(e) => setZipCode(e.target.value)}
                    placeholder={t('zipCode')}
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={phoneCode}
                    onChange={(e) => setPhoneCode(e.target.value)}
                    placeholder={t('phoneCode')}
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                    placeholder={t('phoneNumber')}
                  />
                </div>
              </fieldset>

              <fieldset className="fieldset border border-base-200/60 rounded-box p-4">
                <legend className="fieldset-legend">{t('yourDetails')}</legend>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={customerName}
                    onChange={(e) => setCustomerName(e.target.value)}
                    placeholder={t('name')}
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={customerEmail}
                    onChange={(e) => setCustomerEmail(e.target.value)}
                    placeholder={t('email')}
                    inputMode="email"
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={customerPhone}
                    onChange={(e) => setCustomerPhone(e.target.value)}
                    placeholder={t('phone')}
                    inputMode="tel"
                  />
                  <input
                    className="input input-bordered border-base-200 w-full"
                    value={instagramUsername}
                    onChange={(e) => setInstagramUsername(e.target.value)}
                    placeholder={t('instagramUsername')}
                  />
                </div>
              </fieldset>
            </>
          )}
        </div>

        {items.length > 0 ? (
          <div className="border-t border-base-200/60 bg-base-100 p-4 space-y-3">
            <div className="card card-compact bg-base-100 border border-base-200/60">
              <div className="card-body gap-2">
                <div className="flex justify-between text-sm">
                  <span className="opacity-70">{t('subtotal')}</span>
                  <span className="font-semibold">
                    {formatMoney(subtotal, currency, i18n.language)}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="opacity-70">{t('shipping')}</span>
                  <span className="font-semibold">
                    {shipping.greaterThan(0)
                      ? formatMoney(shipping, currency, i18n.language)
                      : t('free')}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="opacity-70">{t('total')}</span>
                  <span className="font-bold">
                    {formatMoney(total, currency, i18n.language)}
                  </span>
                </div>
              </div>
            </div>

            <WhatsAppButton
              storefrontPublicId={props.storefrontPublicId}
              business={props.catalog.business}
              request={request}
              whatsappNumber={whatsappNumber}
              displayItems={items.map((it) => ({
                variantId: it.variantId,
                title: `${it.productName} — ${it.variantName}`,
                quantity: it.quantity,
                note: it.note,
              }))}
              onSuccess={() => {
                clearCart();
                closeCart();
              }}
            />
          </div>
        ) : null}
      </div>

      <form method="dialog" className="modal-backdrop">
        <button onClick={closeCart}>{t('close')}</button>
      </form>
    </dialog>
  );
}
