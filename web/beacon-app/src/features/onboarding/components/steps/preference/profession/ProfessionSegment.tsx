import CheckCircleIcon from '@/components/icons/check-circle';

import { getProfessionOptions } from '../../../../shared/utils';
import Header from './Header';
type ProfessionSegmentProps = {
  onChange?: (value: any) => void;
  selectedValue?: string;
};
const ProfessionSegment = ({ onChange, selectedValue }: ProfessionSegmentProps) => {
  const PROFESSION_OPTIONS = getProfessionOptions();
  return (
    <div>
      <Header />
      <div className="my-5">
        <ul className="grid w-full gap-20 md:grid-cols-3" data-cy="profession-segment">
          {PROFESSION_OPTIONS?.map((option: any, idx: any) => (
            <li key={idx}>
              <input
                onChange={() => onChange && onChange(option)}
                id={option.id}
                type="radio"
                value={option.value}
                name="profession_segment"
                className="peer hidden"
                data-cy={`profession-${option.value}`}
                // checked if value is equal to selected value
                checked={option.value === selectedValue}
              />
              <label
                htmlFor={option.id}
                className="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white p-5 text-gray-500 hover:bg-gray-100 hover:text-gray-600 peer-checked:border-blue-600 peer-checked:text-blue-600"
              >
                <div className="block">
                  <div className="w-full text-lg font-bold">{option.label}</div>
                </div>
                <CheckCircleIcon />
              </label>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default ProfessionSegment;
