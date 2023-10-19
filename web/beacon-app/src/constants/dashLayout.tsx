import { t, Trans } from '@lingui/macro';
import { AiOutlineHome, AiOutlineProject, AiOutlineTeam } from 'react-icons/ai';
import { BsCodeSlash } from 'react-icons/bs';
import { CgProfile } from 'react-icons/cg';
import { HiOutlineDocument, HiOutlineLightBulb } from 'react-icons/hi';
import { IoSchool } from 'react-icons/io5';
import { MdOutlineContactSupport } from 'react-icons/md';
import { TbPlayFootball } from 'react-icons/tb';

import { EXTERNAL_LINKS, PATH_DASHBOARD } from '@/application';
import { MenuItem } from '@/types/MenuItem';

export const SIDEBAR_WIDTH = 250;
export const TOPBAR_HEIGHT = 60;

export const menuItems: MenuItem[] = [
  {
    name: t`Home`,
    icon: <AiOutlineHome fontSize={24} />,
    href: PATH_DASHBOARD.ROOT,
  },
  {
    name: t`Projects`,
    icon: <AiOutlineProject fontSize={24} />,
    href: PATH_DASHBOARD.PROJECTS,
    href_linked: PATH_DASHBOARD.TOPICS,
  },
  {
    name: t`Team`,
    icon: <AiOutlineTeam fontSize={24} />,
    href: PATH_DASHBOARD.TEAMS,
  },
  {
    name: t`Profile`,
    icon: <CgProfile fontSize={24} />,
    href: PATH_DASHBOARD.PROFILE,
    dropdownItems: [],
  },
];

export const otherMenuItems: MenuItem[] = [
  {
    name: t`Ensign U`,
    icon: <IoSchool fontSize={24} />,
    href: EXTERNAL_LINKS.ENSIGN_UNIVERSITY,
    isExternal: true,
  },
  {
    name: t`Use Cases`,
    icon: <HiOutlineLightBulb fontSize={24} />,
    href: EXTERNAL_LINKS.USE_CASES,
    isExternal: true,
  },
  {
    name: t`Docs`,
    icon: <HiOutlineDocument fontSize={24} />,
    href: EXTERNAL_LINKS.DOCS,
    isExternal: true,
  },
  {
    name: t`Data Playground`,
    icon: <TbPlayFootball fontSize={24} />,
    href: EXTERNAL_LINKS.DATA_PLAYGROUND,
    isExternal: true,
  },
  {
    name: t`SDKs`,
    icon: <BsCodeSlash fontSize={24} />,
    href: EXTERNAL_LINKS.SDK_DOCUMENTATION,
    isExternal: true,
  },

  {
    name: t`Support`,
    icon: <MdOutlineContactSupport fontSize={24} />,
    href: EXTERNAL_LINKS.SUPPORT,
    isExternal: true,
    isMail: true,
  },
];

export const footerItems = [
  {
    name: t`About`,
    href: EXTERNAL_LINKS.ABOUT,
    isExternal: true,
  },
  {
    name: t`Contact Us`,
    href: EXTERNAL_LINKS.CONTACT,
    isExternal: true,
  },
  {
    name: t`Server Status`,
    href: EXTERNAL_LINKS.SERVER,
    isExternal: true,
  },
  {
    name: <Trans>&copy; Rotational Labs, Inc</Trans>,
    href: EXTERNAL_LINKS.ROTATIONAL,
    isExternal: true,
  },
];
