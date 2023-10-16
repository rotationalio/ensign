export interface GridVariant {
  columns?: number;
  gapX?: number;
  gapY?: number;
  gap?: number;
  justify?: 'start' | 'end' | 'center' | 'stretch' | 'between' | 'around';
  align?: 'start' | 'end' | 'center' | 'stretch';
  flow?: 'row' | 'column' | 'dense' | 'row dense' | 'column dense';
}
