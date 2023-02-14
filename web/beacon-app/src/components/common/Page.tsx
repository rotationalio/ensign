import { Container } from '@rotational/beacon-core';
import React, { ReactNode } from 'react';

type PageProps = {
  children: ReactNode;
};

function Page({ children }: PageProps) {
  return <Container className="mt-8 px-28">{children}</Container>;
}

export default Page;
