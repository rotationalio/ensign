import { ToastProps as RadixToastProps } from '@radix-ui/react-toast';
export type ToastProps = {
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info';
  size?: 'small' | 'medium' | 'large';
  hasIcon?: boolean;
  icon?: React.ReactNode;
  [key: string]: any;
  title?: string;
  placement?: 'up' | 'down' | 'left' | 'right';
  description?: string;
  children?: React.ReactNode;
  onClose?: () => void;
};

export type ToastWithRadixProps = ToastProps & RadixToastProps;
