/* eslint-disable jsx-a11y/alt-text */
import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import { CaretRightIcon, Link2Icon, MixerHorizontalIcon, PersonIcon } from '@radix-ui/react-icons';
import { Button } from '@rotational/beacon-core';
import { clsx } from 'clsx';
import React, { ReactNode } from 'react';

import { ChevronDown } from '@/components/icons/chevron-down';

interface RadixMenuItem {
  label: string;
  shortcut?: string;
  icon?: ReactNode;
}

interface User {
  name: string;
  url?: string;
}

const generalMenuItems: RadixMenuItem[] = [
  {
    label: 'Settings',
    icon: <MixerHorizontalIcon className="h-3.5 w-3.5 mr-2" />,
  },
];

const users: User[] = [
  {
    name: 'Noorazi',
    url: 'https://github.com/adamwathan.png',
  },
  {
    name: 'Rotational LLC',
    url: 'https://github.com/steveschoger.png',
  },
  {
    name: 'Patrick INC',
    url: 'https://github.com/robinmalfait.png',
  },
];

const MenuDropdownMenu = () => {
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
            {generalMenuItems.map(({ label, icon, shortcut }, i) => (
              <DropdownMenuPrimitive.Item
                key={`${label}-${i}`}
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
                  {users.map(({ name, url }, i) => (
                    <DropdownMenuPrimitive.Item
                      key={`${name}-${i}`}
                      className={clsx(
                        'flex w-full cursor-default select-none items-center truncate rounded-md px-2 py-2 text-xs outline-none',
                        'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
                      )}
                    >
                      {url ? (
                        <img className="mr-2.5 h-6 w-6 rounded-full" src={url} />
                      ) : (
                        <PersonIcon className="mr-2.5 h-6 w-6" />
                      )}
                      <span className="text-gray-700 dark:text-gray-300">{name}</span>
                    </DropdownMenuPrimitive.Item>
                  ))}
                </DropdownMenuPrimitive.SubContent>
              </DropdownMenuPrimitive.Portal>
            </DropdownMenuPrimitive.Sub>
            <DropdownMenuPrimitive.Separator className="my-1 h-px bg-gray-200 dark:bg-gray-700" />
            <DropdownMenuPrimitive.Item
              className={clsx(
                'flex cursor-default select-none items-center rounded-md px-2 py-2 text-xs outline-none',
                'focus:bg-gray-50 text-gray-400 dark:text-gray-500 dark:focus:bg-gray-900'
              )}
            >
              <Link2Icon className="h-3.5 w-3.5 mr-2" />
              <span className="flex-grow text-gray-700 dark:text-gray-300">Logout</span>
            </DropdownMenuPrimitive.Item>
          </DropdownMenuPrimitive.Content>
        </DropdownMenuPrimitive.Portal>
      </DropdownMenuPrimitive.Root>
    </div>
  );
};

export { MenuDropdownMenu };
