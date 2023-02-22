import { Trans } from '@lingui/macro';
import { Card, Heading } from '@rotational/beacon-core';
import { ReactNode } from 'react';
import { twMerge } from 'tailwind-merge';
// temporary component to replace later when we have a real ListItem component in the design system

export interface ItemDataProps {
  label: string;
  value: string;
}
export interface CardListItemProps {
  title?: string;
  subtitle?: string;
  children?: ReactNode;
  data?: ItemDataProps[];
  tableClassName?: string;
  contentClassName?: string;
  className?: string;
}
function ListItemCard({
  title,
  children,
  data,
  contentClassName,
  tableClassName,
  className,
}: CardListItemProps) {
  return (
    <>
      <Card
        contentClassName={twMerge('w-full', contentClassName)}
        className={twMerge('min-h-[200px] w-full rounded-md border-2 p-4', className)}
      >
        {title && (
          <Card.Header>
            <Heading as="h3" className="px-2 font-bold">
              {title}
            </Heading>
          </Card.Header>
        )}
        <Card.Body>
          <div className="space-y-3">
            {children}
            {data && Object.keys(data).length > 0 && (
              <table
                className={twMerge(
                  'border-separate border-spacing-x-2 border-spacing-y-1 text-sm',
                  tableClassName
                )}
              >
                {data.map((item: ItemDataProps, index: number) => (
                  <tr key={index}>
                    <td className="font-bold">{item.label}</td>
                    <td>{item.value}</td>
                  </tr>
                ))}
              </table>
            )}
            {data && Object.keys(data).length === 0 && (
              <div className="ml-5 mt-5">
                <p className="text-sm font-bold text-danger-500">
                  <Trans>No data available, please try again later or contact support.</Trans>
                </p>
              </div>
            )}
          </div>
        </Card.Body>
      </Card>
    </>
  );
}

export default ListItemCard;
