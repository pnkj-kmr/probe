import React, { useEffect, useState } from "react";
import { getPollingData } from "@/services/api";
import { Card } from "antd";
import {
  BarChart,
  Bar,
  //   Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ReferenceLine,
  ResponsiveContainer,
} from "recharts";

// const chartData = [
//   { name: "safari", count: 100000, fill: "var(--color-safari)" },
// ]

// Sample data

export const PollPeriodStat = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getPollingData()
          .then((response: any) => {
            const res = response.data;
            // console.log("--->response.data", res);
            // const data = [
            //   {
            //     name: "ICMP/60",
            //     polled: 4000,
            //     total: 2400,
            //     // amt: 2400,
            //   },
            // ];
            setData(res);
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
        title="Polling Statistics"
        bordered={false}
        className="w-[100%] border-1"
      >
        <ResponsiveContainer className="min-h-[400px] w-[100%]">
          <BarChart
            width={500}
            height={300}
            data={data}
            stackOffset="sign"
            margin={{
              top: 5,
              right: 30,
              left: 20,
              bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip />
            <Legend />
            <ReferenceLine y={0} stroke="#000" />
            <Bar dataKey="total" fill="#8884d8" stackId="stack" barSize={50} />
            <Bar dataKey="polled" fill="#82ca9d" stackId="stack" barSize={50} />
          </BarChart>
        </ResponsiveContainer>
      </Card>
    </div>
  );
};

// export default App
