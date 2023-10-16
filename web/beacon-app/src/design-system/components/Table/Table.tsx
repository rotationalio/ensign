import { Column, Row, useFilters, usePagination, useSortBy, useTable } from 'react-table';

import mergeClassnames from '../../utils/mergeClassnames';
import { SortDownIcon, SortIcon, SortUpIcon } from '../Icon/Icons';
import Loader from '../Loader/Loader';
import { ActionPill } from './shared/ActionPill';
import { PaginateButton } from './shared/PaginateButton';
import { StatusPill } from './shared/StatusPill';
export interface TableProps {
  columns: Column[];
  data: any;
  initialState?: any;
  className?: string;
  tableClassName?: string;
  theadClassName?: string;
  tbodyClassName?: string;
  trClassName?: string;
  tdClassName?: string;
  thClassName?: string;
  statusClassName?: string;
  actionsClassName?: string;
  isLoading?: boolean;
  showPaginationAfter?: number;
  [key: string]: any;
  onRowClick?: (params: Row) => void;
}

function Table({
  columns,
  data,
  onRowClick,
  tableClassName,
  theadClassName,
  trClassName,
  thClassName,
  initialState,
  tdClassName,
  showPaginationAfter = 10,
  statusClassName,
  actionsClassName,
  isLoading = false,
  ...rest
}: TableProps) {
  // Use the state and functions returned from useTable to build your UI
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    page,
    canPreviousPage,
    canNextPage,
    pageOptions,

    nextPage,
    previousPage,
    setPageSize,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      initialState: { pageIndex: 0, pageSize: 10, ...initialState },
      data,
    },
    useFilters, // useFilters!
    useSortBy,
    usePagination // new
  );

  return (
    <>
      {/*  header group (filter row) */}
      <div className="sm:flex sm:gap-x-2">
        {headerGroups.map((headerGroup) =>
          headerGroup.headers.map((column) =>
            column.Filter ? (
              <div className="mt-2 sm:mt-0" key={column.id}>
                {column.render('Filter')}
              </div>
            ) : null
          )
        )}
      </div>
      {/* table */}
      <div className="mt-4 flex flex-col">
        <div className="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            <div className="shadow overflow-hidden border-b border-gray-200 text-[12px] sm:rounded-lg">
              <table
                {...getTableProps()}
                className={mergeClassnames('min-w-full divide-y divide-gray-200', tableClassName)}
                {...rest}
              >
                <thead className={mergeClassnames('bg-white font-bold', theadClassName)}>
                  {headerGroups.map((headerGroup) => (
                    <tr {...headerGroup.getHeaderGroupProps()}>
                      {headerGroup.headers.map((column) => (
                        <th
                          scope="col"
                          className={mergeClassnames(
                            'group border-gray-400/40 border-x-[1.5] px-6 py-3 text-left  text-[12px] font-bold capitalize tracking-wider',
                            thClassName
                          )}
                          {...column.getHeaderProps(column.getSortByToggleProps())}
                        >
                          <div className="flex items-center justify-between text-[12px]">
                            {column.render('Header')}
                            {/* Add a sort direction indicator */}
                            <span>
                              {column.isSorted ? (
                                column.isSortedDesc ? (
                                  <SortDownIcon className="h-4 w-4 text-gray-400" />
                                ) : (
                                  <SortUpIcon className="h-4 w-4 text-gray-400" />
                                )
                              ) : (
                                <SortIcon className="h-4 w-4 text-gray-400 opacity-0 group-hover:opacity-100" />
                              )}
                            </span>
                          </div>
                        </th>
                      ))}
                    </tr>
                  ))}
                </thead>
                <tbody {...getTableBodyProps()} className="divide-gray-[#1C1C1C] divide-y bg-white">
                  {isLoading && (
                    <tr className="items-center text-center">
                      <td
                        colSpan={columns.length}
                        className="px-auto h-[100px] items-center whitespace-nowrap py-4 text-center"
                      >
                        <div className="flex flex-col items-center justify-center text-center">
                          <Loader size="md" />
                        </div>
                        <p className="mt-2 text-sm text-gray-500">
                          Please wait while we load the data
                        </p>
                      </td>
                    </tr>
                  )}
                  {(page.length > 0 &&
                    page.map((row, i) => {
                      prepareRow(row);

                      return (
                        <tr
                          {...row.getRowProps()}
                          onClick={() => onRowClick && onRowClick(row)}
                          className={mergeClassnames(
                            onRowClick && 'cursor-pointer hover:bg-gray-100',
                            trClassName
                          )}
                        >
                          {!isLoading &&
                            row.cells.length > 0 &&
                            row.cells.map((cell) => {
                              return (
                                <td
                                  {...cell.getCellProps()}
                                  className={mergeClassnames(
                                    'whitespace-nowrap px-6 py-4 text-[12px]',
                                    tdClassName
                                  )}
                                  role="cell"
                                >
                                  {{
                                    status: (
                                      <StatusPill value={cell.value} className={statusClassName} />
                                    ),
                                    actions: (
                                      <ActionPill
                                        actions={cell.value}
                                        className={actionsClassName}
                                      />
                                    ),
                                    default: cell.render('Cell'),
                                  }[cell.column.id] || cell.render('Cell')}
                                </td>
                              );
                            })}
                        </tr>
                      );
                    })) || (
                    <tr className="items-center text-center">
                      {!isLoading && (
                        <td
                          colSpan={columns.length}
                          className="px-auto h-[100px] items-center whitespace-nowrap py-4 text-center"
                        >
                          <div className="flex flex-col items-center justify-center text-center">
                            <p className="mt-2 text-sm font-semibold text-gray-800">
                              No data available
                            </p>
                          </div>
                        </td>
                      )}
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
        {/* display pagination only when data is up to pagesize value */}
        {data.length > showPaginationAfter && (
          <div className="px-auto flex items-center justify-between py-3">
            <PaginateButton onClick={() => previousPage()} disabled={!canPreviousPage}>
              Previous
            </PaginateButton>

            <div className="flex items-center justify-between">
              <div className="flex items-baseline gap-x-2">
                <span className="text-sm text-gray-700">
                  Page <span className="font-medium">{pageIndex + 1}</span> of{' '}
                  <span className="font-medium">{pageOptions.length}</span>
                </span>
                <label>
                  <span className="sr-only">Items Per Page</span>
                  <select
                    className="shadow-sm mt-1 block w-full  rounded-md border-gray-300 text-[12px] font-medium leading-[18px] text-gray-700 focus:border-blue-500 focus:outline-none focus:ring focus:ring-blue-500 focus:ring-opacity-50 sm:text-sm"
                    value={pageSize}
                    onChange={(e) => {
                      setPageSize(Number(e.target.value));
                    }}
                  >
                    {[5, 10, 20].map((pageSize) => (
                      <option key={pageSize} value={pageSize}>
                        Show {pageSize}
                      </option>
                    ))}
                  </select>
                </label>
              </div>
            </div>
            <PaginateButton onClick={() => nextPage()} disabled={!canNextPage}>
              Next
            </PaginateButton>
          </div>
        )}
      </div>
    </>
  );
}

export default Table;
