import * as RadixTooltip from '@radix-ui/react-tooltip';
import { ReactNode } from 'react';

export type TooltipProps = {
  title?: ReactNode;
} & RadixTooltip.TooltipProviderProps;
