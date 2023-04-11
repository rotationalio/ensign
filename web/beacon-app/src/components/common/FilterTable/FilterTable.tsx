import ArrowDownUp from '@/components/icons/arrow-down-up';
import FunnelSimple from '@/components/icons/funnel-simple';
import ThreeDots from '@/components/icons/three-dots';
import Union from '@/components/icons/union';
const FilterTable = () => {
  return (
    <ul className="flex items-center gap-3">
      <li className="flex items-center justify-center">
        <button>
          <Union className="fill-[#1C1C1C]" />
        </button>
      </li>
      <li className="flex items-center justify-center">
        <button>
          <FunnelSimple />
        </button>
      </li>
      <li className="flex items-center justify-center">
        <button>
          <ArrowDownUp />
        </button>
      </li>
      <li className="flex items-center justify-center">
        <button>
          <ThreeDots />
        </button>
      </li>
    </ul>
  );
};

export default FilterTable;
