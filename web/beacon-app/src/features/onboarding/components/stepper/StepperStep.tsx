import React from 'react';

import Indicator from './Indicator';

export interface StepperProps {
  title: string;
  description: string;
  value?: string;
}

const StepperStep = ({ title, description, value }: StepperProps) => {
  return (
    <>
      <li className="mb-10 ml-6">
        <Indicator />
        <h3 className="font-medium leading-tight">{title}</h3>
        <p className="text-sm ">{description}</p>
        {value && <p className="text-sm font-bold">{value}</p>}
      </li>
    </>
  );
};

export default StepperStep;
