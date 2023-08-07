import { mergeClassnames } from '@rotational/beacon-core';
import React from 'react';
export interface TagProps {
  children: React.ReactNode;
  variant?: 'primary' | 'success' | 'warning' | 'error' | 'ghost' | 'secondary' | string;
  className?: string;
  size?: 'small' | 'medium' | 'large';
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  onClick?: () => void;
}

const Tag: React.FC<TagProps> = ({
  children,
  variant = 'primary',
  className,
  size = 'small',
  leftIcon,
  rightIcon,
  onClick,
}) => {
  const getVariantClasses = () => {
    switch (variant) {
      case 'primary':
        return 'bg-primary text-white';
      case 'secondary':
        return 'bg-primary-400 text-white';
      case 'success':
        return 'bg-green-800 text-white';
      case 'warning':
        return 'bg-warning-600 text-white';
      case 'error':
        return 'bg-danger-600 text-white';
      case 'ghost':
        return 'bg-transparent text-black';

      default:
        return 'bg-gray-400 text-black';
    }
  };

  const getSizeClasses = () => {
    switch (size) {
      case 'small':
        return 'text-xs py-1 px-2';
      case 'medium':
        return 'text-sm py-1 px-2';
      case 'large':
        return 'text-lg py-2 px-4';
      default:
        return 'text-sm py-1 px-2';
    }
  };

  const combinedClassName = mergeClassnames(
    'rounded-full py-1 px-2 inline-block',
    getSizeClasses(),
    getVariantClasses(),
    className
  );
  return (
    <>
      <span className={combinedClassName} {...(onClick && { onClick })}>
        {leftIcon && <span className="mr-1">{leftIcon}</span>}
        {children}
        {rightIcon && <span className="ml-1">{rightIcon}</span>}
      </span>
    </>
  );
};

export default Tag;
