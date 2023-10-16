import StyledGrid from './Grid.styles';
import { GridVariant } from './Grid.types';

export type GridProps = {
  className?: string;
  children: React.ReactNode;
} & GridVariant;

export const Grid = ({
  className,
  columns,
  gapX,
  gapY,
  gap,
  justify = 'start',
  align = 'stretch',
  flow,
  children,
  ...props
}: React.PropsWithChildren<GridProps>) => {
  //   const gridColumns = columns ? getGridColumns(columns) : '';
  //   const gridFlow = flow ? getGridFlow(flow) : '';
  //   const gridGap = gap && getGridGap(gap);
  //   const gridGapX = gapX ? getGridGap(gapX) : '';
  //   const gridGapY = gapY ? getGridGap(gapY) : '';
  //   const gridJustify = justify ? `justify-${justify}` : '';
  //   const gridAlign = align ? `align-${align}` : '';
  //   const gridClassName = mergeClassnames(
  //     gridColumns,
  //     gridFlow,
  //     gridGap,
  //     gridGapX,
  //     gridGapY,
  //     gridJustify,
  //     gridAlign,
  //     className
  //   );
  return (
    <StyledGrid
      columns={columns}
      gapX={gapX}
      gapY={gapY}
      gap={gap}
      justify={justify}
      align={align}
      flow={flow}
      {...props}
    >
      {children}
    </StyledGrid>
  );
};

export default Grid;
