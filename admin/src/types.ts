export type AppointmentStatus =
  | 'scheduled'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'no_show'

export type Appointment = {
  id: string
  tenantId: string
  professionalId: string
  serviceId: string
  resourceId?: string
  customerId: string
  startsAt: string
  endsAt: string
  status: AppointmentStatus
  notes: string
  createdBy: string
  seriesId?: string
  paymentStatus?: string
  managementCode?: string
  noShowFeeCents?: number
  source?: string
  cancellationReason?: string
  checkInAt?: string
  noShowAt?: string
  reminderSentAt?: string
  recurrenceRule?: string
  timezone?: string
  serviceIds?: string[]
  createdAt: string
  updatedAt: string
}

export type AppointmentAnalytics = {
  totalAppointments: number
  scheduled: number
  confirmed: number
  completed: number
  cancelled: number
  noShow: number
  noShowRate: number
  revenueCents: number
  byService: { serviceId: string; serviceName: string; count: number; revenueCents: number }[]
  byProfessional: { professionalId: string; professionalName: string; count: number }[]
  byDay: { date: string; count: number }[]
}

export type AppointmentSummary = {
  date: string
  total: number
  scheduled?: number
  confirmed?: number
  in_progress?: number
  completed?: number
  cancelled?: number
  no_show?: number
}

export type Customer = {
  id: string
  tenantId: string
  name: string
  phone: string
  email: string
  birthday: string
  notes: string
  createdAt: string
  updatedAt: string
}

export type Service = {
  id: string
  tenantId: string
  name: string
  description: string
  durationMinutes: number
  priceCents: number
  color: string
  active: boolean
  maxCapacity?: number
  createdAt: string
  updatedAt: string
}

export type Professional = {
  id: string
  tenantId: string
  userId: string
  name: string
  email: string
  phone: string
  color: string
  commissionPercentage?: number
  active: boolean
  maxCapacity?: number
  createdAt: string
  updatedAt: string
}

export type Resource = {
  id: string
  tenantId: string
  name: string
  type: string
  capacity: number
  active: boolean
  notes: string
  createdAt: string
  updatedAt: string
}

export type WorkingHours = {
  dayOfWeek: number
  openTime: string
  closeTime: string
  enabled: boolean
}

export type ClosedDay = {
  date: string
  name: string
  allDay: boolean
  openTime?: string
  closeTime?: string
}

export type Configuration = {
  tenantId: string
  establishmentName: string
  workingHours: WorkingHours[]
  appointmentIntervalMinutes: number
  defaultDurationMinutes: number
  allowOnlineBooking: boolean
  advanceNoticeHours: number
  bufferMinutes: number
  closedDays: ClosedDay[]
  notificationsEnabled: boolean
  notificationReminderHours: number
  notificationChannel: string
  depositEnabled: boolean
  depositType: string
  depositAmountCents: number
  depositPercentage: number
  noShowFeeCents: number
  cancellationPolicyHours: number
  calendarSyncEnabled: boolean
  logoUrl?: string
  primaryColor?: string
  address?: string
  addressNumber?: string
  neighborhood?: string
  city?: string
  state?: string
  postalCode?: string
  contactPhone?: string
  contactWhatsApp?: string
  instagram?: string
  facebook?: string
  website?: string
  description?: string
  updatedAt: string
}

export type NotificationLog = {
  id: string
  tenantId: string
  appointmentId?: string
  type: string
  channel: string
  recipient: string
  subject: string
  body: string
  status: string
  error?: string
  createdAt: string
  sentAt: string
}

export type Payment = {
  id: string
  tenantId: string
  appointmentId?: string
  amountCents: number
  method: string
  provider: string
  providerPaymentId?: string
  status: string
  description?: string
  paidAt?: string
  createdAt: string
}

export type Review = {
  id: string
  tenantId: string
  tenantSlug?: string
  appointmentId?: string
  customerId?: string
  customerName?: string
  status: string
  rating: number
  comment: string
  invitedAt: string
  submittedAt?: string
  createdAt: string
}

export type WaitingListEntry = {
  id: string
  tenantId: string
  serviceId?: string
  professionalId?: string
  customerId?: string
  name: string
  phone: string
  status: string
  createdAt: string
  notifiedAt?: string
}

export type Package = {
  id: string
  tenantId: string
  name: string
  description: string
  priceCents: number
  sessionsCount: number
  durationDays: number
  active: boolean
  createdAt: string
  updatedAt: string
}

export type PackageCredit = {
  id: string
  tenantId: string
  packageId?: string
  packageName?: string
  customerId: string
  customerName?: string
  sessionsTotal: number
  sessionsUsed: number
  expiresAt?: string
  createdAt: string
}

export type GiftCard = {
  id: string
  tenantId: string
  code: string
  valueCents: number
  balanceCents: number
  customerId?: string
  customerName?: string
  status: string
  expiresAt?: string
  createdAt: string
}

export type Membership = {
  id: string
  tenantId: string
  name: string
  description: string
  priceCents: number
  billingPeriodDays: number
  active: boolean
  createdAt: string
  updatedAt: string
}

export type MemberPlan = {
  id: string
  tenantId: string
  membershipId?: string
  membershipName?: string
  customerId: string
  customerName?: string
  status: string
  startedAt: string
  expiresAt?: string
  createdAt: string
}

export type AdminUser = {
  id: string
  name: string
  email: string
  phone: string
  role: string
  tenantId: string
  mfaEnabled: boolean
  createdAt: string
}

export type AdminInviteResponse = {
  user: AdminUser
  inviteUrl: string
  inviteEmailSent: boolean
  temporaryPassword?: string
}

export type Tenant = {
  id: string
  name: string
  slug: string
  active: boolean
  createdAt: string
  updatedAt: string
}

// --- ERP Inventory & Supplies ---
export type Product = {
  id: string
  tenantId: string
  name: string
  sku: string
  category: string
  unit: string
  costPriceCents: number
  salePriceCents: number
  stockQuantity: number
  minStockAlert: number
  isSupply: boolean
  active: boolean
  createdAt: string
  updatedAt: string
}

export type StockMovementType = 'in' | 'out' | 'adjustment'

export type StockMovement = {
  id: string
  tenantId: string
  productId: string
  productName?: string
  movementType: StockMovementType
  quantity: number
  previousStock: number
  newStock: number
  reason: string
  referenceType?: string
  referenceId?: string
  performedBy?: string
  createdAt: string
}

export type ServiceSupply = {
  id: string
  tenantId: string
  serviceId: string
  productId: string
  quantity: number
}

// --- ERP Commissions ---
export type CommissionType = 'percentage' | 'fixed'

export type CommissionRule = {
  id: string
  tenantId: string
  professionalId?: string
  serviceId?: string
  commissionType: CommissionType
  value: number
  active: boolean
  createdAt: string
}

export type CommissionRecordStatus = 'pending' | 'settled'

export type CommissionRecord = {
  id: string
  tenantId: string
  appointmentId: string
  professionalId: string
  amountCents: number
  status: CommissionRecordStatus
  settledAt?: string
  createdAt: string
}

export type CommissionSummary = {
  professionalId: string
  professionalName: string
  pendingCents: number
  settledCents: number
  totalCents: number
  count: number
}

// --- ERP Finance ---
export type TransactionType = 'income' | 'expense'
export type TransactionStatus = 'pending' | 'paid' | 'cancelled'

export type FinancialTransaction = {
  id: string
  tenantId: string
  type: TransactionType
  category: string
  description: string
  amountCents: number
  dueDate: string
  paidAt?: string
  status: TransactionStatus
  referenceType?: string
  referenceId?: string
  paymentMethod?: string
  createdAt: string
  updatedAt: string
}

export type CategoryTotal = {
  category: string
  amountCents: number
}

export type FinancialSummary = {
  totalIncomeCents: number
  totalExpenseCents: number
  netProfitCents: number
  pendingIncomeCents: number
  pendingExpenseCents: number
  incomeByCategory: CategoryTotal[]
  expenseByCategory: CategoryTotal[]
}

// --- Calendar Sync ---
export type CalendarConnection = {
  id: string
  tenantId: string
  professionalId: string
  provider: 'google' | 'caldav'
  externalCalendarId: string
  syncToken?: string
  active: boolean
  lastSyncAt?: string
  createdAt: string
}

// --- Cash Register & POS ---
export type CashRegisterStatus = 'open' | 'closed'
export type CashRegisterEntryType = 'sangria' | 'suprimento'

export type CashRegister = {
  id: string
  tenantId: string
  openedByUserId: string
  closedByUserId?: string
  status: CashRegisterStatus
  initialBalanceCents: number
  closingBalanceCents: number
  systemBalanceCents: number
  differenceCents: number
  openedAt: string
  closedAt?: string
  notes?: string
}

export type CashRegisterEntry = {
  id: string
  tenantId: string
  cashRegisterId: string
  type: CashRegisterEntryType
  amountCents: number
  reason: string
  createdByUserId: string
  createdAt: string
}

export type ShiftSummary = {
  cashRegister: CashRegister
  entries: CashRegisterEntry[]
  totalSangriasCents: number
  totalSuprimentosCents: number
  totalCashIncomeCents: number
  totalCardIncomeCents: number
  totalPixIncomeCents: number
  totalOtherIncomeCents: number
  currentDrawerCents: number
}

export type CurrentCashRegisterResponse = {
  status: 'open' | 'closed'
  open: boolean
  summary?: ShiftSummary
}

// --- Checkout ---
export type CheckoutProductItem = {
  productId: string
  quantity: number
}

export type CheckoutInput = {
  appointmentId: string
  products?: CheckoutProductItem[]
  discountCents?: number
  surchargeCents?: number
  paymentMethod: 'pix' | 'cash' | 'credit_card' | 'debit_card' | 'other'
  notes?: string
}

export type CheckoutResult = {
  appointment: Appointment
  totalChargedCents: number
  serviceTotalCents: number
  productsTotalCents: number
  commissionCents: number
}
