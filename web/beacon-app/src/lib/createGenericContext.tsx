import React from 'react';
interface Props {
  value: any;
  children: React.ReactNode;
}
export const createGenericContext = (defaultValue: any) => {
  const context = React.createContext(defaultValue);
  const Provider = ({ value, children }: Props) => {
    return React.createElement(context.Provider, { value }, children);
  };
  return { context, Provider };
};
