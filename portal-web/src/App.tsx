import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { RequireAuth } from './components/routing/RequireAuth';
import DesignSystem from './routes/design-system';
import LoginPage from './routes/login';
import DashboardPage from './routes/dashboard';
import OAuthCallbackPage from './routes/oauth-callback';
import ForgotPasswordPage from './routes/forgot-password';
import ResetPasswordPage from './routes/reset-password';
import OnboardingRoot from './routes/onboarding';
import PlanSelectionPage from './routes/onboarding/plan';
import EmailEntryPage from './routes/onboarding/email';
import VerifyEmailPage from './routes/onboarding/verify';
import BusinessSetupPage from './routes/onboarding/business';
import PaymentPage from './routes/onboarding/payment';
import CompletePage from './routes/onboarding/complete';
import OnboardingOAuthCallbackPage from './routes/onboarding/oauth-callback';
import CustomersPage from './routes/dashboard/customers';
import HomePage from './routes/home';

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public Home (Unauthenticated) */}
          <Route path="/welcome" element={<DesignSystem />} />
          
          {/* Authenticated Home */}
          <Route
            path="/"
            element={
              <RequireAuth>
                <HomePage />
              </RequireAuth>
            }
          />

          {/* Auth Routes */}
          <Route path="/auth/login" element={<LoginPage />} />
          <Route path="/auth/forgot-password" element={<ForgotPasswordPage />} />
          <Route path="/auth/reset-password" element={<ResetPasswordPage />} />
          <Route path="/auth/oauth/callback" element={<OAuthCallbackPage />} />
          
          {/* Legacy redirects for backward compatibility */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/forgot-password" element={<ForgotPasswordPage />} />
          <Route path="/reset-password" element={<ResetPasswordPage />} />
          <Route path="/oauth/callback" element={<OAuthCallbackPage />} />
          
          {/* Onboarding Routes */}
          <Route path="/onboarding" element={<OnboardingRoot />}>
            <Route index element={<PlanSelectionPage />} />
            <Route path="plan" element={<PlanSelectionPage />} />
            <Route path="email" element={<EmailEntryPage />} />
            <Route path="verify" element={<VerifyEmailPage />} />
            <Route path="business" element={<BusinessSetupPage />} />
            <Route path="payment" element={<PaymentPage />} />
            <Route path="complete" element={<CompletePage />} />
            <Route path="oauth-callback" element={<OnboardingOAuthCallbackPage />} />
          </Route>
          
          {/* Design System (Dev Only) */}
          <Route path="/design-system" element={<DesignSystem />} />

          {/* Business Routes */}
          <Route path="/businesses/:businessDescriptor">
            <Route
              path="dashboard"
              element={
                <RequireAuth>
                  <DashboardPage />
                </RequireAuth>
              }
            />
            <Route
              path="customers"
              element={
                <RequireAuth>
                  <CustomersPage />
                </RequireAuth>
              }
            />
          </Route>

          {/* Legacy dashboard redirect */}
          <Route path="/dashboard" element={<HomePage />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
