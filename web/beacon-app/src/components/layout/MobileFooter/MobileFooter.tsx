import { Link } from 'react-router-dom';

import { footerItems } from '@/constants/dashLayout';

function MobileFooter() {
  return (
    <footer className="fixed bottom-0 left-0 w-screen bg-[#1D65A6] py-1 md:hidden md:pl-[250px]">
      <ul className="flex items-center justify-around space-y-1 text-xs text-white">
        {footerItems.map((item) => (
          <li key={item.name}>
            <Link to={item.href} target="_blank">
              {item.name}
            </Link>
          </li>
        ))}
      </ul>
    </footer>
  );
}

export default MobileFooter;
