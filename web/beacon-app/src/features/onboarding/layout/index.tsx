import { Container } from '@rotational/beacon-core';
import React from 'react';

type Props = {
  children: React.ReactNode;
};

const OnboardingFormLayout = ({ children }: Props) => {
  return (
    <>
      <Container max={696} centered className="mt-10 flex w-full px-4 xl:m-10 xl:mt-20 xl:px-28">
        <div className="w-[996px] border-[1px] border-[#72A2C0] p-10 px-5">{children}</div>
      </Container>
    </>
  );
};

export default OnboardingFormLayout;
