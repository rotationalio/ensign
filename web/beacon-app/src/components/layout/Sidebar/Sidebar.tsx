import { Link } from 'react-router-dom';
import { Avatar } from '@rotational/beacon-core';

import { MenuItem } from '@/components/ui/CollapsibleMenu';
import { footerItems, menuItems, otherMenuItems, SIDEBAR_WIDTH } from '@/constants/dash-layout';

function SideBar() {
  return (
    <aside
      className={`fixed top-0 left-0 right-0 z-40 flex h-screen flex-col bg-[#F7F9FB] pt-5 pb-10`}
      style={{
        maxWidth: SIDEBAR_WIDTH,
      }}
    >
      <div className="relative flex items-center gap-2 overflow-hidden py-2 pl-4 text-sm">
        <Avatar alt="Acme Systems" />
        <h1>
          Acme <br /> Systems
        </h1>
      </div>
      <div className="grow pt-8">
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
  );
}

export default SideBar;
