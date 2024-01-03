import { Avatar, Heading, Loader } from '@rotational/beacon-core';
import cn from 'classnames';
import { Link } from 'react-router-dom';

import { appConfig } from '@/application/config';
import ExternalIcon from '@/components/icons/external-icon';
import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { footerItems, menuItems, otherMenuItems } from '@/constants/dashLayout';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useFetchOrg } from '@/features/organization/hooks/useFetchOrgDetail';
import { useOrgStore } from '@/store';

type SandboxSidebarProps = {
  className?: string;
};

function SandboxSidebar({ className }: SandboxSidebarProps) {
  const { version: appVersion, revision: gitRevision } = appConfig;
  const { profile: userInfo } = useFetchProfile();
  const appState = useOrgStore((state: any) => state) as any;
  const { org, isFetchingOrg } = useFetchOrg(appState?.orgID);

  return (
    <>
      <aside
        className={cn(
          `fixed left-0 top-0 flex h-screen flex-col bg-[#1D65A6] pb-10 pt-5 text-white md:w-[250px]`,
          className
        )}
      >
        <div
          className="flex h-full flex-col"
          data-testid="sandbox-sidebar"
          data-cy="sandbox-sidebar"
        >
          <div className="grow">
            <Heading as="h1" className="ml-8 space-y-3 text-lg font-bold">
              Ensign Sandbox
            </Heading>
            <div className="ml-8 flex items-center gap-3 pt-6 text-sm">
              <Avatar
                alt={appState?.name || userInfo?.organization}
                src={appState?.picture || userInfo?.picture}
                className="w-64"
                data-testid="sandbox-avatar"
                data-cy="sandbox-avatar"
              />
              <Heading as="h2" data-testid="org-name" data-cy="org-name">
                {!org?.name && isFetchingOrg && <Loader className="flex" />}
                {org?.name?.split(' ')[0]}
                <br />
                {org?.name?.split(' ').slice(1).join(' ')}
              </Heading>
            </div>
            <div className="pt-8">
              {menuItems.map((item, index) => (
                <MenuItem
                  href={item?.href}
                  key={item?.name + index}
                  name={item?.name}
                  icon={item?.icon}
                  href_linked={item?.href_linked}
                  dropdownItems={item?.dropdownItems}
                  isExternal={item?.isExternal}
                  isMail={item?.isMail}
                />
              ))}
              <div className="mx-8 my-5 border-b border-b-white"></div>
              {otherMenuItems.map((item, index) => (
                <MenuItem
                  href={item?.href}
                  key={item?.name + index}
                  name={item?.name}
                  icon={item?.icon}
                  href_linked={item?.href_linked}
                  dropdownItems={item?.dropdownItems}
                  isExternal={item?.isExternal}
                  isMail={item?.isMail}
                />
              ))}
            </div>
          </div>
          <div className="ml-8 mt-5 space-y-3">
            <ul className="space-y-1 text-xs text-white">
              {footerItems.map((item) => (
                <li key={`${item?.name}`}>
                  <Link to={item?.href} target="_blank" className="flex">
                    {item?.name}
                    {item?.isExternal && <ExternalIcon className="ml-1 h-3 w-3 text-white" />}
                  </Link>
                </li>
              ))}
            </ul>
            <section className="w-full text-xs text-white">
              {appVersion && <span>App Version {appVersion} & </span>}
              {gitRevision && <span>Git Revision {gitRevision}</span>}
            </section>
          </div>
        </div>
      </aside>
    </>
  );
}

export default SandboxSidebar;
