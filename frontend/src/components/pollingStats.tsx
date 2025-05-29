// src/App.tsx
// import React from "react"
import { Card } from "antd"
import {
  Label,
  PolarGrid,
  PolarRadiusAxis,
  RadialBar,
  RadialBarChart,
  ResponsiveContainer,
} from "recharts"


// const chartData = [
//   { name: "safari", count: 100000, fill: "var(--color-safari)" },
// ]

// Sample data
const data = [
    {name: "Resources", count: 1000000}
]

export function PollingStat() {
  return (
    <div >
      <Card title="Performace Stats" bordered={false} className="w-[100%] border-1">
        <ResponsiveContainer className="min-h-[220px] w-[100%]">
          <RadialBarChart
            data={data}
            endAngle={320}
            startAngle={0}
            innerRadius={100}
            outerRadius={140}
          >
            <PolarGrid
              gridType="circle"
              radialLines={false}
              stroke="none"
              className="first:fill-muted last:fill-background"
              polarRadius={[86, 74]}
            />
            <RadialBar dataKey="count" background />
            <PolarRadiusAxis tick={false} tickLine={false} axisLine={false}>
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="fill-foreground text-4xl font-bold"
                        >
                          {data[0].count.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground"
                        >
                          Resources
                        </tspan>
                      </text>
                    )
                  }
                }}
              />
            </PolarRadiusAxis>
          </RadialBarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  )
}

// export default App
