import Navbar from "@/components/navbar"
import Footer from "@/components/footer"


export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col px-0 md:px-[6%]">
        <Navbar />

        <div>About agent version and agent ID if any</div>
        
        <Footer />
    </div>
  )
}
