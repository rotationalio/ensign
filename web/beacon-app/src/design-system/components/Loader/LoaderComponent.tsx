import React, { ReactNode } from 'react';
import mergeClassnames from 'utils/mergeClassnames';

import StyledLoader from './Loader';

export type LoaderProps = {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  color?: string;
  className?: string;
  [key: string]: any;
  content?: ReactNode;
} & React.HTMLAttributes<HTMLDivElement>;

const Loader = (props: LoaderProps) => {
  const { size = 'md', className, variant, content, ...rest } = props;
  return (
    <div className="flex flex-col ">
      <StyledLoader
        className={mergeClassnames('flex flex-row items-center justify-center', className)}
        size={size}
        variant={variant}
        {...rest}
      >
        {content}
      </StyledLoader>
      {content && <p className="text-center text-sm">{content}</p>}
    </div>
  );
};

export default Loader;
