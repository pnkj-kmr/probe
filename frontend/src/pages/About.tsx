import Navbar from "@/components/navbar"
import Footer from "@/components/footer"
import AboutWidget from "@/components/AboutWidget"


export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col px-0 md:px-[6%]">
        <Navbar />

        <section className="grid md:grid-cols-1 gap-6 py-2 pt-8">
          <AboutWidget />
        </section>
        
        <Footer />
    </div>
  )
}
