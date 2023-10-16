import styled from 'styled-components';

import { GridVariant } from './Grid.types';

const StyledGrid = styled.div.attrs(
  ({ columns, gapX, gapY, gap, justify, align, flow }: GridVariant) => ({
    className: `grid ${columns ? `grid-cols-${columns}` : ''} ${gapX ? `gap-x-${gapX}` : ''} ${
      gapY ? `gap-y-${gapY}` : ''
    } ${gap ? `gap-${gap}` : ''} ${justify ? `justify-${justify}` : ''} ${
      align ? `align-${align}` : ''
    } ${flow ? `grid-${flow}` : ''}`,
  })
)<GridVariant>``;

export default StyledGrid;
