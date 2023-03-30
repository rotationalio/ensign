import { FOOTER, PATH_DASHBOARD, ROUTES } from '@/application';
import DocsIcon from '@/components/icons/docs';
import FolderIcon from '@/components/icons/folder';
import HomeIcon from '@/components/icons/home-icon';
import ProfileIcon from '@/components/icons/profile';
import SupportIcon from '@/components/icons/support';
import TeamIcon from '@/components/icons/team';
import { MenuItem } from '@/types/MenuItem';

export const SIDEBAR_WIDTH = 250;
export const TOPBAR_HEIGHT = 60;

export const menuItems: MenuItem[] = [
  {
    name: 'Home',
    icon: <HomeIcon />,
    href: PATH_DASHBOARD.ROOT,
  },
  {
    name: 'Projects',
    icon: <FolderIcon />,
    href: PATH_DASHBOARD.PROJECTS,
  },
  {
    name: 'Teams',
    icon: <TeamIcon />,
    href: PATH_DASHBOARD.TEAMS,
  },
];

export const otherMenuItems: MenuItem[] = [
  {
    name: 'Docs',
    icon: <DocsIcon />,
    href: ROUTES.DOCS,
    isExternal: true,
  },
  {
    name: 'Support',
    icon: <SupportIcon />,
    href: ROUTES.SUPPORT,
    isExternal: true,
    isMail: true,
  },
  {
    name: 'Profile',
    icon: <ProfileIcon />,
    href: PATH_DASHBOARD.PROFILE,
    dropdownItems: [],
  },
];

export const footerItems = [
  {
    name: 'About',
    href: FOOTER.ABOUT,
    isExternal: true,
  },
  {
    name: 'Contact Us',
    href: FOOTER.CONTACT,
    isExternal: true,
  },
  {
    name: 'Server Status',
    href: FOOTER.SERVER,
    isExternal: true,
  },
  {
    name: <>&copy; Rotational Labs, Inc</>,
    href: 'https://rotational.io/',
    isExternal: true,
  },
];
