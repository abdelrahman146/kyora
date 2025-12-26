import { useState } from 'react';
import { Search, Mail, Eye, EyeOff } from 'lucide-react';
import { Button, Input, Badge, Skeleton, SkeletonText } from '@/components/atoms';

export default function DesignSystem() {
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleLoadingDemo = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="container mx-auto px-4 max-w-6xl">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-base-content mb-2">
            Kyora Design System
          </h1>
          <p className="text-lg text-neutral-500">
            Portal Web App - Atomic Components
          </p>
        </div>

        {/* Color Palette */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">Color Palette</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <ColorSwatch color="primary" label="Primary (Teal)" />
            <ColorSwatch color="secondary" label="Secondary (Gold)" />
            <ColorSwatch color="accent" label="Accent" />
            <ColorSwatch color="neutral" label="Neutral" />
            <ColorSwatch color="success" label="Success" />
            <ColorSwatch color="warning" label="Warning" />
            <ColorSwatch color="error" label="Error" />
            <ColorSwatch color="info" label="Info" />
          </div>
        </section>

        {/* Buttons */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">Buttons</h2>
          <div className="bg-base-100 rounded-lg p-6 space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-4">Variants</h3>
              <div className="flex flex-wrap gap-4">
                <Button variant="primary">Primary Button</Button>
                <Button variant="secondary">Secondary Button</Button>
                <Button variant="ghost">Ghost Button</Button>
                <Button variant="outline">Outline Button</Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">Sizes</h3>
              <div className="flex flex-wrap items-center gap-4">
                <Button size="sm">Small</Button>
                <Button size="md">Medium</Button>
                <Button size="lg">Large</Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">States</h3>
              <div className="flex flex-wrap gap-4">
                <Button loading={loading} onClick={handleLoadingDemo}>
                  {loading ? 'Loading...' : 'Click to Load'}
                </Button>
                <Button disabled>Disabled</Button>
                <Button fullWidth>Full Width Button</Button>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">RTL Test (with Icons)</h3>
              <div className="flex flex-wrap gap-4" dir="rtl">
                <Button>زر أساسي</Button>
                <Button variant="secondary">زر ثانوي</Button>
              </div>
            </div>
          </div>
        </section>

        {/* Inputs */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">Input Fields</h2>
          <div className="bg-base-100 rounded-lg p-6 space-y-6">
            <div className="grid md:grid-cols-2 gap-6">
              <Input
                label="Email Address"
                type="email"
                placeholder="Enter your email"
                startIcon={<Mail size={18} />}
              />
              <Input
                label="Password"
                type={showPassword ? 'text' : 'password'}
                placeholder="Enter your password"
                endIcon={
                  <button
                    type="button"
                    onClick={() => {
                      setShowPassword(!showPassword);
                    }}
                    className="cursor-pointer"
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                }
              />
            </div>

            <Input
              label="Search"
              placeholder="Search products..."
              startIcon={<Search size={18} />}
              helperText="Enter at least 3 characters"
            />

            <Input
              label="Error State"
              placeholder="This field has an error"
              error="This field is required"
            />

            <div dir="rtl">
              <Input
                label="البحث (RTL)"
                placeholder="ابحث عن المنتجات..."
                startIcon={<Search size={18} />}
              />
            </div>
          </div>
        </section>

        {/* Badges */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">Badges</h2>
          <div className="bg-base-100 rounded-lg p-6 space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-4">Variants</h3>
              <div className="flex flex-wrap gap-3">
                <Badge variant="default">Default</Badge>
                <Badge variant="primary">Primary</Badge>
                <Badge variant="secondary">Secondary</Badge>
                <Badge variant="success">Success</Badge>
                <Badge variant="warning">Warning</Badge>
                <Badge variant="error">Error</Badge>
                <Badge variant="info">Info</Badge>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">Sizes</h3>
              <div className="flex flex-wrap items-center gap-3">
                <Badge size="sm">Small</Badge>
                <Badge size="md">Medium</Badge>
                <Badge size="lg">Large</Badge>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">Use Cases</h3>
              <div className="flex flex-wrap gap-3">
                <Badge variant="success">Paid</Badge>
                <Badge variant="warning">Pending</Badge>
                <Badge variant="error">Out of Stock</Badge>
                <Badge variant="info">New</Badge>
              </div>
            </div>
          </div>
        </section>

        {/* Skeletons */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">Loading Skeletons</h2>
          <div className="bg-base-100 rounded-lg p-6 space-y-6">
            <div>
              <h3 className="text-lg font-semibold mb-4">Text Skeleton</h3>
              <SkeletonText lines={3} />
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">Variants</h3>
              <div className="space-y-4">
                <Skeleton variant="rectangular" height={100} />
                <div className="flex gap-4">
                  <Skeleton variant="circular" width={40} height={40} />
                  <div className="flex-1 space-y-2">
                    <Skeleton variant="text" width="60%" />
                    <Skeleton variant="text" width="40%" />
                  </div>
                </div>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-semibold mb-4">Card Loading State</h3>
              <div className="border border-base-300 rounded-lg p-4">
                <div className="flex gap-4">
                  <Skeleton variant="rectangular" width={80} height={80} />
                  <div className="flex-1 space-y-3">
                    <Skeleton variant="text" />
                    <Skeleton variant="text" width="70%" />
                    <div className="flex gap-2 mt-2">
                      <Skeleton variant="rectangular" width={60} height={24} />
                      <Skeleton variant="rectangular" width={60} height={24} />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* RTL/LTR Toggle Demo */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">RTL Support Test</h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="bg-base-100 rounded-lg p-6" dir="ltr">
              <h3 className="text-lg font-semibold mb-4">LTR (English)</h3>
              <div className="space-y-4">
                <Button fullWidth>Click Here</Button>
                <Input placeholder="Enter text" startIcon={<Search size={18} />} />
                <div className="flex gap-2">
                  <Badge variant="success">Active</Badge>
                  <Badge variant="warning">Pending</Badge>
                </div>
              </div>
            </div>

            <div className="bg-base-100 rounded-lg p-6" dir="rtl">
              <h3 className="text-lg font-semibold mb-4">RTL (العربية)</h3>
              <div className="space-y-4">
                <Button fullWidth>انقر هنا</Button>
                <Input placeholder="أدخل النص" startIcon={<Search size={18} />} />
                <div className="flex gap-2">
                  <Badge variant="success">نشط</Badge>
                  <Badge variant="warning">قيد الانتظار</Badge>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}

function ColorSwatch({ color, label }: { color: string; label: string }) {
  const colorClasses: Record<string, string> = {
    primary: 'bg-primary',
    secondary: 'bg-secondary',
    accent: 'bg-accent',
    neutral: 'bg-neutral',
    success: 'bg-success',
    warning: 'bg-warning',
    error: 'bg-error',
    info: 'bg-info',
  };

  return (
    <div className="space-y-2">
      <div className={`h-20 rounded-lg border border-base-300 ${colorClasses[color] || ''}`} />
      <p className="text-sm font-medium text-center">{label}</p>
    </div>
  );
}
