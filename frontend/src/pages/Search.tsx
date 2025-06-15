import React, { useEffect, useState } from "react";
import { getSearchModels, getSearchResult } from "@/services/api";
import { Input } from "antd";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import { Select } from "antd";
import type { GetProps } from "antd";

type SearchProps = GetProps<typeof Input.Search>;
const { Search } = Input;

export default function SearchPage() {
  const [models, setModels] = useState<any>(null);
  const [selected, setSelected] = useState<any>(null);
  const [result, setResult] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        getSearchModels()
          .then((response: any) => {
            const res = response.data;
            setModels(res);
          })
          .catch((error) => setError(error.message))
          .finally(() => setLoading(false));
      } catch (err: any) {
        setError(err.message);
      }
    };
    fetchData();
  }, []);

  const handleChange = (value: { value: string; label: React.ReactNode }) => {
    setSelected(value);
  };
  const onSearch: SearchProps["onSearch"] = (value, _e, _info) => {
    // console.log(info?.source, value);
    try {
      getSearchResult({ params: { table: selected.value, cid: value } })
        .then((response: any) => {
          const res = response.data;
          // console.log("--->response.data", res);
          // const models = res["stats"];
          setResult(res);
        })
        .catch((error) => setError(error.message))
        .finally(() => setLoading(false));
    } catch (err: any) {
      setError(err.message);
    }
  };

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;
  return (
    <div className="min-h-screen flex flex-col px-0 md:px-[6%]">
      <Navbar />

      <section className="grid md:grid-cols-2 gap-6 py-2 pt-8">
        <div className="flex items-center space-x-4 gap-4 w-min-4">
          <Select
            labelInValue
            defaultValue={{ value: "-", label: "- SELECT -" }}
            style={{ width: 120 }}
            onChange={handleChange}
            options={models}
          />

          <Search
            placeholder="enter ID to search"
            enterButton="Search"
            onSearch={onSearch}
            // size="large"
            // loading={{ loading }}
          />
        </div>
      </section>

      <section className="grid md:grid-cols-2 gap-6 py-2 pt-8 h-min-[300px]">
        {result && (
          <>
            <div className="w-[100%] h-[100%]">
              <pre
                style={{
                  background: "#f4f4f4",
                  padding: "1em",
                  borderRadius: "5px",
                }}
              >
                {
                  <p className="text-[#0d51d9]">
                    <b>{selected.label} Configuration</b>
                  </p>
                }
                {JSON.stringify(result["config"], null, 2)}
              </pre>
            </div>
            <div className="w-[100%] h-[100%]">
              <pre
                style={{
                  background: "#f4f4f4",
                  padding: "1em",
                  borderRadius: "5px",
                }}
              >
                {
                  <p className="text-[#0d51d9]">
                    <b>{selected.label} Polled result if any</b>
                  </p>
                }
                {JSON.stringify(result["polled"], null, 2)}
              </pre>
            </div>
          </>
        )}
        {!result && <p>Search result will be displayed here...</p>}
      </section>

      <Footer />
    </div>
  );
}
