import { t } from '@lingui/macro';

import { PATH_DASHBOARD, ROUTES } from '@/application';
import DocsIcon from '@/components/icons/docs';
import FolderIcon from '@/components/icons/folder';
import HomeIcon from '@/components/icons/home-icon';
import ProfileIcon from '@/components/icons/profile';
import SupportIcon from '@/components/icons/support';
import { MenuItem } from '@/types/MenuItem';

export const SIDEBAR_WIDTH = 250;
export const TOPBAR_HEIGHT = 60;

export const menuItems: MenuItem[] = [
  {
    name: t`Home`,
    icon: <HomeIcon />,
    href: PATH_DASHBOARD.ROOT,
  },
  {
    name: t`Projects`,
    icon: <FolderIcon />,
    href: PATH_DASHBOARD.PROJECTS,
  },
];

export const otherMenuItems: MenuItem[] = [
  {
    name: t`Docs`,
    icon: <DocsIcon />,
    href: ROUTES.DOCS,
    isExternal: false,
  },
  {
    name: t`Support`,
    icon: <SupportIcon />,
    href: ROUTES.SUPPORT,
    isExternal: true,
  },
  {
    name: t`Profile`,
    icon: <ProfileIcon />,
    href: PATH_DASHBOARD.PROFILE,
    dropdownItems: [],
  },
];

export const footerItems = [
  {
    name: t`About`,
    href: '/#',
  },
  {
    name: t`Contact Us`,
    href: '/#',
    isExternal: true,
  },
  {
    name: `Server Status`,
    href: '/#',
    isExternal: true,
  },
];
