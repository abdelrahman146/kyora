import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { QRCodeCanvas } from 'qrcode.react';
import type { CreateOrderRequest, CreateOrderResponse, PublicBusiness } from '../../api/types';
import { storefrontApi } from '../../api/storefront';
import { ApiError } from '../../api/client';
import { buildWhatsAppLink, isProbablyMobileDevice } from '../../utils/whatsapp';

function randomIdempotencyKey(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (crypto as any).randomUUID();
  }
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

function buildMessage(params: {
  businessName: string;
  order: CreateOrderResponse;
  request: CreateOrderRequest;
  displayItems?: Array<{ title: string; quantity: number; note?: string; variantId: string }>;
}): string {
  const lines: string[] = [];

  lines.push(`Order #${params.order.orderNumber}`);
  lines.push(params.businessName);
  lines.push('');

  lines.push('Items:');
  const display = params.displayItems;
  if (display && display.length > 0) {
    for (const item of display) {
      const note = (item.note || '').trim();
      lines.push(`- ${item.quantity} × ${item.title}${note ? ` (${note})` : ''}`);
    }
  } else {
    for (const item of params.request.items) {
      const note = (item.specialRequest || '').trim();
      lines.push(`- ${item.quantity} × ${item.variantId}${note ? ` (${note})` : ''}`);
    }
  }

  lines.push('');
  lines.push(`Total: ${params.order.total} ${params.order.currency}`);
  lines.push('');

  lines.push('Customer:');
  lines.push(`- Name: ${params.request.customer.name}`);
  lines.push(`- Email: ${params.request.customer.email}`);
  if (params.request.customer.phoneNumber) lines.push(`- Phone: ${params.request.customer.phoneNumber}`);
  if (params.request.customer.instagramUsername)
    lines.push(`- Instagram: ${params.request.customer.instagramUsername}`);

  lines.push('');
  lines.push('Shipping:');
  lines.push(`- Country: ${params.request.shippingAddress.countryCode}`);
  lines.push(`- State: ${params.request.shippingAddress.state}`);
  lines.push(`- City: ${params.request.shippingAddress.city}`);
  if (params.request.shippingAddress.street) lines.push(`- Street: ${params.request.shippingAddress.street}`);
  if (params.request.shippingAddress.zipCode) lines.push(`- ZIP: ${params.request.shippingAddress.zipCode}`);
  lines.push(`- Phone: ${params.request.shippingAddress.phoneCode}${params.request.shippingAddress.phoneNumber}`);

  return lines.join('\n');
}

export function WhatsAppButton(props: {
  storefrontPublicId: string;
  business: PublicBusiness;
  request: CreateOrderRequest;
  whatsappNumber?: string;
  displayItems?: Array<{ title: string; quantity: number; note?: string; variantId: string }>;
  onSuccess?: () => void;
}) {
  const { t } = useTranslation();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string>('');

  const [qrOpen, setQrOpen] = useState(false);
  const [whatsAppUrl, setWhatsAppUrl] = useState<string>('');
  const [message, setMessage] = useState<string>('');
  const [copied, setCopied] = useState(false);

  const canSubmit = useMemo(() => {
    const r = props.request;
    return (
      r.items.length > 0 &&
      !!r.customer.name &&
      !!r.customer.email &&
      !!r.shippingAddress.countryCode &&
      !!r.shippingAddress.state &&
      !!r.shippingAddress.city &&
      !!r.shippingAddress.phoneCode &&
      !!r.shippingAddress.phoneNumber
    );
  }, [props.request]);

  async function onClick() {
    setError('');
    setCopied(false);

    const number = (props.whatsappNumber || '').trim();
    if (!number) {
      setError(t('whatsappNotConfigured'));
      return;
    }

    if (!canSubmit) {
      return;
    }

    setIsSubmitting(true);
    try {
      const idempotencyKey = randomIdempotencyKey();
      const order = await storefrontApi.createOrder(props.storefrontPublicId, props.request, idempotencyKey);

      const msg = buildMessage({
        businessName: props.business.name,
        order,
        request: props.request,
        displayItems: props.displayItems,
      });

      const url = buildWhatsAppLink(number, msg);

      setMessage(msg);
      setWhatsAppUrl(url);

      if (isProbablyMobileDevice()) {
        window.location.href = url;
        props.onSuccess?.();
      } else {
        setQrOpen(true);
        props.onSuccess?.();
      }
    } catch (e) {
      if (e instanceof ApiError) {
        setError(e.problem?.detail || e.problem?.title || `HTTP ${e.status}`);
      } else {
        setError('Failed to create order');
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(message);
      setCopied(true);
      setTimeout(() => setCopied(false), 1200);
    } catch {
      // ignore
    }
  }

  return (
    <>
      {error ? <div role="alert" className="alert alert-error alert-soft">{error}</div> : null}

      <button
        type="button"
        className="btn btn-primary btn-block btn-lg"
        onClick={onClick}
        disabled={isSubmitting || !canSubmit}
      >
        {isSubmitting ? t('creatingOrder') : t('sendOrderOnWhatsapp')}
      </button>

      <dialog className={`modal ${qrOpen ? 'modal-open' : ''}`} onClose={() => setQrOpen(false)}>
        <div className="modal-box">
          <h3 className="font-bold text-lg mb-4">{t('scanQr')}</h3>

          <div className="flex justify-center">
            <div className="p-3 bg-base-100 rounded-box border">
              <QRCodeCanvas value={whatsAppUrl || 'https://wa.me/'} size={220} />
            </div>
          </div>

          <div className="mt-4">
            <div className="text-sm opacity-70 mb-2">{t('orCopy')}</div>
            <textarea className="textarea textarea-bordered w-full" rows={6} readOnly value={message} />

            <div className="mt-3 flex gap-2">
              <button type="button" className="btn" onClick={copyToClipboard}>
                {copied ? t('copied') : t('copy')}
              </button>
              <a className="btn btn-primary" href={whatsAppUrl} target="_blank" rel="noreferrer">
                {t('openWhatsapp')}
              </a>
            </div>
          </div>

          <div className="modal-action">
            <button type="button" className="btn" onClick={() => setQrOpen(false)}>
              {t('close')}
            </button>
          </div>
        </div>
        <form method="dialog" className="modal-backdrop">
          <button onClick={() => setQrOpen(false)}>{t('close')}</button>
        </form>
      </dialog>
    </>
  );
}
