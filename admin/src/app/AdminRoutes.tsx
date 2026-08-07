import { lazy, Suspense } from 'react'
import { Navigate, Route, Routes } from 'react-router-dom'

import type { AdminAreaKey } from '../lib/adminPermissions'

const AgendaPage = lazy(() => import('../pages/AgendaPage').then((module) => ({ default: module.AgendaPage })))
const AppointmentsPage = lazy(() => import('../pages/AppointmentsPage').then((module) => ({ default: module.AppointmentsPage })))
const CustomersPage = lazy(() => import('../pages/CustomersPage').then((module) => ({ default: module.CustomersPage })))
const ServicesPage = lazy(() => import('../pages/ServicesPage').then((module) => ({ default: module.ServicesPage })))
const ResourcesPage = lazy(() => import('../pages/ResourcesPage').then((module) => ({ default: module.ResourcesPage })))
const ProfessionalsPage = lazy(() => import('../pages/ProfessionalsPage').then((module) => ({ default: module.ProfessionalsPage })))
const ConfigurationPage = lazy(() => import('../pages/ConfigurationPage').then((module) => ({ default: module.ConfigurationPage })))
const TeamPage = lazy(() => import('../pages/TeamPage').then((module) => ({ default: module.TeamPage })))
const TenantsPage = lazy(() => import('../pages/TenantsPage').then((module) => ({ default: module.TenantsPage })))
const AnalyticsPage = lazy(() => import('../pages/AnalyticsPage').then((module) => ({ default: module.AnalyticsPage })))
const NotificationsPage = lazy(() => import('../pages/NotificationsPage').then((module) => ({ default: module.NotificationsPage })))
const PaymentsPage = lazy(() => import('../pages/PaymentsPage').then((module) => ({ default: module.PaymentsPage })))
const PackagesPage = lazy(() => import('../pages/PackagesPage').then((module) => ({ default: module.PackagesPage })))
const GiftCardsPage = lazy(() => import('../pages/GiftCardsPage').then((module) => ({ default: module.GiftCardsPage })))
const MembershipsPage = lazy(() => import('../pages/MembershipsPage').then((module) => ({ default: module.MembershipsPage })))
const WaitingListPage = lazy(() => import('../pages/WaitingListPage').then((module) => ({ default: module.WaitingListPage })))
const ReviewsPage = lazy(() => import('../pages/ReviewsPage').then((module) => ({ default: module.ReviewsPage })))
const WidgetPage = lazy(() => import('../pages/WidgetPage').then((module) => ({ default: module.WidgetPage })))

type AdminRoutesProps = {
  canAccessRoute: (area: AdminAreaKey) => boolean
  defaultRoute: string
  role?: string
}

export function AdminRoutes({ canAccessRoute, defaultRoute }: AdminRoutesProps) {
  return (
    <Suspense fallback={<AdminRouteFallback />}>
      <Routes>
        <Route element={<Navigate replace to={defaultRoute} />} path="/" />
        <Route element={canAccessRoute('agenda') ? <AgendaPage /> : <AccessDeniedPage />} path="/agenda" />
        <Route element={canAccessRoute('appointments') ? <AppointmentsPage /> : <AccessDeniedPage />} path="/appointments" />
        <Route element={canAccessRoute('customers') ? <CustomersPage /> : <AccessDeniedPage />} path="/customers" />
        <Route element={canAccessRoute('services') ? <ServicesPage /> : <AccessDeniedPage />} path="/services" />
        <Route element={canAccessRoute('resources') ? <ResourcesPage /> : <AccessDeniedPage />} path="/resources" />
        <Route element={canAccessRoute('professionals') ? <ProfessionalsPage /> : <AccessDeniedPage />} path="/professionals" />
        <Route element={canAccessRoute('configuration') ? <ConfigurationPage /> : <AccessDeniedPage />} path="/configuration" />
        <Route element={canAccessRoute('team') ? <TeamPage /> : <AccessDeniedPage />} path="/team" />
        <Route element={canAccessRoute('tenants') ? <TenantsPage /> : <AccessDeniedPage />} path="/tenants" />
        <Route element={canAccessRoute('analytics') ? <AnalyticsPage /> : <AccessDeniedPage />} path="/analytics" />
        <Route element={canAccessRoute('notifications') ? <NotificationsPage /> : <AccessDeniedPage />} path="/notifications" />
        <Route element={canAccessRoute('payments') ? <PaymentsPage /> : <AccessDeniedPage />} path="/payments" />
        <Route element={canAccessRoute('packages') ? <PackagesPage /> : <AccessDeniedPage />} path="/packages" />
        <Route element={canAccessRoute('gift-cards') ? <GiftCardsPage /> : <AccessDeniedPage />} path="/gift-cards" />
        <Route element={canAccessRoute('memberships') ? <MembershipsPage /> : <AccessDeniedPage />} path="/memberships" />
        <Route element={canAccessRoute('waiting-list') ? <WaitingListPage /> : <AccessDeniedPage />} path="/waiting-list" />
        <Route element={canAccessRoute('reviews') ? <ReviewsPage /> : <AccessDeniedPage />} path="/reviews" />
        <Route element={canAccessRoute('widget') ? <WidgetPage /> : <AccessDeniedPage />} path="/widget" />
        <Route element={<Navigate replace to={defaultRoute} />} path="*" />
      </Routes>
    </Suspense>
  )
}

function AccessDeniedPage() {
  return (
    <section className="rounded-[1.75rem] border border-red-200 bg-red-50 p-6">
      <h2 className="font-headline text-2xl font-extrabold text-red-700">Acesso negado</h2>
      <p className="mt-2 text-sm text-red-700">
        Seu perfil não tem permissão para acessar esta área do painel.
      </p>
    </section>
  )
}

function AdminRouteFallback() {
  return (
    <section className="rounded-[1.75rem] border border-stone-200 bg-white px-6 py-5 text-sm font-bold text-stone-700 shadow-sm">
      Carregando tela...
    </section>
  )
}
