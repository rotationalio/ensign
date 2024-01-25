/* eslint-disable unused-imports/no-unused-vars */
/* eslint-disable jsx-a11y/alt-text */
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { Link2Icon } from '@radix-ui/react-icons';
import { clsx } from 'clsx';
import React from 'react';
interface RadixDropdownMenuProps {
  items: any;
  trigger?: React.ReactNode;
  isOpen?: boolean;
  'data-cy'?: string;
  onOpenChange?: (isOpen: boolean) => void;
}

const MenuDropdownMenu = ({
  items,
  trigger,
  isOpen,
  'data-cy': dataCy,
  onOpenChange,
}: RadixDropdownMenuProps) => {
  //console.log('items menu', items);
  return (
    <div className="relative">
      <DropdownMenuPrimitive.Root open={isOpen} onOpenChange={onOpenChange}>
        <DropdownMenuPrimitive.Trigger className="border-none focus:ring-0" data-cy={dataCy}>
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
          >
            {items?.generalMenuItems?.map(({ label, icon, shortcut, onClick }: any, i: any) => (
              <DropdownMenuPrimitive.Item
                key={`${label}-${i}`}
                onClick={onClick}
                className={clsx(
                  'flex cursor-default select-none items-center rounded-md px-2 py-2 text-base outline-none',
                  'text-gray-400 focus:bg-gray-200 dark:text-gray-500 dark:focus:bg-gray-900'
                )}
              >
                {icon}
                <span className="flex-grow text-gray-700 dark:text-gray-300">{label}</span>
                {shortcut && <span className="text-base">{shortcut}</span>}
              </DropdownMenuPrimitive.Item>
            ))}

            {items?.logoutMenuItem && (
              <DropdownMenuPrimitive.Item
                onClick={items?.logoutMenuItem?.onClick}
                className={clsx(
                  'flex cursor-default select-none items-center rounded-md px-2 py-2 text-base outline-none',
                  'text-gray-400 focus:bg-gray-200 dark:text-gray-500 dark:focus:bg-gray-900'
                )}
              >
                <Link2Icon className="h-3.5 w-3.5 mr-2" />
                <span className="flex-grow text-gray-700 dark:text-gray-300">
                  {items?.logoutMenuItem?.label}
                </span>
              </DropdownMenuPrimitive.Item>
            )}
            <DropdownMenuPrimitive.Arrow />
          </DropdownMenuPrimitive.Content>
        </DropdownMenuPrimitive.Portal>
      </DropdownMenuPrimitive.Root>
    </div>
  );
};

export { MenuDropdownMenu };
