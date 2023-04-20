/* eslint-disable unused-imports/no-unused-vars */
/* eslint-disable jsx-a11y/alt-text */
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { CaretRightIcon, Link2Icon } from '@radix-ui/react-icons';
import { Button } from '@rotational/beacon-core';
import { clsx } from 'clsx';
import React from 'react';

import { ChevronDown } from '@/components/icons/chevron-down';
interface RadixDropdownMenuProps {
  items: any;
}

const MenuDropdownMenu = ({ items }: RadixDropdownMenuProps) => {
  console.log('items menu', items);
  return (
    <div className="relative">
      <DropdownMenuPrimitive.Root>
        <DropdownMenuPrimitive.Trigger asChild>
          <Button variant="ghost" className="bg-transparent w-16 border-none text-white">
            <ChevronDown />
          </Button>
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
                  'flex cursor-default select-none items-center rounded-md px-2 py-2 text-xs outline-none',
                  'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
                )}
              >
                {icon}
                <span className="flex-grow text-gray-700 dark:text-gray-300">{label}</span>
                {shortcut && <span className="text-xs">{shortcut}</span>}
              </DropdownMenuPrimitive.Item>
            ))}

            {items?.organizationMenuItems?.length > 0 && (
              <>
                <DropdownMenuPrimitive.Separator className="my-1 h-px bg-gray-200 dark:bg-gray-700" />
                <DropdownMenuPrimitive.Label className="select-none px-2 py-2 text-xs text-gray-700 dark:text-gray-200">
                  Organization Management
                </DropdownMenuPrimitive.Label>
                <DropdownMenuPrimitive.Sub>
                  <DropdownMenuPrimitive.SubTrigger
                    className={clsx(
                      'flex w-full cursor-default select-none items-center rounded-md px-2 py-2 text-xs outline-none',
                      'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
                    )}
                  >
                    <Link2Icon className="h-3.5 w-3.5 mr-2" />
                    <span className="flex-grow text-gray-700 dark:text-gray-300">
                      Switch Organization
                    </span>
                    <CaretRightIcon className="h-3.5 w-3.5" />
                  </DropdownMenuPrimitive.SubTrigger>
                  <DropdownMenuPrimitive.Portal>
                    <DropdownMenuPrimitive.SubContent
                      className={clsx(
                        'radix-side-right:animate-scale-in·origin-radix-dropdown-menu',
                        'shadow-md w-48 rounded-md px-1 py-1 text-xs',
                        'bg-white dark:bg-gray-800'
                      )}
                    >
                      {items?.organizationMenuItems?.map(
                        ({ name, orgId, handleSwitch }: any, i: any) => (
                          <DropdownMenuPrimitive.Item
                            onClick={() => handleSwitch(orgId)}
                            key={`${name}-${i}`}
                            className={clsx(
                              'flex w-full cursor-default select-none items-center truncate rounded-md px-2 py-2 text-xs outline-none',
                              'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
                            )}
                          >
                            <span className="text-gray-700 dark:text-gray-300">{name}</span>
                          </DropdownMenuPrimitive.Item>
                        )
                      )}
                    </DropdownMenuPrimitive.SubContent>
                  </DropdownMenuPrimitive.Portal>
                </DropdownMenuPrimitive.Sub>
                <DropdownMenuPrimitive.Separator className="my-1 h-px bg-gray-200 dark:bg-gray-700" />
              </>
            )}
            {items?.logoutMenuItem && (
              <DropdownMenuPrimitive.Item
                onClick={items?.logoutMenuItem?.onClick}
                className={clsx(
                  'flex cursor-default select-none items-center rounded-md px-2 py-2 text-xs outline-none',
                  'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
                )}
              >
                <Link2Icon className="h-3.5 w-3.5 mr-2" />
                <span className="flex-grow text-gray-700 dark:text-gray-300">
                  {items?.logoutMenuItem?.label}
                </span>
              </DropdownMenuPrimitive.Item>
            )}
          </DropdownMenuPrimitive.Content>
        </DropdownMenuPrimitive.Portal>
      </DropdownMenuPrimitive.Root>
    </div>
  );
};

export { MenuDropdownMenu };