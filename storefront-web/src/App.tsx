import { useEffect } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { StorefrontPage } from './pages/StorefrontPage';

function IndexPage() {
  return (
    <div className="min-h-dvh flex items-center justify-center p-6 bg-base-200">
      <div className="text-center space-y-2">
        <div className="text-xl font-bold text-base-content">Kyora Storefront</div>
        <div className="text-base-content/70">Open a storefront link to continue.</div>
      </div>
    </div>
  );
}

export default function App() {
  const { i18n } = useTranslation();

  // Set HTML dir attribute based on language
  useEffect(() => {
    const dir = i18n.language === 'ar' ? 'rtl' : 'ltr';
    document.documentElement.dir = dir;
    document.documentElement.lang = i18n.language;
  }, [i18n.language]);

  return (
    <Routes>
      <Route path="/" element={<IndexPage />} />
      <Route path="/:storefrontPublicId" element={<StorefrontPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
