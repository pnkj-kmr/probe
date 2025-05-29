// src/App.tsx
// import React from "react"
import { Card } from "antd"
import { PolarAngleAxis, PolarGrid, Radar, RadarChart, ResponsiveContainer } from "recharts"


// const chartData = [
//   { name: "safari", count: 100000, fill: "var(--color-safari)" },
// ]

// Sample data
const data = [
     { name: "ICMP", total: 186 },
    { name: "SNMP", total: 305 },
    { name: "SSH", total: 237 },
    { name: "SDWAN", total: 273 },
    { name: "HTTP", total: 273 },
]

export function PollTypeStat() {
  return (
    <div >
      <Card title="Poll Types Stat" bordered={false} className="w-[100%] border-1">
        <ResponsiveContainer className="min-h-[220px] w-[100%]">
          
          <RadarChart data={data}>
            {/* <ChartTooltip cursor={false} content={<ChartTooltipContent />} /> */}
            <PolarAngleAxis dataKey="name" />
            <PolarGrid />
            <Radar
              dataKey="total"
            //   fill="var(--color-desktop)"
                // fill="blue"
              fillOpacity={0.6}
              dot={{
                r: 4,
                fillOpacity: 1,
              }}
            />
          </RadarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  )
}

// export default App
