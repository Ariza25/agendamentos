import { HashRouter, Route, Routes } from 'react-router-dom'

import { BookingPage } from './pages/BookingPage'
import { ManagePage } from './pages/ManagePage'
import { SlugEntryPage } from './pages/SlugEntryPage'
import { WaitingListPage } from './pages/WaitingListPage'

export default function App() {
  return (
    <HashRouter>
      <Routes>
        <Route element={<SlugEntryPage />} path="/" />
        <Route element={<BookingPage />} path="/booking/:slug" />
        <Route element={<ManagePage />} path="/manage/:slug" />
        <Route element={<WaitingListPage />} path="/waiting-list/:slug" />
        <Route element={<SlugEntryPage />} path="*" />
      </Routes>
    </HashRouter>
  )
}
