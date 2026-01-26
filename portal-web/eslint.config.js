//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config'

export default [
  ...tanstackConfig,
  {
    ignores: ['eslint.config.js', 'prettier.config.js'],
  },
  {
    files: [
      'src/components/**/*.{ts,tsx}',
      'src/features/**/*.{ts,tsx}',
      'src/routes/**/*.{ts,tsx}',
      'src/lib/form/components/**/*.{ts,tsx}',
    ],
    rules: {
      'no-restricted-imports': [
        'error',
        {
          paths: [
            {
              name: 'ky',
              message:
                'Do not import ky in UI code. Use TanStack Query hooks from src/api/**.',
            },
            {
              name: '@/api/client',
              message:
                'Do not import the HTTP client in UI code. Use TanStack Query hooks from src/api/**.',
            },
            {
              name: '@/api/auth',
              importNames: ['authApi'],
              message:
                'Use TanStack Query mutation hooks from src/api/auth.ts (e.g. useForgotPasswordMutation, useGoogleAuthUrlMutation).',
            },
            {
              name: '@/api/business',
              importNames: ['businessApi'],
              message:
                'Use TanStack Query hooks from src/api/business.ts (e.g. useBusinessesQuery).',
            },
            {
              name: '@/api/customer',
              importNames: ['customerApi'],
              message:
                'Use TanStack Query hooks from src/api/customer.ts (e.g. useCustomersQuery/useCustomerQuery).',
            },
            {
              name: '@/api/order',
              importNames: ['orderApi'],
              message: 'Use TanStack Query hooks from src/api/order.ts.',
            },
            {
              name: '@/api/inventory',
              importNames: ['inventoryApi'],
              message: 'Use TanStack Query hooks from src/api/inventory.ts.',
            },
            {
              name: '@/api/user',
              importNames: ['userApi'],
              message: 'Use TanStack Query hooks from src/api/user.ts.',
            },
            {
              name: '@/api/onboarding',
              importNames: ['onboardingApi'],
              message: 'Use TanStack Query hooks from src/api/onboarding.ts.',
            },
            {
              name: '@/api/address',
              importNames: ['addressApi'],
              message: 'Use TanStack Query hooks from src/api/address.ts.',
            },
            {
              name: '@/api/metadata',
              importNames: ['metadataApi'],
              message: 'Use TanStack Query hooks from src/api/metadata.ts.',
            },
          ],
        },
      ],
      '@typescript-eslint/no-unnecessary-condition': 'off',
    },
  },
]
