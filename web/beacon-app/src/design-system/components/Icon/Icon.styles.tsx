import styled from 'styled-components';

import { StyledSVGProps } from './Icon.types';
const StyledSVG = styled.svg<StyledSVGProps>((props) => ({
  ...(props.variant === 'default' && {
    fill: 'none',
    stroke: 'currentColor',
  }),
  ...(props.variant === 'primary' && {
    fill: 'none',
    stroke: 'var(--colors-primary)',
  }),
  ...(props.variant === 'secondary' && {
    fill: 'none',
    stroke: 'var(--colors-secondary)',
  }),
  ...(props.variant === 'success' && {
    fill: 'none',
    stroke: 'var(--colors-green)',
  }),
  ...(props.variant === 'danger' && {
    fill: 'none',
    stroke: 'var(--colors-red)',
  }),
  ...(props.variant === 'warning' && {
    fill: 'none',
    stroke: 'var(--colors-yellow)',
  }),
  ...(props.variant === 'info' && {
    fill: 'none',
    stroke: 'var(--colors-blue)',
  }),
  ...(props.outline && {
    fill: 'none !important',
  }),
  ...(props.size &&
    {
      1: {
        width: 'var(--spacings-1)',
        height: 'var(--spacings-1)',
      },
      2: {
        width: 'var(--spacings-2)',
        height: 'var(--spacings-2)',
      },
      3: {
        width: 'var(--spacings-3)',
        height: 'var(--spacings-3)',
      },
      4: {
        width: 'var(--spacings-4)',
        height: 'var(--spacings-4)',
      },
      5: {
        width: 'var(--spacings-5)',
        height: 'var(--spacings-5)',
      },
      6: {
        width: 'var(--spacings-6)',
        height: 'var(--spacings-6)',
      },
      7: {
        width: 'var(--spacings-7)',
        height: 'var(--spacings-7)',
      },
      8: {
        width: 'var(--spacings-8)',
        height: 'var(--spacings-8)',
      },
      9: {
        width: 'var(--spacings-9)',
        height: 'var(--spacings-9)',
      },
      10: {
        width: 'var(--spacings-10)',
        height: 'var(--spacings-10)',
      },
    }[props.size]),
}));

StyledSVG.defaultProps = {
  variant: 'default',
  outline: true,
  size: 5,
};

StyledSVG.displayName = 'StyledSVG';

export default StyledSVG;
