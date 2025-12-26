import { useMemo, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Trash2, PackageOpen } from 'lucide-react';
import type {
  CatalogResponse,
  CreateOrderRequest,
  PublicBusiness,
  PublicShippingZone,
} from '../../api/types';
import { storefrontApi } from '../../api/storefront';
import { cartItemsArray, cartTotalQuantity, useCartStore } from '../../cart/useCartStore';
import { money } from '../../utils/money';
import { CartItem, FormInput, FormSelect, OrderSummary } from '../molecules';
import { BottomSheet } from './BottomSheet';
import { WhatsAppButton } from './WhatsAppButton';

function getDefaultWhatsappNumber(business: PublicBusiness): string {
  return business.whatsappNumber || business.phoneNumber || '';
}

function firstNonEmpty(...values: Array<string | undefined | null>): string {
  for (const v of values) {
    if (typeof v === 'string' && v.trim()) return v;
  }
  return '';
}

const checkoutSchema = z.object({
  shippingZoneId: z.string().min(1, 'Required'),
  countryCode: z.string().min(1, 'Required'),
  state: z.string().min(1, 'Required'),
  city: z.string().min(1, 'Required'),
  street: z.string().optional(),
  zipCode: z.string().optional(),
  phoneCode: z.string().min(1, 'Required'),
  phoneNumber: z.string().min(1, 'Required'),
  customerName: z.string().min(1, 'Required'),
  customerEmail: z.string().email('Invalid email'),
  customerPhone: z.string().optional(),
  instagramUsername: z.string().optional(),
});

type CheckoutFormData = z.infer<typeof checkoutSchema>;

export function CartDrawer(props: { storefrontPublicId: string; catalog: CatalogResponse }) {
  const { t } = useTranslation();

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

  const { register, watch, setValue, formState: { errors } } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      shippingZoneId: '',
      countryCode: firstNonEmpty(props.catalog.business.countryCode),
      state: '',
      city: '',
      street: '',
      zipCode: '',
      phoneCode: '',
      phoneNumber: '',
      customerName: '',
      customerEmail: '',
      customerPhone: '',
      instagramUsername: '',
    },
  });

  const formValues = watch();
  
  const selectedZone = useMemo(() => {
    if (!zones.length) return undefined;
    return zones.find((z) => z.id === formValues.shippingZoneId) || zones[0];
  }, [zones, formValues.shippingZoneId]);

  const availableCountries = useMemo(
    () => selectedZone?.countries || [],
    [selectedZone],
  );

  const effectiveCountryCode = useMemo(() => {
    if (availableCountries.length === 0) return formValues.countryCode;
    if (availableCountries.includes(formValues.countryCode)) return formValues.countryCode;
    return availableCountries[0];
  }, [availableCountries, formValues.countryCode]);

  // Auto-select first zone
  useEffect(() => {
    if (zones.length > 0 && !formValues.shippingZoneId) {
      setValue('shippingZoneId', zones[0].id);
    }
  }, [zones, formValues.shippingZoneId, setValue]);

  // Auto-correct country code
  useEffect(() => {
    if (availableCountries.length > 0 && !availableCountries.includes(formValues.countryCode)) {
      setValue('countryCode', availableCountries[0]);
    }
  }, [availableCountries, formValues.countryCode, setValue]);

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
        email: formValues.customerEmail.trim(),
        name: formValues.customerName.trim(),
        phoneNumber: formValues.customerPhone?.trim() || undefined,
        instagramUsername: formValues.instagramUsername?.trim() || undefined,
      },
      shippingAddress: {
        countryCode: effectiveCountryCode.trim(),
        state: formValues.state.trim(),
        city: formValues.city.trim(),
        street: formValues.street?.trim() || undefined,
        zipCode: formValues.zipCode?.trim() || undefined,
        phoneCode: formValues.phoneCode.trim(),
        phoneNumber: formValues.phoneNumber.trim(),
      },
      items: items.map((it) => ({
        variantId: it.variantId,
        quantity: it.quantity,
        specialRequest: (it.note || '').trim() || undefined,
      })),
    }),
    [items, formValues, effectiveCountryCode],
  );

  const whatsappNumber = getDefaultWhatsappNumber(props.catalog.business);

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={closeCart}
      title={t('cart')}
    >
      {/* Cart Content */}
      <div className="flex flex-col gap-4">
        {/* Header with item count and clear button */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-neutral-500">
              {qty} {qty === 1 ? t('item') : t('items')}
            </p>
          </div>
          {items.length > 0 && (
            <button
              type="button"
              className="btn btn-ghost btn-sm gap-2 active-scale focus-ring text-error hover:bg-error hover:text-error-content"
              onClick={clearCart}
            >
              <Trash2 className="w-4 h-4" strokeWidth={2} />
              <span className="text-sm font-medium">{t('clearCart')}</span>
            </button>
          )}
        </div>

        {/* Cart Items */}
        <div className="space-y-3">
          {items.length === 0 ? (
            <div className="text-center text-neutral-500 py-12">
              <PackageOpen className="w-16 h-16 mx-auto mb-4 opacity-30" strokeWidth={1.5} />
              <p className="text-base font-medium">{t('emptyCart')}</p>
            </div>
          ) : (
            <>
              {items.map((it) => (
                <CartItem key={it.variantId} {...it} />
              ))}

              {/* Shipping Zone Section */}
              <div className="bg-neutral-50 rounded-lg border border-neutral-200 p-4 space-y-3">
                <h3 className="font-semibold text-sm text-neutral-900">{t('shippingZone')}</h3>

                <FormSelect
                  {...register('shippingZoneId')}
                  error={errors.shippingZoneId}
                  disabled={zonesQuery.isLoading || zones.length === 0}
                  options={zones.length === 0 ? [{ value: '', label: t('select') }] : zones.map((z) => ({ value: z.id, label: z.name }))}
                  onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                    register('shippingZoneId').onChange(e);
                    const zone = zones.find((z) => z.id === e.target.value);
                    if (zone?.countries?.length && !zone.countries.includes(formValues.countryCode)) {
                      setValue('countryCode', zone.countries[0]);
                    }
                  }}
                />

                {availableCountries.length > 0 ? (
                  <FormSelect
                    {...register('countryCode')}
                    error={errors.countryCode}
                    options={availableCountries.map((cc) => ({ value: cc, label: cc }))}
                  />
                ) : (
                  <FormInput
                    {...register('countryCode')}
                    error={errors.countryCode}
                    placeholder={t('countryCode')}
                  />
                )}

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <FormInput {...register('state')} error={errors.state} placeholder={t('state')} />
                  <FormInput {...register('city')} error={errors.city} placeholder={t('city')} />
                  <FormInput {...register('street')} error={errors.street} placeholder={t('street')} />
                  <FormInput {...register('zipCode')} error={errors.zipCode} placeholder={t('zipCode')} />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <FormInput {...register('phoneCode')} error={errors.phoneCode} placeholder={t('phoneCode')} />
                  <FormInput {...register('phoneNumber')} error={errors.phoneNumber} placeholder={t('phoneNumber')} />
                </div>
              </div>

              {/* Customer Details Section */}
              <div className="bg-neutral-50 rounded-lg border border-neutral-200 p-4 space-y-3">
                <h3 className="font-semibold text-sm text-neutral-900">{t('yourDetails')}</h3>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  <FormInput {...register('customerName')} error={errors.customerName} placeholder={t('name')} />
                  <FormInput {...register('customerEmail')} error={errors.customerEmail} placeholder={t('email')} inputMode="email" />
                  <FormInput {...register('customerPhone')} error={errors.customerPhone} placeholder={t('phone')} inputMode="tel" />
                  <FormInput {...register('instagramUsername')} error={errors.instagramUsername} placeholder={t('instagramUsername')} />
                </div>
              </div>
            </>
          )}
        </div>

        {/* Order Summary and WhatsApp Button */}
        {items.length > 0 && (
          <div className="bg-white border-t border-neutral-200 p-4 space-y-4">
            <OrderSummary subtotal={subtotal} shipping={shipping} total={total} currency={currency} />

            {/* WhatsApp Checkout Button */}
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
        )}
      </div>
      {/* End of flex-col container */}
    </BottomSheet>
  );
}
