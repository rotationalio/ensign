/* eslint-disable unused-imports/no-unused-vars */
/* eslint-disable jsx-a11y/alt-text */
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { clsx } from 'clsx';
import React from 'react';
interface RadixDropdownMenuProps {
  items: any;
  trigger?: React.ReactNode;
  isOpen?: boolean;
  'data-cy'?: string;
  onOpenChange?: (isOpen: boolean) => void;
}

const OrganizationMenuDropdown = ({
  items,
  trigger,
  isOpen,
  'data-cy': dataCy,
  onOpenChange,
}: RadixDropdownMenuProps) => {
  //console.log('items menu', items);
  return (
    <div className="relative mt-4">
      <DropdownMenuPrimitive.Root open={isOpen} onOpenChange={onOpenChange}>
        <DropdownMenuPrimitive.Trigger className="border-none focus:ring-0">
          {trigger}
        </DropdownMenuPrimitive.Trigger>

        <DropdownMenuPrimitive.Portal>
          <DropdownMenuPrimitive.Content
            align="end"
            sideOffset={5}
            className={clsx(
              'radix-side-bottom:animate-slide-down radix-side-top:animate-slide-up',
              'shadow-md w-48 rounded-lg px-1.5 py-1 md:w-56',
              'bg-white dark:bg-gray-800'
            )}
            data-cy={dataCy}
          >
            {items?.organizationMenuItems?.map(({ name, orgId, handleSwitch }: any, i: any) => (
              <DropdownMenuPrimitive.Item
                onClick={() => handleSwitch(orgId)}
                key={`${name}-${i}`}
                className={clsx(
                  'flex w-full cursor-default select-none items-center truncate rounded-md p-2 text-base outline-none',
                  'text-gray-400 focus:bg-gray-200 dark:text-gray-500 dark:focus:bg-gray-900'
                )}
              >
                <span className="text-gray-700 dark:text-gray-300">{name}</span>
              </DropdownMenuPrimitive.Item>
            ))}

            <DropdownMenuPrimitive.Arrow />
          </DropdownMenuPrimitive.Content>
        </DropdownMenuPrimitive.Portal>
      </DropdownMenuPrimitive.Root>
    </div>
  );
};

export { OrganizationMenuDropdown };
