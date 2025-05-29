// src/App.tsx
// import React from "react"
import { Card } from "antd"
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
  LabelList,
} from "recharts"

// Sample data
const data = [
  { name: "CPU", usage: 20},
  { name: "MEMORY", usage: 35},
  { name: "DISK", usage: 57 },
]

export function PerfStat() {
  return (
    <div >
      <Card title="Performace Stats" bordered={false} className="w-[100%] border-1">
        <ResponsiveContainer className="min-h-[220px] w-[100%]">
          <BarChart 
            accessibilityLayer
            layout="vertical" 
            data={data} 
            // margin={{ top: 20, right: 30, left: 0, bottom: 5 }}
             margin={{
              right: 16,
            }}
        >
            {/* <CartesianGrid strokeDasharray="3 3" /> */}
            <CartesianGrid horizontal={false} />
            {/* <XAxis dataKey="name" /> */}
            <XAxis dataKey="usage" type="number" hide />
            {/* <YAxis /> */}
            <YAxis
              dataKey="name"
              type="category"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={(value) => value.slice(0, 3)}
              hide
            />
            <Tooltip />
            {/* <Bar dataKey="sales" fill="#1890ff" barSize={40} /> */}

            <Bar
              dataKey="usage"
              layout="vertical"
            //   fill="var(--color-usage)"
              barSize={40}
              radius={4}
            >
              <LabelList
                dataKey="name"
                position="insideLeft"
                offset={8}
                // className="fill-[--color-label]"
                fontSize={12}
              />
              <LabelList
                dataKey="usage"
                position="right"
                offset={8}
                // className="fill-foreground"
                fontSize={12}
              />
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  )
}

// export default App
