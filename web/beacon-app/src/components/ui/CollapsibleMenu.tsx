// import { useRouter } from 'next/router';
import cn from 'classnames';
import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import useMeasure from 'react-use/lib/useMeasure';
import { twMerge } from 'tailwind-merge';

import { ChevronDown } from '../icons/chevron-down';
import ExternalIcon from '../icons/external-icon';

type MenuItemProps = {
  name?: string;
  icon: React.ReactNode;
  href: string;
  isMail?: boolean;
  dropdownItems?: DropdownItemProps[];
  isExternal?: boolean;
};

type DropdownItemProps = {
  name: string;
  href: string;
};

export function MenuItem({ name, icon, href, dropdownItems, isExternal, isMail }: MenuItemProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [ref, { height }] = useMeasure<HTMLUListElement>();
  const location = useLocation();

  const isCurrentPath = location.pathname === href;

  const isChildrenActive =
    dropdownItems && dropdownItems.some((item) => item.href === location.pathname);

  useEffect(() => {
    if (isChildrenActive) {
      setIsOpen(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="mb-2 min-h-[8px] list-none text-white last:mb-0">
      {dropdownItems?.length ? (
        <>
          <div
            onClick={() => setIsOpen(!isOpen)}
            role="button"
            aria-hidden="true"
            className={cn(
              'relative flex h-12 cursor-pointer items-center justify-between whitespace-nowrap rounded-lg  px-4 text-sm text-white transition-all hover:font-bold',
              isChildrenActive
                ? 'text-white'
                : 'hover:text-brand text-gray-500 dark:hover:text-white'
            )}
          >
            <span className="z-[1] flex items-center ltr:mr-3 rtl:ml-3">
              <span className="ltr:mr-3 rtl:ml-3">{icon}</span>
              {name}
            </span>
            <span
              className={`z-[1] transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`}
            >
              <ChevronDown />
            </span>

            {isChildrenActive && (
              <motion.span
                className="bg-brand shadow-large absolute bottom-0 left-0 right-0 h-full w-full rounded-lg"
                layoutId="menu-item-active-indicator"
              />
            )}
          </div>
          <div
            style={{
              height: isOpen ? height : 0,
            }}
            className="ease-[cubic-bezier(0.33, 1, 0.68, 1)] overflow-hidden transition-all duration-[350ms]"
          >
            <ul ref={ref}>
              {dropdownItems.map((item, index) => (
                <li className="first:pt-2" key={index}>
                  <NavLink
                    to={href}
                    className={({ isActive }) =>
                      !isActive
                        ? 'hover:text-brand flex items-center rounded-lg p-3 text-sm text-gray-500 transition-all before:h-1 before:w-1 before:rounded-full before:bg-gray-500 ltr:pl-6 before:ltr:mr-5 rtl:pr-6 before:rtl:ml-5 dark:hover:text-white'
                        : '!text-brand before:!bg-brand  before:-ml-0.5 before:!h-2 before:!w-2 before:ltr:!mr-[18px] before:rtl:!ml-[18px] dark:!text-white dark:before:!bg-white'
                    }
                  >
                    {item.name}
                  </NavLink>
                </li>
              ))}
            </ul>
          </div>
        </>
      ) : (
        <>
          {isMail ? (
            <a
              href={`mailto:${href}`}
              className="flex h-12 items-center whitespace-nowrap pl-8 text-sm"
            >
              <span className="relative z-[1] mr-3 w-[24px] text-white">{icon}</span>
              <span className={'relative z-[1] flex'}>
                {name}
                <ExternalIcon className="ml-1 h-3 w-3 text-white" />
              </span>
            </a>
          ) : (
            <NavLink
              to={href}
              target={isExternal ? '_blank' : '_self'}
              rel="noopener noreferrer"
              className={({ isActive }) =>
                cn(
                  `${isActive ? 'transition-all' : 'text-secondary-900'}`,
                  'relative flex h-12 items-center whitespace-nowrap pl-8 text-sm text-secondary-900 text-white'
                )
              }
            >
              <span className="relative z-[1] mr-3 w-[24px] text-white">{icon}</span>
              <span
                className={twMerge(
                  'relative z-[1] flex',
                  isCurrentPath ? 'font-bold' : 'font-normal'
                )}
              >
                {name} {isExternal && <ExternalIcon className="ml-1 h-3 w-3 text-white" />}
              </span>

              {isCurrentPath && (
                <motion.span
                  className="absolute bottom-0 left-0 right-0 h-full w-full border-l-4 border-white bg-blue-500 font-bold shadow-1"
                  layoutId="menu-item-active-indicator"
                />
              )}
            </NavLink>
          )}
        </>
      )}
    </div>
  );
}
