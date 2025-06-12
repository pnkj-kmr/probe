import React from 'react';
import type { CollapseProps } from 'antd';
import { Collapse } from 'antd';

const text = `
  A dog is a type of domesticated animal.
  Known for its loyalty and faithfulness,
  it can be found as a welcome guest in many households across the world.
`;

const items: CollapseProps['items'] = [
  {
    key: '1',
    label: 'Info',
    children: <p>{text}</p>,
  },
  {
    key: '2',
    label: 'Processes',
    children: <p>{text}</p>,
  },
  {
    key: '3',
    label: 'Environment Variables',
    children: <p>{text}</p>,
  },
];

const AboutWidget: React.FC = () => {
  const onChange = (key: string | string[]) => {
    console.log(key);
  };

  return <Collapse items={items} defaultActiveKey={['1']} onChange={onChange} />;
};

export default AboutWidget;