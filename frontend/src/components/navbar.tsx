import { Link } from "react-router-dom";
import {
  HomeOutlined,
  InfoCircleOutlined,
  FileSearchOutlined,
  //   UnlockOutlined,
  ApiOutlined,
} from "@ant-design/icons";

export default function Navbar() {
  return (
    <header className="border-b shadow-sm p-4 flex justify-between items-center">
      <h1 className="text-xl font-bold border-2 px-1 text-[#0d51d9]">
        <strong>P R O B E</strong>
      </h1>
      <nav className="space-x-8">
        <Link to="/">
          <HomeOutlined />
        </Link>
        <Link to="/search">
          <FileSearchOutlined />
        </Link>
        <Link to="/about">
          <InfoCircleOutlined />
        </Link>
        <Link to="/docs">
          <ApiOutlined />
        </Link>
        {/* <UnlockOutlined /> */}
      </nav>
    </header>
  );
}
