type GridPosition = 'x' | 'y';
export const setGridColumns = (columns: number) => {
  switch (columns) {
    case 1:
      return 'grid-cols-1';

    case 2:
      return 'grid-cols-2';

    case 3:
      return 'grid-cols-3';
  }
};

export const setGridRows = (rows: number) => {
  switch (rows) {
    case 1:
      return 'grid-rows-1';

    case 2:
      return 'grid-rows-2';

    case 3:
      return 'grid-rows-3';
  }
};

export const setGridGap = (gap: number, position?: GridPosition) => {
  switch (gap) {
    case 1:
      if (position === 'x') {
        return 'gap-x-1';
      }
      if (position === 'y') {
        return 'gap-y-1';
      }
      return 'gap-1';
    case 2:
      if (position === 'x') {
        return 'gap-x-2';
      }
      if (position === 'y') {
        return 'gap-y-2';
      }
      return 'gap-2';
    case 3:
      if (position === 'x') {
        return 'gap-x-3';
      }
      if (position === 'y') {
        return 'gap-y-3';
      }
      return 'gap-3';
  }
};
export const setGridAutoFlow = (autoFlow: string) => {
  switch (autoFlow) {
    case 'row':
      return 'grid-flow-row';

    case 'column':
      return 'grid-flow-col';
  }
};

export const setGridAutoRows = (autoRows: number) => {
  switch (autoRows) {
    case 1:
      return 'auto-rows-1';

    case 2:
      return 'auto-rows-2';

    case 3:
      return 'auto-rows-3';
  }
};
export const setGridAutoColumns = (autoColumns: number) => {
  switch (autoColumns) {
    case 1:
      return 'auto-cols-1';

    case 2:
      return 'auto-cols-2';

    case 3:
      return 'auto-cols-3';
  }
};

export const setGridTemplateRows = (templateRows: number) => {
  switch (templateRows) {
    case 1:
      return 'grid-rows-1';

    case 2:
      return 'grid-rows-2';

    case 3:
      return 'grid-rows-3';
  }
};

export const setGridTemplateColumns = (templateColumns: number) => {
  switch (templateColumns) {
    case 1:
      return 'grid-cols-1';

    case 2:
      return 'grid-cols-2';

    case 3:
      return 'grid-cols-3';
  }
};

export const setGridFlow = (flow: string) => {
  switch (flow) {
    case 'row':
      return 'grid-flow-row';

    case 'column':
      return 'grid-flow-col';
    case 'row-dense':
      return 'grid-flow-row-dense';
    case 'column-dense':
      return 'grid-flow-col-dense';
  }
};
