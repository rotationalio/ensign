import { ReactNode } from 'react';

export interface AlertProps {
  children: ReactNode;
}

function Alert({ children }: AlertProps) {
  return <>{children}</>;
}

export default Alert;
