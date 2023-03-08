import { Link, useLocation } from 'react-router-dom';
import styled from 'styled-components';

import { menuItems, otherMenuItems } from '@/constants/dashLayout';

const StyledLink = styled<any>(Link)`
  position: relative;
  &:after {
    content: ' ';
    position: absolute;
    width: 100;
    border-radius: 5px;
    bottom: 0;
  }
  svg {
    fill: ${(props: any) => (props.isActive ? '#fff' : '')};
  }
`;

function MobileNav() {
  const location = useLocation();

  return (
    <div className="relative flex items-center gap-3 px-4 text-white md:hidden">
      {menuItems.map((item) => (
        <StyledLink to={item.href} key={item.name} isActive={location.pathname === item.href}>
          {item.icon}
        </StyledLink>
      ))}
      {otherMenuItems.map((item) => (
        <StyledLink to={item.href} key={item.name} isActive={location.pathname === item.href}>
          {item.icon}
        </StyledLink>
      ))}
    </div>
  );
}

export default MobileNav;
