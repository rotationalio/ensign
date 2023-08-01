export type MenuItem = {
  name: string;
  icon: JSX.Element;
  href: string;
  href_linked?: string;
  isExternal?: boolean;
  isMail?: boolean;
  dropdownItems?: Pick<MenuItem, 'name' | 'href'>[];
};
