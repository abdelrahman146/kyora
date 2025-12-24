import { Navigate, Route, Routes } from 'react-router-dom'
import { StorefrontPage } from './pages/StorefrontPage'

function IndexPage() {
  return (
    <div className="min-h-dvh flex items-center justify-center p-6">
      <div className="text-center space-y-2">
        <div className="text-xl font-bold">Kyora Storefront</div>
        <div className="opacity-70">Open a storefront link to continue.</div>
      </div>
    </div>
  )
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<IndexPage />} />
      <Route path="/:storefrontPublicId" element={<StorefrontPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
