import { Avatar, Loader, useMenu } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { Link, useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import { ChevronDown } from '@/components/icons/chevron-down';
import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { Dropdown as Menu } from '@/components/ui/Dropdown';
import { footerItems, menuItems, otherMenuItems, SIDEBAR_WIDTH } from '@/constants/dashLayout';
import { useFetchOrg } from '@/features/organization/hooks/useFetchOrgDetail';
import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useAuth } from '@/hooks/useAuth';
import { useOrgStore } from '@/store';
function SideBar() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const getOrg = useOrgStore.getState() as any;
  const { tenants } = useFetchTenants();
  const { projects, wasProjectsFetched } = useFetchTenantProjects(tenants?.tenants[0]?.id);
  const { org, isFetchingOrg } = useFetchOrg(getOrg?.org);
  if (wasProjectsFetched) {
    getOrg.setProjectID(projects?.tenant_projects[0]?.id);
  }
  if (org) {
    getOrg.setOrgName(org.name);
  }
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'profile-menu' });
  const handleLogout = () => {
    logout();
    navigate('/');
  };
  const redirectToSettings = () => {
    navigate(PATH_DASHBOARD.ORGANIZATION);
  };

  return (
    <>
      <aside
        className={`fixed top-0 left-0 right-0 flex h-screen flex-col bg-[#1D65A6] pt-5 pb-10 text-white`}
        style={{
          maxWidth: SIDEBAR_WIDTH,
        }}
      >
        <ErrorBoundary fallback={<div className="flex">Reload</div>}>
          <div
            onClick={open}
            role="button"
            tabIndex={0}
            aria-hidden="true"
            className="flex w-full flex-row items-center justify-between py-2 pr-5 pl-8 text-sm"
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
                {org?.name?.split(' ')[1]}
              </h1>
            </div>
            <div className="flex-end">
              <ChevronDown />
            </div>
          </div>
        </ErrorBoundary>
        <div className="grow pt-8">
          <div>
            {menuItems.map((item, index) => (
              <MenuItem
                href={
                  item.href === PATH_DASHBOARD.PROJECTS
                    ? `${PATH_DASHBOARD.PROJECTS}/${getOrg?.projectID}`
                    : item.href
                }
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
        <div className="ml-8 space-y-3">
          <ul className="space-y-1 text-xs text-white">
            {footerItems.map((item) => (
              <li key={item.name}>
                <Link to={item.href} target="_blank">
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
          <p className="text-xs text-white">&copy; Rotational Labs, Inc</p>
        </div>
      </aside>
      <div className="flex">
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          <Menu.Item onClick={redirectToSettings} data-testid="settings">
            Settings
          </Menu.Item>
          <Menu.Item onClick={handleLogout} data-testid="logoutButton">
            Logout
          </Menu.Item>
        </Menu>
      </div>
    </>
  );
}

export default SideBar;
