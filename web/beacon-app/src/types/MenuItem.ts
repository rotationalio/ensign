export type MenuItem = {
  name: string;
  icon: JSX.Element;
  href: string;
  isExternal?: boolean;
  isMail?: boolean;
  dropdownItems?: Pick<MenuItem, 'name' | 'href'>[];
};
