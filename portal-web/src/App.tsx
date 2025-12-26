import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import DesignSystem from './routes/design-system';

function Home() {
  const { t, i18n } = useTranslation();

  const toggleLanguage = () => {
    const newLang = i18n.language === 'ar' ? 'en' : 'ar';
    i18n.changeLanguage(newLang);
  };

  return (
    <div className="min-h-screen bg-base-100">
      <div className="container mx-auto p-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-4xl font-bold text-primary">
            Kyora Portal
          </h1>
          <div className="flex gap-4">
            <Link to="/design-system" className="btn btn-secondary">
              Design System
            </Link>
            <button 
              onClick={toggleLanguage}
              className="btn btn-primary"
            >
              {i18n.language === 'ar' ? 'English' : 'عربي'}
            </button>
          </div>
        </div>
        
        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('dashboard.welcome')}</h2>
            <p>{t('common.loading')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary">{t('common.save')}</button>
              <button className="btn btn-ghost">{t('common.cancel')}</button>
            </div>
          </div>
        </div>

        <div className="alert alert-info mt-8">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <span>Portal web app initialized successfully! Ready for development.</span>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/design-system" element={<DesignSystem />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
