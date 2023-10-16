import { Listbox } from '@headlessui/react';
import { FC, useEffect, useState } from 'react';

import { SelectOptionType, SelectProps } from './Select.types';

type selectedOptionType = SelectOptionType | undefined;

const Select: FC<SelectProps<SelectOptionType[], SelectOptionType>> = ({
  options,
  value,
  onChange,
  placeholder,
  label,
  isError,
  hintText,
  isMulti = false,
  disabled,
  size = 'md',
  optionWidth = '100%',
  ...props
}) => {
  const [option, setOption] = useState<selectedOptionType>(value);

  const handleChange = (option: SelectOptionType) => {
    if (onChange) {
      onChange(option);
    }
    if (isMulti) {
      // concat the new option to the existing options with , as a separator
      // take the option value split and add the new string to the array
      setOption([...option, option.value.split(',')]);
    }
    setOption(option);
  };

  useEffect(() => {
    if (value) {
      handleChange(value);
    }
  }, [value]);

  return (
    <Listbox value={option} onChange={handleChange} {...(isMulti && { multiple: true })} {...props}>
      {({ open }) => (
        <>
          {label && (
            <Listbox.Label className="block text-sm font-medium text-gray-700">
              {label}
            </Listbox.Label>
          )}
          <div className="relative mt-1">
            <Listbox.Button className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 relative w-full cursor-default rounded-md border border-gray-300 bg-white py-2 pl-3 pr-10 text-left focus:outline-none focus:ring-1 sm:text-sm">
              <span className="block truncate">
                {isMulti &&
                  Array.isArray(option) &&
                  option.length > 0 &&
                  option.map((item) => item.label).join(', ')}
                {(!isMulti && !Array.isArray(option) && option?.label) || placeholder}
              </span>
              <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <svg
                  className="h-5 w-5 text-gray-400"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fillRule="evenodd"
                    d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                    clipRule="evenodd"
                  />
                </svg>
              </span>
            </Listbox.Button>
            <Listbox.Options className="shadow-lg absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
              {options.map((option) => (
                <Listbox.Option
                  key={option.value}
                  className={({ active }) =>
                    `${active ? 'bg-indigo-600 text-white' : 'text-gray-900'}
                                            relative cursor-default select-none py-2 pl-3 pr-9`
                  }
                  value={option}
                >
                  {({ selected, active }) => (
                    <>
                      <div className="flex items-center">
                        <span
                          className={`${selected ? 'font-semibold' : 'font-normal'} block truncate`}
                        >
                          {option.label}
                        </span>
                      </div>
                      {selected ? (
                        <span
                          className={`${active ? 'text-white' : 'text-indigo-600'}
                                                            absolute inset-y-0 right-0 flex items-center pr-4`}
                        >
                          <svg
                            className="h-5 w-5"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            aria-hidden="true"
                          >
                            <path
                              fillRule="evenodd"
                              d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                              clipRule="evenodd"
                            />
                          </svg>
                        </span>
                      ) : null}
                    </>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </div>
        </>
      )}
    </Listbox>
  );
};

Select.displayName = 'Select';

export default Select;
