import * as RadixTooltip from '@radix-ui/react-tooltip';
import styled, { keyframes } from 'styled-components';

export const slideUpAndFade = keyframes({
  '0%': { opacity: 0, transform: 'translateY(2px)' },
  '100%': { opacity: 1, transform: 'translateY(0)' },
});

export const slideRightAndFade = keyframes({
  '0%': { opacity: 0, transform: 'translateX(-2px)' },
  '100%': { opacity: 1, transform: 'translateX(0)' },
});

export const slideDownAndFade = keyframes({
  '0%': { opacity: 0, transform: 'translateY(-2px)' },
  '100%': { opacity: 1, transform: 'translateY(0)' },
});

export const slideLeftAndFade = keyframes({
  '0%': { opacity: 0, transform: 'translateX(2px)' },
  '100%': { opacity: 1, transform: 'translateX(0)' },
});

export const StyledContent = styled(RadixTooltip.Content)`
  border-radius: 4px;
  padding: 10px 15px;
  font-size: 12px;
  line-height: 1.25;
  color: #000;
  background-color: white;
  max-width: 250px;
  box-shadow: hsl(206 22% 7% / 35%) 0px 10px 38px -10px, hsl(206 22% 7% / 20%) 0px 10px 20px -15px;
  user-select: none;
  animation-duration: 400ms;
  animation-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  will-change: transform, opacity;
  '&[data-state="delayed-open"]': {
    '&[data-side="top"]': { animationName: ${slideDownAndFade} },
    '&[data-side="right"]': { animationName: ${slideLeftAndFade} },
    '&[data-side="bottom"]': { animationName: ${slideUpAndFade} },
    '&[data-side="left"]': { animationName: ${slideRightAndFade} },
  },
`;

export const StyledArrow = styled(RadixTooltip.Arrow)`
  fill: #fff;
`;
