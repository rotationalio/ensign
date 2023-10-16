import { lighten } from 'polished';
import styled, { css } from 'styled-components';

import { BtnProps } from './Button';

const primaryColor = '#1D65A6';
const secondaryColor = '#E66809';
const grayColor = '#DEE2E6';
const tertiaryColor = '#238553';

export const StyledButton = styled.button<BtnProps>`
  /* base */
  transition: background-color 200ms ease;
  border-radius: 5px;

  /* variants */
  ${(props) => getVariantStyles(props.variant)}

  /* sizes */
  ${(props) => getSizeStyles(props.size)}

  /* disabled & loading */
  ${(props) => (props.isLoading || props.disabled) && getDisabledStyles(props.variant)}
`;

const getVariantColor = (variant: BtnProps['variant'] = 'primary') => {
  const colorVarsMap = {
    primary: primaryColor,
    secondary: secondaryColor,
    ghost: grayColor,
    tertiary: tertiaryColor,
  };
  return colorVarsMap[variant];
};

const getSizeStyles = (size: BtnProps['size'] = 'medium') => {
  return {
    small: css`
      min-width: 60px;
      font-size: 12px;
      & > svg {
        width: 20px;
        height: auto;
      }
    `,
    medium: css`
      min-height: 32px;
      min-width: 100px;
      padding: 8px 12px;
    `,
    large: css`
      min-height: 56px;
      min-width: 150px;
      padding: 8px 12px;
      font-size: 18px;
    `,
    custom: css``,
  }[size];
};

const getVariantStyles = (variant: BtnProps['variant'] = 'primary') => {
  return {
    primary: css`
      background-color: ${primaryColor};
      :hover {
        background-color: ${() => lighten(0.2)(getVariantColor(variant))};
      }
    `,
    secondary: css`
      background-color: ${secondaryColor};
      :hover {
        background-color: ${() => lighten(0.2)(getVariantColor(variant))};
      }
    `,
    ghost: '',
    tertiary: css`
      background-color: ${tertiaryColor};
      :hover {
        background-color: ${() => lighten(0.2)(getVariantColor(variant))};
      }
    `,
  }[variant];
};

const getDisabledStyles = (variant: BtnProps['variant'] = 'primary') => {
  return css`
    background-color: ${lighten(0.4)(getVariantColor(variant))};
    cursor: not-allowed;
    :hover {
      background-color: ${lighten(0.4)(getVariantColor(variant))};
    }
  `;
};
