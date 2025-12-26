import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { MessageCircle } from 'lucide-react';

interface FloatingWhatsAppButtonProps {
  phoneNumber: string;
  businessName: string;
  defaultMessage?: string;
}

/**
 * FloatingWhatsAppButton Molecule - WhatsApp FAB for conversions
 * Memoized to prevent unnecessary re-renders
 * Optimized with useMemo for WhatsApp URL
 * Follows KDS 5.4 with high visibility and proper RTL support
 */
export const FloatingWhatsAppButton = memo<FloatingWhatsAppButtonProps>(function FloatingWhatsAppButton({
  phoneNumber,
  businessName,
  defaultMessage,
}) {
  const { t } = useTranslation();

  // Memoize WhatsApp URL to prevent recalculation
  const whatsappUrl = useMemo(() => {
    if (!phoneNumber) return '';
    const message = defaultMessage || `Hi ${businessName}, I am interested in...`;
    return `https://wa.me/${phoneNumber}?text=${encodeURIComponent(message)}`;
  }, [phoneNumber, businessName, defaultMessage]);

  if (!phoneNumber) return null;

  return (
    <a
      href={whatsappUrl}
      target="_blank"
      rel="noopener noreferrer"
      className="fixed bottom-20 end-4 z-30 w-14 h-14 rounded-full bg-[#25D366] text-white shadow-float flex items-center justify-center hover:scale-110 active:scale-95 transition-transform focus-ring"
      aria-label={t('sendOrderOnWhatsapp')}
    >
      <MessageCircle className="w-7 h-7" strokeWidth={2} fill="currentColor" />
    </a>
  );
});
