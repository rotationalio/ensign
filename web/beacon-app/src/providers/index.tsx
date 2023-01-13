import { ReactElement, ReactNode } from 'react';

type AppProvidersProps = {
  children: ReactNode;
};

function AppProviders({ children }: AppProvidersProps): ReactElement {
  return <div>{children}</div>;
}

export default AppProviders;
