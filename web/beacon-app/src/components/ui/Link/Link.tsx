import { mergeClassnames } from '@rotational/beacon-core';
import React from 'react';
interface LinkProps {
  href?: string;
  children: React.ReactNode;
  className?: string;
  openInNewTab?: boolean;
  onClick?: () => void;
}

const renderLink = (
  href: string,
  combinedClasses: string,
  combinedProps: any,
  children: React.ReactNode
) => {
  return (
    <a href={href} className={combinedClasses} {...combinedProps}>
      {children}
    </a>
  );
};

const renderButton = (combinedClasses: string, combinedProps: any, children: React.ReactNode) => {
  return (
    <button type="button" className={combinedClasses} {...combinedProps}>
      {children}
    </button>
  );
};

const Link: React.FC<LinkProps> = ({ href, children, className, openInNewTab, onClick }) => {
  const defaultClasses = 'text-[#1F4CED] hover:underline';
  const buttonClasses = 'hover:cursor-pointer text-[#1F4CED] hover:underline';
  const targetProps = openInNewTab ? { target: '_blank', rel: 'noopener noreferrer' } : {};

  const clickProps = onClick ? { onClick } : {};

  const combinedClasses = mergeClassnames(onClick ? buttonClasses : defaultClasses, className);
  const combinedProps = { ...targetProps, ...clickProps };

  return (
    <>
      {href && renderLink(href, combinedClasses, combinedProps, children)}
      {!href && onClick && renderButton(combinedClasses, combinedProps, children)}
    </>
  );
};

export default Link;
