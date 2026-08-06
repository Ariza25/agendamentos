import type { LucideIcon } from 'lucide-react'
import {
  Building2,
  CalendarDays,
  CalendarClock,
  CreditCard,
  Gift,
  History,
  LineChart,
  Package,
  RefreshCcw,
  Settings,
  SlidersHorizontal,
  Users,
  UserCog,
  Scissors,
  Clock3,
} from 'lucide-react'

import type { AdminAreaKey } from '../lib/adminPermissions'

export type AdminNavigationItem = {
  to: string
  label: string
  icon: LucideIcon
  area: AdminAreaKey
  roles?: string[]
  masterOnly?: boolean
}

export const navigationItems: AdminNavigationItem[] = [
  { to: '/agenda', label: 'Agenda', icon: CalendarDays, area: 'agenda' },
  { to: '/appointments', label: 'Agendamentos', icon: CalendarClock, area: 'appointments' },
  { to: '/customers', label: 'Clientes', icon: Users, area: 'customers' },
  { to: '/services', label: 'Serviços', icon: Scissors, area: 'services' },
  { to: '/professionals', label: 'Profissionais', icon: UserCog, area: 'professionals' },
  { to: '/analytics', label: 'Dashboard', icon: LineChart, area: 'analytics' },
  { to: '/payments', label: 'Pagamentos', icon: CreditCard, area: 'payments' },
  { to: '/packages', label: 'Pacotes', icon: Package, area: 'packages' },
  { to: '/gift-cards', label: 'Gift cards', icon: Gift, area: 'gift-cards' },
  { to: '/memberships', label: 'Assinaturas', icon: RefreshCcw, area: 'memberships' },
  { to: '/waiting-list', label: 'Lista de espera', icon: Clock3, area: 'waiting-list' },
  { to: '/notifications', label: 'Notificações', icon: History, area: 'notifications' },
  { to: '/configuration', label: 'Configuração', icon: SlidersHorizontal, area: 'configuration' },
  { to: '/team', label: 'Equipe', icon: Settings, area: 'team' },
  { to: '/tenants', label: 'Estabelecimentos', icon: Building2, area: 'tenants', masterOnly: true },
]

export const navigationGroups = [
  {
    label: 'Operação',
    paths: ['/agenda', '/appointments'],
  },
  {
    label: 'Cadastro',
    paths: ['/customers', '/services', '/professionals'],
  },
  {
    label: 'Gestão',
    paths: ['/analytics', '/payments', '/packages', '/gift-cards', '/memberships', '/waiting-list', '/notifications'],
  },
  {
    label: 'Administração',
    paths: ['/configuration', '/team', '/tenants'],
  },
]
