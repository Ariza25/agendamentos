import { HashRouter, Route, Routes } from 'react-router-dom'

import { BookingPage } from './pages/BookingPage'
import { SlugEntryPage } from './pages/SlugEntryPage'

export default function App() {
  return (
    <HashRouter>
      <Routes>
        <Route element={<SlugEntryPage />} path="/" />
        <Route element={<BookingPage />} path="/booking/:slug" />
        <Route element={<SlugEntryPage />} path="*" />
      </Routes>
    </HashRouter>
  )
}
