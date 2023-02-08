export type MenuItem = {
  name: string;
  icon: JSX.Element;
  href: string;
  isExternal?: boolean;
  dropdownItems?: Pick<MenuItem, 'name' | 'href'>[];
};
