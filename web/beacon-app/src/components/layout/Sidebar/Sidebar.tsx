import { Trans } from '@lingui/macro';
import { Avatar, Loader } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import cn from 'classnames';
import invariant from 'invariant';
import { useEffect, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { appConfig } from '@/application/config';
import ExternalIcon from '@/components/icons/external-icon';
import { OrganizationMenuDropdown } from '@/components/MenuDropdown/OrganizationMenuDropdown';
import { useDropdownMenu } from '@/components/MenuDropdown/useDropdownMenu';
import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { footerItems, menuItems, otherMenuItems } from '@/constants/dashLayout';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useFetchOrganizations } from '@/features/organization/hooks/useFetchOrganizations';
import { useFetchOrg } from '@/features/organization/hooks/useFetchOrgDetail';
import { useAuth } from '@/hooks/useAuth';
import { useOrgStore } from '@/store';

type SidebarProps = {
  className?: string;
};

function SideBar({ className }: SidebarProps) {
  const { profile: userInfo } = useFetchProfile();
  const { version: appVersion, revision: gitRevision } = appConfig;

  const navigate = useNavigate();
  const { logout } = useAuth();
  const appState = useOrgStore((state: any) => state) as any;
  const refreshOnceRef = useRef(false);
  const { org, isFetchingOrg, error, getOrgDetail } = useFetchOrg(appState?.orgID);
  const { organizations } = useFetchOrganizations();
  const [isOpen, setIsOpen] = useState(false);
  const { menuItems: dropdownItems } = useDropdownMenu({
    organizationsList: organizations?.organizations,
    currentOrg: appState?.orgID,
  });

  const onOpenChange = () => {
    setIsOpen(!isOpen);
  };

  const handleOpen = () => {
    setIsOpen(true);
  };

  useEffect(() => {
    if (appState?.orgID && !refreshOnceRef.current && error) {
      getOrgDetail();
      refreshOnceRef.current = true;
    }
  }, [appState?.orgID, getOrgDetail, error]);

  useEffect(() => {
    if (error?.status === 401) {
      // ('error?.status', error?.status);
      logout();
      navigate('/');
    }
  }, [error, logout, navigate]);
  // set the orgname in the store
  useEffect(() => {
    if (org?.name) {
      useOrgStore.setState({ orgName: org?.name });
    }
  }, [org]);

  // make sure we have the orgID
  useEffect(() => {
    invariant(appState?.orgID, 'orgID is not defined');
  }, [appState?.orgID]);

  return (
    <>
      <aside
        className={cn(
          `fixed top-0 left-0 flex h-screen flex-col bg-[#1D65A6] pt-5 pb-10 text-white md:w-[250px]`,
          className
        )}
      >
        <div className="flex h-full flex-col" data-cy="sidebar">
          <div className="grow">
            <ErrorBoundary
              fallback={
                <div className="flex">
                  <Trans>Something went wrong. Please try again later.</Trans>
                </div>
              }
            >
              <div
                onClick={handleOpen}
                role="button"
                tabIndex={0}
                aria-hidden="true"
                className="flex w-full flex-row items-center justify-between py-2 pr-5 pl-8 text-sm outline-none"
                data-testid="menu"
                data-cy="menu"
              >
                <div className="flex items-center gap-3 ">
                  <Avatar
                    alt={appState?.name || userInfo?.organization}
                    src={appState?.picture || userInfo?.picture}
                    className="flex w-64  "
                    data-testid="avatar"
                    data-cy="avatar"
                  />
                  <h1 className="flex" data-testid="orgName" data-cy="org-name">
                    {!org?.name && isFetchingOrg && <Loader className="flex" />}
                    {org?.name?.split(' ')[0]}
                    <br />
                    {org?.name?.split(' ').slice(1).join(' ')}
                  </h1>
                </div>
                <div className="flex-end">
                  {dropdownItems?.organizationMenuItems?.length > 0 && (
                    <OrganizationMenuDropdown
                      items={dropdownItems}
                      onOpenChange={onOpenChange}
                      isOpen={isOpen}
                      data-cy="org-menu"
                    />
                  )}
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
                    href_linked={item?.href_linked}
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
          <div className="ml-8 mt-5 space-y-3">
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
            <p className="w-full text-xs text-white">
              {appVersion && <span>App Version {appVersion} </span>}
              {gitRevision && <span>& Git Revision {gitRevision} </span>}
            </p>
          </div>
        </div>
      </aside>
    </>
  );
}

export default SideBar;
