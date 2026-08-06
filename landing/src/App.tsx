import { CtaSection } from './components/CtaSection'
import { FeaturesSection } from './components/FeaturesSection'
import { Footer } from './components/Footer'
import { Header } from './components/Header'
import { HeroSection } from './components/HeroSection'
import { HowItWorksSection } from './components/HowItWorksSection'

export function App() {
  return (
    <div className="relative min-h-screen text-ink selection:bg-primary selection:text-white">
      <Header />
      <main>
        <HeroSection />
        <FeaturesSection />
        <HowItWorksSection />
        <CtaSection />
      </main>
      <Footer />
    </div>
  )
}
