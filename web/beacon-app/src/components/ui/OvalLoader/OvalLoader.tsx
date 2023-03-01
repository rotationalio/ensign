import { ReactNode } from 'react';

import Oval from '@/components/icons/oval';

type ContainerProps = {
  className?: string;
  style?: React.CSSProperties;
};

type OvalLoaderProps = {
  containerProps?: ContainerProps;
  children?: ReactNode;
} & React.SVGProps<SVGSVGElement>;

function OvalLoader({ containerProps, children, ...rest }: OvalLoaderProps) {
  return (
    <div {...containerProps}>
      <Oval {...rest} />
      <p>{children}</p>
    </div>
  );
}

export default OvalLoader;
