import { useEffect, useState } from "react";
import { getResourceData } from "@/services/api";
import { Card } from "antd";
import {
  PolarAngleAxis,
  PolarGrid,
  Radar,
  RadarChart,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

export const PollTypeStat = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getResourceData()
          .then((response: any) => {
            const res = response.data;
            // console.log("--->response.data", res);
            const data = res["stats"];
            // data.push({ name: "DUMMY", total: 100 });
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
        title="Poll Types Statistics"
        bordered={false}
        className="w-[100%] border-1"
      >
        <ResponsiveContainer className="min-h-[220px] w-[100%]">
          <RadarChart data={data}>
            {/* <ChartTooltip cursor={false} content={<ChartTooltipContent />} /> */}
            <PolarAngleAxis dataKey="name" />
            <PolarGrid />
            <Tooltip />
            <Radar
              fill="#0d51d9"
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
  );
};
