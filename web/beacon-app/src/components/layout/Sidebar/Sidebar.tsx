import { Link } from 'react-router-dom';

import { routes } from '@/application';
import DocsIcon from '@/components/icons/docs';
import FolderIcon from '@/components/icons/folder';
import HomeIcon from '@/components/icons/home-icon';
import ProfileIcon from '@/components/icons/profile';
import SupportIcon from '@/components/icons/support';
import Avatar from '@/components/ui/Avatar';
import { MenuItem } from '@/components/ui/CollapsibleMenu';

type MenuItem = {
  name: string;
  icon: JSX.Element;
  href: string;
  isExternal?: boolean;
  dropdownItems?: Pick<MenuItem, 'name' | 'href'>[];
};

const menuItems: MenuItem[] = [
  {
    name: 'Home',
    icon: <HomeIcon />,
    href: routes.home,
  },
  {
    name: 'Projects',
    icon: <FolderIcon />,
    href: routes.projects,
  },
];

const otherMenuItems: MenuItem[] = [
  {
    name: 'Docs',
    icon: <DocsIcon />,
    href: routes.docs,
    isExternal: false,
  },
  {
    name: 'Support',
    icon: <SupportIcon />,
    href: routes.support,
    isExternal: true,
  },
  {
    name: 'Profile',
    icon: <ProfileIcon />,
    href: routes.profile,
    dropdownItems: [],
  },
];

const footerItems = [
  {
    name: 'About',
    href: '/#',
  },
  {
    name: 'Contact Us',
    href: '/#',
    isExternal: true,
  },
  {
    name: 'Server Status',
    href: '/#',
    isExternal: true,
  },
];

function SideBar() {
  return (
    <aside
      className={
        'xs:w-80 top-0 left-0 right-0 z-40 flex h-screen w-full max-w-full flex-col bg-[#F7F9FB] px-5 pt-5 pb-10 md:w-72 xl:fixed 2xl:w-80'
      }
    >
      <div className="relative flex items-center gap-2 overflow-hidden py-2 pl-4 text-sm">
        <Avatar alt="Acme Systems" />
        <h1>
          Acme <br /> Systems
        </h1>
      </div>
      <div className="grow pt-8">
        <div>
          {menuItems.map((item, index) => (
            <MenuItem
              href={item.href}
              key={'default' + item.name + index}
              name={item.name}
              icon={item.icon}
              dropdownItems={item?.dropdownItems}
              isExternal={item.isExternal}
            />
          ))}
        </div>
        <hr className="my-5 mx-8"></hr>
        <div>
          {otherMenuItems.map((item, index) => (
            <MenuItem
              href={item.href}
              key={'default' + item.name + index}
              name={item.name}
              icon={item.icon}
              dropdownItems={item?.dropdownItems}
              isExternal={item.isExternal}
            />
          ))}
        </div>
      </div>
      <div className="ml-8 space-y-3">
        <ul className="space-y-1 text-xs text-neutral-600">
          {footerItems.map((item) => (
            <li key={item.name}>
              <Link to={item.href}>{item.name}</Link>
            </li>
          ))}
        </ul>
        <p className="text-xs text-neutral-600">&copy; Rotational Labs, Inc</p>
      </div>
    </aside>
  );
}

export default SideBar;
