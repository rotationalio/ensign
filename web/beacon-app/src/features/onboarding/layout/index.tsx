import { Container } from '@rotational/beacon-core';
import React from 'react';

type Props = {
  children: React.ReactNode;
};

const OnboardingFormLayout = ({ children }: Props) => {
  return (
    <>
      <Container max={696} centered className="m-10 mt-20 px-4 xl:mt-10 xl:px-28">
        <div className="border-[1px] border-[#72A2C0] p-10">{children}</div>
      </Container>
    </>
  );
};

export default OnboardingFormLayout;
