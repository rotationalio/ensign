import React from 'react';
import { useBreadcrumbItem, useBreadcrumbs } from 'react-aria';
import { twMerge } from 'tailwind-merge';

import { BreadcrumbItemProps, BreadcrumbsChildProps, BreadcrumbsProps } from './Breadcrumb.type';

function Breadcrumbs(props: BreadcrumbsProps) {
  const { navProps } = useBreadcrumbs(props);
  const children = React.Children.toArray(props.children);

  return (
    <nav {...navProps} className={twMerge('flex items-center gap-2', props.className)}>
      <ol style={{ display: 'flex', listStyle: 'none', margin: 0, padding: 0 }}>
        {children.map(
          (child, index) =>
            React.isValidElement(child) &&
            React.cloneElement(child as React.ReactElement<BreadcrumbsChildProps>, {
              isCurrent: index === children.length - 1,
              separator: props.separator,
            })
        )}
      </ol>
    </nav>
  );
}

function BreadcrumbItem(props: BreadcrumbItemProps) {
  const { href, separator = '›', children } = props;
  const ref = React.useRef(null);
  const { itemProps } = useBreadcrumbItem(props, ref);

  return (
    <li className="text-sm">
      <a
        {...itemProps}
        ref={ref}
        href={href}
        style={{
          color: props.isDisabled || !props.isCurrent ? 'gray' : 'black',
          cursor: props.isCurrent || props.isDisabled ? 'default' : 'pointer',
        }}
        className={props.className}
      >
        {children}
      </a>
      {!props.isCurrent && (
        <span aria-hidden="true" className="px-3">
          {separator}
        </span>
      )}
    </li>
  );
}

Breadcrumbs.Item = BreadcrumbItem;

export default Breadcrumbs;
