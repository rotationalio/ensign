import { ReactNode } from 'react';
import { AriaBreadcrumbItemProps, AriaBreadcrumbsProps } from 'react-aria';

export type BreadcrumbsProps = {
  children: React.ReactNode;
  separator?: ReactNode;
  onClick?: (e: React.MouseEventHandler<HTMLAnchorElement>) => void;
} & AriaBreadcrumbsProps &
  React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;

export type BreadcrumbItemProps = {
  href?: string;
  separator?: ReactNode;
  onClick?: (e: React.MouseEventHandler<HTMLAnchorElement>) => void;
} & AriaBreadcrumbItemProps &
  React.DetailedHTMLProps<React.AnchorHTMLAttributes<HTMLAnchorElement>, HTMLAnchorElement>;

export type BreadcrumbsChildProps = {
  isCurrent: boolean;
  separator?: ReactNode;
};
