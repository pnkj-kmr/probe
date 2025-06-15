import { useEffect, useState } from "react";
import { getSystemData } from "@/services/api";
import { Card } from "antd";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
  LabelList,
} from "recharts";

export const PerfStat = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getSystemData()
          .then((response: any) => {
            // console.log("--->response.data", response.data);
            const res = response.data;
            const data = [
              {
                name: "CPU",
                display: `Cores: ${res["cpu"]["cores"]} | ${res["cpu"]["usage"]}%`,
                usage: res["cpu"]["usage"],
              },
              {
                name: "MEMORY",
                display: `${res["memory"]["used"]}/${res["memory"]["total"]} | ${res["memory"]["used_percent"]}%`,
                usage: res["memory"]["used_percent"],
              },
              {
                name: "DISK",
                display: `${res["disk"]["used"]}/${res["disk"]["total"]} | ${res["disk"]["used_percent"]}%`,
                usage: res["disk"]["used_percent"],
              },
            ];
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
        title="System Performace"
        bordered={false}
        className="w-[100%] border-1"
      >
        <ResponsiveContainer className="min-h-[220px] w-[100%]">
          <BarChart
            accessibilityLayer
            layout="vertical"
            data={data}
            // margin={{ top: 20, right: 30, left: 0, bottom: 5 }}
            margin={{
              right: 30,
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
            {/* <Tooltip /> */}
            <Tooltip
              formatter={(_v, _n, x) => [
                `${JSON.stringify(x.payload?.display)}`,
              ]}
              labelFormatter={(label) => label}
            />
            {/* <Bar dataKey="sales" fill="#1890ff" barSize={40} /> */}

            <Bar
              dataKey="usage"
              layout="vertical"
              barSize={40}
              radius={4}
              fill="#0d51d9"
            >
              <LabelList
                dataKey="name"
                position="insideLeft"
                offset={8}
                fontSize={12}
                fill="#ffffff"
              />
              <LabelList
                dataKey="usage"
                position="right"
                offset={8}
                fontSize={12}
              />
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  );
};
