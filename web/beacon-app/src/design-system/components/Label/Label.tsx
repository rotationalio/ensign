import cn from 'clsx';
import React from 'react';

import StyledLabel from './Label.styles';

export type LabelProps = {
  className?: string;
  children: React.ReactNode;
  htmlFor?: string;
};

const Label = ({ className, children, htmlFor, ...props }: LabelProps) => {
  return (
    <StyledLabel className={cn('beacon-text-sm', className)} htmlFor={htmlFor} {...props}>
      {children}
    </StyledLabel>
  );
};

export default Label;
