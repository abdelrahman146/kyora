import { z } from 'zod'

export const OrdersSearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional(),
  pageSize: z.number().optional(),
  sortBy: z
    .enum(['orderNumber', 'total', 'status', 'paymentStatus', 'orderedAt'])
    .optional(),
  sortOrder: z.enum(['asc', 'desc']).optional(),
  status: z.array(z.string()).optional(),
  paymentStatus: z.array(z.string()).optional(),
  socialPlatforms: z
    .array(
      z.enum(['instagram', 'tiktok', 'facebook', 'x', 'snapchat', 'whatsapp']),
    )
    .optional(),
  customerId: z.string().optional(),
  from: z.string().optional(),
  to: z.string().optional(),
})

export type OrdersSearch = z.infer<typeof OrdersSearchSchema>
