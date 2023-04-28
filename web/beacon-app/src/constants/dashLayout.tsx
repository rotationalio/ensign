import { t, Trans } from '@lingui/macro';

import { EXTRENAL_LINKS, PATH_DASHBOARD, ROUTES } from '@/application';
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
    name: t`Home`,
    icon: <HomeIcon />,
    href: PATH_DASHBOARD.ROOT,
  },
  {
    name: t`Projects`,
    icon: <FolderIcon />,
    href: PATH_DASHBOARD.PROJECTS,
  },
  {
    name: t`Team`,
    icon: <TeamIcon />,
    href: PATH_DASHBOARD.TEAMS,
  },
];

export const otherMenuItems: MenuItem[] = [
  {
    name: t`Docs`,
    icon: <DocsIcon />,
    href: ROUTES.DOCS,
    isExternal: true,
  },
  {
    name: t`Support`,
    icon: <SupportIcon />,
    href: ROUTES.SUPPORT,
    isExternal: true,
    isMail: true,
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
    href: EXTRENAL_LINKS.ABOUT,
    isExternal: true,
  },
  {
    name: t`Contact Us`,
    href: EXTRENAL_LINKS.CONTACT,
    isExternal: true,
  },
  {
    name: t`Server Status`,
    href: EXTRENAL_LINKS.SERVER,
    isExternal: true,
  },
  {
    name: <Trans>&copy; Rotational Labs, Inc</Trans>,
    href: EXTRENAL_LINKS.ROTATIONAL,
    isExternal: true,
  },
];
