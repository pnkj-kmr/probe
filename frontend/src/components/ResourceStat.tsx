import { useEffect, useState } from "react";
import { getResourceData } from "@/services/api";
import { Card } from "antd";
import {
  Label,
  PolarGrid,
  PolarRadiusAxis,
  RadialBar,
  RadialBarChart,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

export const ResourceStat = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getResourceData()
          .then((response: any) => {
            // console.log("--->response.data", response.data);
            const res = response.data;
            const data = [{ name: "Resources", count: res["total"] }];
            setData(data);
          })
          .catch((error) => setError(error.message))
          .finally(() => setLoading(false));
      } catch (err: any) {
        setError(err.message);
      }
    };
    fetchData();
  }, []);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;

  return (
    <div>
      <Card
        title="Resource Statistics"
        bordered={false}
        className="w-[100%] border-1"
      >
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
            <RadialBar dataKey="count" background fill="#0d51d9" />
            <Tooltip />
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
                    );
                  }
                }}
              />
            </PolarRadiusAxis>
          </RadialBarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  );
};
