import { forwardRef } from 'react';

import mergeClassnames from '../../utils/mergeClassnames';
import { ContainerVariant } from './Container.types';
import { CONTAINER_VARIANT, setVariantStyle } from './utils';
export type ContainerProps = {
  children: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
  as?: React.ElementType;
  [key: string]: any; // key value pair for any other props like data-testid or aria-label etc
  variant?: ContainerVariant;
  dash?: boolean;
  max?: string; // max-width in px
  min?: string; // min-width in px
  centred?: boolean; // centred container
  brk?: 'sm' | 'md' | 'lg' | 'xl'; // breakpoint
} & React.HTMLAttributes<HTMLDivElement>;

const Container = forwardRef((props: ContainerProps, ref: any) => {
  const { children, className, style, variant, dash, brk, as, min, max, centred, ...rest } = props;
  const Component = as || 'div';
  const mergeClassname = mergeClassnames(
    brk ? `container-${brk} mx-auto px-4` : 'container mx-auto px-4',
    max && `max-w-[${max}px]`,
    min && `min-w-[${min}px]`,
    centred && 'flex justify-center items-center place-content-center',
    dash && setVariantStyle(CONTAINER_VARIANT.DASH),
    variant && setVariantStyle(variant || CONTAINER_VARIANT.DEFAULT),
    className
  );
  return (
    <Component className={mergeClassname} style={style} ref={ref} {...rest}>
      {children}
    </Component>
  );
});

export default Container;
