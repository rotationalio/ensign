import React, { ReactNode } from 'react';
import {
  AriaBreadcrumbItemProps,
  AriaBreadcrumbsProps,
  useBreadcrumbItem,
  useBreadcrumbs,
} from 'react-aria';
import { twMerge } from 'tailwind-merge';

type BreadcrumbsChilProps = {
  isCurrent: boolean;
  separator?: ReactNode;
};

type BreadcrumbsProps = {
  children: React.ReactNode;
  separator?: ReactNode;
  onClick?: (e: React.MouseEventHandler<HTMLAnchorElement>) => void;
} & AriaBreadcrumbsProps &
  React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;

function Breadcrumbs(props: BreadcrumbsProps) {
  const { navProps } = useBreadcrumbs(props);
  const children = React.Children.toArray(props.children);

  return (
    <nav {...navProps} className={twMerge('flex items-center gap-2', props.className)}>
      <ol style={{ display: 'flex', listStyle: 'none', margin: 0, padding: 0 }}>
        {children.map(
          (child, index) =>
            React.isValidElement(child) &&
            React.cloneElement(child as React.ReactElement<BreadcrumbsChilProps>, {
              isCurrent: index === children.length - 1,
              separator: props.separator,
            })
        )}
      </ol>
    </nav>
  );
}

type BreadcrumbItemProps = {
  href?: string;
  separator?: ReactNode;
  onClick?: (e: React.MouseEventHandler<HTMLAnchorElement>) => void;
} & AriaBreadcrumbItemProps &
  React.DetailedHTMLProps<React.AnchorHTMLAttributes<HTMLAnchorElement>, HTMLAnchorElement>;

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
