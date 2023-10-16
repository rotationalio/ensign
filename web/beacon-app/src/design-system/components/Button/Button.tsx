import React, { forwardRef, ReactNode } from 'react';

import { mergeClassnames } from '../../utils';
import Loader from '../Loader/Loader';
import { StyledButton } from './Button.styles';
import { BtnSize, BtnVariant } from './Button.types';

export type BtnProps = {
  children: ReactNode;
  variant?: BtnVariant;
  size?: BtnSize;
  className?: string;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  isLoading?: boolean;
  tabIndex?: number;
  onclick?: () => void;
  onChange?: () => void;
} & React.DetailedHTMLProps<React.ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement>;

const Button = forwardRef<HTMLButtonElement, BtnProps>((props, ref) => {
  const {
    children,
    leftIcon,
    rightIcon,
    isLoading,
    className,
    variant = 'primary',
    size = 'medium',
    ...rest
  } = props;

  return (
    <StyledButton
      size={size}
      variant={variant}
      isLoading={isLoading}
      className={mergeClassnames(
        'line-height-1.75 font-size-14 bg-inherit rounded-5 min-h-[28px] cursor-pointer px-4 py-2 text-[14px] font-bold text-white transition-colors duration-200 ease-in-out focus:outline-none ',
        variant === 'ghost' &&
          'rounded-5 h-10 w-20 border border-gray-400 text-gray-600 hover:border-gray-300  hover:bg-gray-300 disabled:border-gray-300 disabled:bg-gray-300 disabled:text-gray-400 disabled:hover:bg-gray-300',
        className
      )}
      {...rest}
      ref={ref}
    >
      {isLoading ? (
        <Loader size="xs" className="item-center m-auto text-center" />
      ) : (
        <>
          {leftIcon && <span className="pr-1">{leftIcon}</span>}
          {children}
          {rightIcon && <span className="pl-1">{rightIcon}</span>}
        </>
      )}
    </StyledButton>
  );
});

export default Button;
