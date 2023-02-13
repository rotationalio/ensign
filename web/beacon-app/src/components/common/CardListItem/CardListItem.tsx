import { Button, Card } from '@rotational/beacon-core';
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
  description?: ReactNode;
  data?: ItemDataProps[];
  tableClassName?: string;
  contentClassName?: string;
}
function ListItemCard({
  title,
  subtitle,
  description,
  data,
  contentClassName,
  tableClassName,
}: CardListItemProps) {
  return (
    <>
      <Card
        contentClassName={twMerge(
          'w-full min-h-[200px] border border-primary-900 rounded-md p-4',
          contentClassName
        )}
      >
        {title && (
          <Card.Header>
            <h1 className="px-2 font-bold">{title}</h1>
          </Card.Header>
        )}
        <Card.Body>
          <div className="space-y-3">
            {subtitle ||
              (description && (
                <div className="my-3 mb-5 flex flex-col items-start justify-between gap-3 px-2 sm:mb-0 sm:flex-row sm:gap-0">
                  {subtitle && <p className="text-sm sm:w-4/5">{subtitle}</p>}
                  {description && (
                    <div className="sm:w-1/5">
                      <Button className="h-auto text-sm">Manage Project</Button>
                    </div>
                  )}
                </div>
              ))}
            {data && data.length > 0 && (
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
          </div>
        </Card.Body>
      </Card>
    </>
  );
}

export default ListItemCard;
