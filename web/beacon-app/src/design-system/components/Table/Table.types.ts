export type ColumnProps = {
  Header: string;
  accessor: string;
  Filter?: any;
  Cell?: any;
  disableSortBy?: boolean;
  disableFilters?: boolean;
  className?: string;
  width?: string;
  minWidth?: string;
  maxWidth?: string;
  align?: string;
  status?: boolean;
  actions?: any;
};
