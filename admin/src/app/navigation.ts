import type { LucideIcon } from 'lucide-react'
import {
  Building2,
  CalendarDays,
  CalendarClock,
  Settings,
  SlidersHorizontal,
  Users,
  UserCog,
  Scissors,
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
    label: 'Administração',
    paths: ['/configuration', '/team', '/tenants'],
  },
]
