// import { useRouter } from 'next/router';
import cn from 'classnames';
import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import useMeasure from 'react-use/lib/useMeasure';

import { ChevronDown } from '../icons/chevron-down';
import ExternalIcon from '../icons/external-icon';

type MenuItemProps = {
  name?: string;
  icon: React.ReactNode;
  href: string;
  dropdownItems?: DropdownItemProps[];
  isExternal?: boolean;
};

type DropdownItemProps = {
  name: string;
  href: string;
};

export function MenuItem({ name, icon, href, dropdownItems, isExternal }: MenuItemProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [ref, { height }] = useMeasure<HTMLUListElement>();
  const location = useLocation();

  const isChildrenActive =
    dropdownItems && dropdownItems.some((item) => item.href === location.pathname);

  useEffect(() => {
    if (isChildrenActive) {
      setIsOpen(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="mb-2 min-h-[8px] list-none last:mb-0">
      {dropdownItems?.length ? (
        <>
          <div
            onClick={() => setIsOpen(!isOpen)}
            role="button"
            aria-hidden="true"
            className={cn(
              'relative flex h-12 cursor-pointer items-center justify-between whitespace-nowrap  rounded-lg px-4 text-sm transition-all',
              isChildrenActive
                ? 'text-white'
                : 'text-gray-500 hover:text-brand dark:hover:text-white'
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
                        ? 'text-gray-500 before:bg-gray-500 hover:text-brand flex items-center rounded-lg p-3 text-sm transition-all before:h-1 before:w-1 before:rounded-full ltr:pl-6 before:ltr:mr-5 rtl:pr-6 before:rtl:ml-5 dark:hover:text-white'
                        : '!text-brand before:!bg-brand !font-medium before:-ml-0.5 before:!h-2 before:!w-2 before:ltr:!mr-[18px] before:rtl:!ml-[18px] dark:!text-white dark:before:!bg-white'
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
        <NavLink
          to={href}
          target={isExternal ? '_blank' : '_self'}
          rel="noopener noreferrer"
          className={({ isActive }) =>
            cn(
              `${isActive ? 'transition-all' : 'text-secondary-900'}`,
              'relative flex h-12 items-center whitespace-nowrap pl-8 text-sm text-secondary-900'
            )
          }
        >
          <span className="relative z-[1] mr-3">{icon}</span>
          <span className="relative z-[1] flex font-normal">
            {name} {isExternal && <ExternalIcon className="ml-1 h-3 w-3" />}
          </span>

          {href === location.pathname && (
            <motion.span
              className="absolute bottom-0 left-0 right-0 h-full w-full border-l-4 border-secondary-900 bg-secondary-100 shadow-1"
              layoutId="menu-item-active-indicator"
            />
          )}
        </NavLink>
      )}
    </div>
  );
}
