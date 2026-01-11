import { z } from 'zod'

export const CustomersSearchSchema = z.object({
  search: z.string().optional(),
  page: z.number().optional().default(1),
  pageSize: z.number().optional().default(20),
  sortBy: z.string().optional(),
  sortOrder: z.enum(['asc', 'desc']).optional().default('desc'),
  countryCode: z.string().optional(),
  hasOrders: z.boolean().optional(),
  socialPlatforms: z
    .array(
      z.enum(['instagram', 'tiktok', 'facebook', 'x', 'snapchat', 'whatsapp']),
    )
    .optional(),
})

export type CustomersSearch = z.infer<typeof CustomersSearchSchema>
