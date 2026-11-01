import type { LucideIcon } from 'lucide-react'
import {
  Building2,
  CalendarDays,
  CalendarClock,
  Code2,
  CreditCard,
  DollarSign,
  Boxes,
  Percent,
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
  Star,
  DoorOpen,
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
  { to: '/resources', label: 'Recursos e salas', icon: DoorOpen, area: 'resources' },
  { to: '/professionals', label: 'Profissionais', icon: UserCog, area: 'professionals' },
  { to: '/analytics', label: 'Dashboard', icon: LineChart, area: 'analytics' },
  { to: '/finance', label: 'Financeiro', icon: DollarSign, area: 'finance' },
  { to: '/commissions', label: 'Comissões', icon: Percent, area: 'commissions' },
  { to: '/inventory', label: 'Estoque & Produtos', icon: Boxes, area: 'inventory' },
  { to: '/payments', label: 'Pagamentos', icon: CreditCard, area: 'payments' },
  { to: '/packages', label: 'Pacotes', icon: Package, area: 'packages' },
  { to: '/gift-cards', label: 'Gift cards', icon: Gift, area: 'gift-cards' },
  { to: '/memberships', label: 'Assinaturas', icon: RefreshCcw, area: 'memberships' },
  { to: '/waiting-list', label: 'Lista de espera', icon: Clock3, area: 'waiting-list' },
  { to: '/reviews', label: 'Avaliações', icon: Star, area: 'reviews' },
  { to: '/notifications', label: 'Notificações', icon: History, area: 'notifications' },
  { to: '/widget', label: 'Widget', icon: Code2, area: 'widget' },
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
    paths: ['/customers', '/services', '/resources', '/professionals'],
  },
  {
    label: 'Gestão',
    paths: ['/analytics', '/finance', '/commissions', '/inventory', '/payments', '/packages', '/gift-cards', '/memberships', '/waiting-list', '/reviews', '/notifications'],
  },
  {
    label: 'Administração',
    paths: ['/configuration', '/widget', '/team', '/tenants'],
  },
]
