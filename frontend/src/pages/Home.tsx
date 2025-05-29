import Navbar from "@/components/navbar"
import Footer from "@/components/footer"
import { PerfStat } from "@/components/perfStat"
import { PollingStat } from "@/components/pollingStats"
import { PollTypeStat } from "@/components/pollTypesStat"
import { PollPeriodStat } from "@/components/pollPeriodStat"


export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col px-0 md:px-[6%]">
        <Navbar />

        <section className="grid md:grid-cols-3 gap-6  py-2 pt-8 min-h-[230px]">
            <PerfStat />
            <PollingStat />
            <PollTypeStat />
            
        </section>
        <section className="grid md:grid-cols-1 gap-6 py-2 pb-4">
            <PollPeriodStat />
            
        </section>
        

        <Footer/>
    </div>
  )
}
