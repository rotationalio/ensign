import { Avatar, Loader } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import cn from 'classnames';
import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { appConfig } from '@/application/config';
import ExternalIcon from '@/components/icons/external-icon';
import { MenuDropdownMenu } from '@/components/MenuDropdown/MenuDropdown';
import { useDropdownMenu } from '@/components/MenuDropdown/useDropdownMenu';
import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { footerItems, menuItems, otherMenuItems } from '@/constants/dashLayout';
import { useFetchOrganizations } from '@/features/organization/hooks/useFetchOrganizations';
import { useFetchOrg } from '@/features/organization/hooks/useFetchOrgDetail';
import { useAuth } from '@/hooks/useAuth';
import { useOrgStore } from '@/store';

type SidebarProps = {
  className?: string;
};

function SideBar({ className }: SidebarProps) {
  const { beaconVersion, apiVersion } = appConfig;
  const navigate = useNavigate();
  const { logout } = useAuth();
  const getOrg = useOrgStore.getState() as any;
  const { org, isFetchingOrg, error } = useFetchOrg(getOrg?.org);
  const { organizations } = useFetchOrganizations();
  const [isOpen, setIsOpen] = useState(false);
  const { menuItems: dropdownItems } = useDropdownMenu({
    organizationsList: organizations?.organizations,
    currentOrg: getOrg?.org,
  });

  if (org) {
    getOrg.setOrgName(org.name);
  }

  const onOpenChange = () => {
    setIsOpen(!isOpen);
  };

  const handleOpen = () => {
    setIsOpen(true);
  };

  useEffect(() => {
    //console.log('getOrg?.name', getOrg?.name);
    if (!getOrg?.name) {
      logout();
      navigate('/');
    }
  }, [getOrg, logout, navigate]);

  useEffect(() => {
    if (error?.status === 401) {
      console.log('error?.status', error?.status);
      logout();
      navigate('/');
    }
  }, [error, logout, navigate]);

  return (
    <>
      <aside
        className={cn(
          `fixed top-0 left-0 flex h-screen flex-col bg-[#1D65A6] pt-5 pb-10 text-white md:w-[250px]`,
          className
        )}
      >
        <div className="flex h-full flex-col">
          <div className="grow">
            <ErrorBoundary fallback={<div className="flex">Reload</div>}>
              <div
                onClick={handleOpen}
                role="button"
                tabIndex={0}
                aria-hidden="true"
                className="flex w-full flex-row items-center justify-between py-2 pr-5 pl-8 text-sm outline-none"
                data-testid="menu"
              >
                <div className="flex items-center gap-3 ">
                  <Avatar
                    alt={getOrg?.name}
                    src={getOrg?.picture}
                    className="flex w-64  "
                    data-testid="avatar"
                  />
                  <h1 className="flex" data-testid="orgName">
                    {!org?.name && isFetchingOrg && <Loader className="flex" />}
                    {org?.name?.split(' ')[0]}
                    <br />
                    {org?.name?.split(' ').slice(1).join(' ')}
                  </h1>
                </div>
                <div className="flex-end">
                  <MenuDropdownMenu
                    items={dropdownItems}
                    onOpenChange={onOpenChange}
                    isOpen={isOpen}
                  />
                </div>
              </div>
            </ErrorBoundary>
            <div className="pt-8">
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
                    isMail={item.isMail}
                  />
                ))}
              </div>
            </div>
          </div>
          <div className="ml-8 space-y-3">
            <ul className="space-y-1 text-xs text-white">
              {footerItems.map((item) => (
                <li key={`${item.name}`}>
                  <Link to={item.href} target="_blank" className="flex">
                    {item.name}{' '}
                    {item.isExternal && <ExternalIcon className="ml-1 h-3 w-3 text-white" />}
                  </Link>
                </li>
              ))}
            </ul>
            <p>
              {beaconVersion && <span className="text-xs text-white">Beacon {beaconVersion} </span>}
              {apiVersion && <span className="text-xs text-white">& API {apiVersion} </span>}
            </p>
          </div>
        </div>
      </aside>
    </>
  );
}

export default SideBar;
