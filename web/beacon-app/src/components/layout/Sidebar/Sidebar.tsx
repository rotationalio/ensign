import { Avatar, Button, useMenu } from '@rotational/beacon-core';
import { Link, useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import { ChevronDown } from '@/components/icons/chevron-down';
import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { Dropdown as Menu } from '@/components/ui/Dropdown';
import { footerItems, menuItems, otherMenuItems, SIDEBAR_WIDTH } from '@/constants/dashLayout';
import { useAuth } from '@/hooks/useAuth';
import { useOrgStore } from '@/store';

function SideBar() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'profile-menu' });
  const handleLogout = () => {
    logout();
    navigate('/');
  };
  const redirectToSettings = () => {
    navigate(PATH_DASHBOARD.ORGANIZATION);
  };

  const org = useOrgStore.getState() as any;
  return (
    <>
      <aside
        className={`fixed top-0 left-0 right-0  flex h-screen flex-col bg-[#F7F9FB] pt-5 pb-10`}
        style={{
          maxWidth: SIDEBAR_WIDTH,
        }}
      >
        <div className="flew-row flex w-full items-center gap-2 overflow-hidden py-2 pl-4 text-sm">
          <Avatar alt={org.name} src={org?.picture} className="flex" data-testid="avatar" />
          <h1 className="flex">
            {org?.name.split(' ')[0]}
            <br />
            {org?.name.split(' ')[1]}
          </h1>
          <div className="absolute right-5 flex">
            <Button variant="ghost" className="border-transparent border-none" onClick={open}>
              <ChevronDown />
            </Button>
          </div>
        </div>
        <div className="grow pt-8">
          <div>
            {menuItems.map((item, index) => (
              <MenuItem
                href={
                  item.href === PATH_DASHBOARD.PROJECTS
                    ? `${PATH_DASHBOARD.PROJECTS}/${org.projectID}`
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
      <div className="flex">
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          <Menu.Item onClick={handleLogout}>logout</Menu.Item>
          <Menu.Item onClick={redirectToSettings}>settings</Menu.Item>
        </Menu>
      </div>
    </>
  );
}

export default SideBar;
