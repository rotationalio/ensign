import { MixerHorizontalIcon } from '@radix-ui/react-icons';
import React, { ReactNode, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
// import { switchOrganizationRequest } from '@/features/organization/api/switchOrganizationApi';
import { useSwitchOrganization } from '@/features/organization/hooks/useSwitchOrganization';
import { useAuth } from '@/hooks/useAuth';
import { useOrgStore } from '@/store';
import { decodeToken } from '@/utils/decodeToken';

interface DropdownMenuPrimitiveProps {
  organizationsList: any[];
  currentOrg: string;
}

export interface MenuItem {
  label: string;
  shortcut?: string;
  icon?: ReactNode;
  onClick: () => void;
}

export interface Org {
  name: string;
  profileUrl?: string;
  id: string;
  handleSwitch: (orgId: string) => void | undefined;
}

const useDropdownMenu = ({ organizationsList, currentOrg }: DropdownMenuPrimitiveProps) => {
  const navigate = useNavigate();
  const { switch: switchOrganization, wasSwitchFetched, auth } = useSwitchOrganization();
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  function handleSwitch(orgId: string) {
    return () => {
      switchOrganization(orgId);
    };
  }

  useEffect(() => {
    if (wasSwitchFetched) {
      // persist org state
      const token = decodeToken(auth.access_token) as any;
      useOrgStore.persist.clearStorage();
      useOrgStore.setState({
        org: token?.org,
        user: token?.sub,
        isAuthenticated: !!wasSwitchFetched,
        name: token?.name,
        email: token?.email,
        picture: token?.picture,
        permissions: token?.permissions,
      });

      // reload the page

      window.location.reload();
    }
  }, [wasSwitchFetched, auth?.access_token]);

  const generalMenuItems: MenuItem[] = [
    {
      label: 'Settings',
      icon: <MixerHorizontalIcon className="h-3.5 w-3.5 mr-2" />,
      onClick: () => navigate(PATH_DASHBOARD.PROFILE),
    },
  ];

  const logoutMenuItem: MenuItem = {
    label: 'Logout',
    onClick: handleLogout,
  };

  const organizations = organizationsList?.filter((org: Org) => org.id !== currentOrg);

  const organizationMenuItems = organizations?.map((org: Org) => ({
    name: org.name,
    orgId: org.id,
    handleSwitch: handleSwitch(org.id) as any,
  }));

  const menuItems = {
    generalMenuItems,
    organizationMenuItems,
    logoutMenuItem,
  };

  return {
    menuItems,
  };
};

export { useDropdownMenu };
