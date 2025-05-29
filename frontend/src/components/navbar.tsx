import { Link } from "react-router-dom";
// import {Button} from 'antd';
import { HomeOutlined, InfoCircleOutlined,FileSearchOutlined, UnlockOutlined } from "@ant-design/icons";


export default function Navbar() {
    return (
        <header className="border-b shadow-sm p-4 flex justify-between items-center">
            <h1 className="text-xl font-bold border-2 border-blue-500 text-blue-500 px-1">P R O B E</h1>
            <nav className="space-x-8">
                <Link to="/">
                <HomeOutlined />
                    {/* <Button ></Button> */}
                </Link>
                <Link to="/stats">
                <FileSearchOutlined />
                    {/* <Button ></Button> */}
                </Link>
                <Link to="/about">
                <InfoCircleOutlined />
                    {/* <Button ></Button> */}
                </Link>
                <UnlockOutlined  />
            {/* <Button ></Button> */}
            </nav>
        </header>
    )
}