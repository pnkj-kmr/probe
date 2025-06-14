// import React from 'react';
import { Input } from "antd";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import type { MenuProps } from "antd";
import { Button, Dropdown, Space, message } from "antd";
import { DownOutlined, UserOutlined } from "@ant-design/icons";

const { Search } = Input;

const handleMenuClick: MenuProps["onClick"] = (e) => {
  message.info("Click on menu item.");
  console.log("click", e);
};

const items: MenuProps["items"] = [
  {
    label: "1st menu item",
    key: "1",
    icon: <UserOutlined />,
  },
  {
    label: "2nd menu item",
    key: "2",
    icon: <UserOutlined />,
  },
  {
    label: "3rd menu item",
    key: "3",
    icon: <UserOutlined />,
    danger: true,
  },
  {
    label: "4rd menu item",
    key: "4",
    icon: <UserOutlined />,
    danger: true,
    disabled: true,
  },
];

const menuProps = {
  items,
  onClick: handleMenuClick,
};

export default function StatsPage() {
  return (
    <div className="min-h-screen flex flex-col px-0 md:px-[6%]">
      <Navbar />

      <section className="grid md:grid-cols-2 gap-6 py-2 pt-8">
        <div className="flex items-center space-x-4 gap-4">
          <Dropdown menu={menuProps} className="min-w-[200px] px-4 py-1">
            <Button>
              <Space>
                - SELECT -
                <DownOutlined />
              </Space>
            </Button>
          </Dropdown>

          <Search
            placeholder="input search text"
            enterButton="Search"
            size="large"
            loading
          />
          {/* test */}
        </div>
      </section>

      <section className="grid md:grid-cols-2 gap-6 py-2 pt-8">
        <div className="border w-[100%] h-[100%] h-min-[200px]">test</div>
        <div className="border w-[100%] h-[100%]">test</div>
      </section>

      <Footer />
    </div>
  );
}
