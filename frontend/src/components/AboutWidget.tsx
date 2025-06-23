import React, { useEffect, useState } from "react";
import { getInfo, getEnvVariables, getProcessStatus } from "@/services/api";
import { Table, Splitter, Space } from "antd";

interface TableWrapperProps {
  dataSource: Array<any>;
  title: string;
}

const HiddenHeaderTable: React.FC<TableWrapperProps> = ({
  title,
  dataSource,
}) => {
  const columns = [
    {
      // title: () => null,
      title: title,
      dataIndex: "K",
      key: "K",
      render: (text: string) => <span>{text.toUpperCase()}</span>,
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 8px",
        },
      }),
      width: 250,
    },
    {
      title: () => null,
      dataIndex: "V",
      key: "V",
    },
  ];

  return (
    <div>
      <Table
        // className="custom-header-table"
        dataSource={dataSource}
        columns={columns}
        pagination={false}
        scroll={{ y: 500 }}
        components={{
          body: {
            row: (props: any) => <tr {...props} style={{ height: "30px" }} />,
            cell: (props: any) => (
              <td {...props} style={{ padding: "4px 8px", fontSize: "14px" }} />
            ),
          },
        }}
        // bordered
      />
    </div>
  );
};

const ProcessTable: React.FC<TableWrapperProps> = ({ title, dataSource }) => {
  const columns = [
    {
      // title: () => null,
      title: title,
      dataIndex: "name",
      key: "name",
      render: (text: string) => <span>{text.toUpperCase()}</span>,
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 8px",
        },
      }),
      width: 250,
    },
    {
      title: "RUNCYCLES",
      dataIndex: "counter",
      key: "counter",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
      width: 120,
    },
    {
      title: "TOTAL",
      dataIndex: "given",
      key: "given",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
      width: 120,
    },
    {
      title: "POLLED",
      dataIndex: "polled",
      key: "polled",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
      width: 120,
    },
    {
      title: "TT",
      dataIndex: "t",
      key: "t",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
      width: 120,
    },
    {
      title: "ST",
      dataIndex: "st",
      key: "st",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
    },
    {
      title: "ET",
      dataIndex: "et",
      key: "et",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
    },
    {
      title: "ErrorIfAny",
      dataIndex: "error",
      key: "error",
      onHeaderCell: () => ({
        style: {
          fontSize: "15px",
          padding: "8px 0px",
        },
      }),
    },
  ];

  return (
    <div>
      <Table
        // className="custom-header-table"
        dataSource={dataSource}
        columns={columns}
        pagination={false}
        scroll={{ y: 500 }}
        components={{
          body: {
            row: (props: any) => <tr {...props} style={{ height: "30px" }} />,
            cell: (props: any) => (
              <td {...props} style={{ padding: "4px 8px", fontSize: "14px" }} />
            ),
          },
        }}
        // bordered
      />
    </div>
  );
};

const AboutWidget = () => {
  const [info, setInfo] = useState<any>(null);
  const [env, setEnv] = useState<any>(null);
  const [process, setProcess] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchInfo = async () => {
    try {
      getInfo()
        .then((response: any) => {
          const res = response.data;
          setInfo(res);
        })
        .catch((error) => setError(error.message))
        .finally(() => setLoading(false));
    } catch (err: any) {
      setError(err.message);
    }
  };
  const fetchEnv = async () => {
    try {
      getEnvVariables()
        .then((response: any) => {
          const res = response.data;
          setEnv(res);
        })
        .catch((error) => setError(error.message))
        .finally(() => setLoading(false));
    } catch (err: any) {
      setError(err.message);
    }
  };
  const fetchProcess = async () => {
    try {
      getProcessStatus()
        .then((response: any) => {
          const res = response.data;
          // console.log("--->response.data", res);
          setProcess(res);
        })
        .catch((error) => setError(error.message))
        .finally(() => setLoading(false));
    } catch (err: any) {
      setError(err.message);
    }
  };

  useEffect(() => {
    fetchInfo();
    fetchEnv();
    fetchProcess();
  }, []);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;
  return (
    <>
      <Space direction="vertical" className="min-h-[68vh]">
        <Splitter lazy className="h-[40vh]">
          <Splitter.Panel defaultSize="45%" min="20%" max="70%">
            <HiddenHeaderTable dataSource={info} title="PROBE INFO" />
          </Splitter.Panel>
          <Splitter.Panel>
            <HiddenHeaderTable dataSource={env} title="ENVIRONMENT VARIABLES" />
          </Splitter.Panel>
        </Splitter>
        <Splitter lazy layout="vertical" className="min-h-[40vh]">
          <Splitter.Panel defaultSize="50%" min="30%" max="70%">
            <ProcessTable title="THREAD PROCESSES" dataSource={process} />
          </Splitter.Panel>
        </Splitter>
      </Space>
    </>
  );
};

export default AboutWidget;
