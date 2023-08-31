import chroma from 'chroma-js';
import ReactSelect, { Props } from 'react-select';
interface SelectProps extends Props {
  isMulti?: boolean;
  isDisabled?: boolean;
}

function Select(props: SelectProps) {
  return (
    <ReactSelect
      {...(props.isMulti && { isMulti: true })}
      components={{
        IndicatorSeparator: () => null,
      }}
      menuPosition={'fixed'}
      styles={{
        container: (base) => ({
          ...base,
          position: 'relative',
          width: '100%',
        }),
        control: (base) => ({
          ...base,
          fontSize: 14,
          padding: 5,
          borderRadius: '0.375rem',
          borderColor: '#000',
          '&:hover': {
            borderColor: '#000',
          },
        }),
        placeholder: (base) => ({
          ...base,
          color: 'gray',
        }),
        menu: (base) => ({
          ...base,
          zIndex: 9999999,
          fontSize: 14,
        }),
        menuPortal: (base) => ({
          ...base,
          zIndex: 9999999,
        }),
        multiValue: (base) => {
          const color = chroma('#C5EDFF').alpha(0.5).css();
          return {
            ...base,
            backgroundColor: color,
            borderRadius: '0.375rem',
            color: '#000',
            fontSize: 14,
            padding: 5,
          };
        },
        multiValueRemove: (base) => ({
          ...base,
          color: '#000',
          ':hover': {
            backgroundColor: '#C5EDFF',
            color: '#000',
          },
        }),
        // option: (base, { isFocused }) => {
        //   const color = chroma('#C5EDFF').alpha(0.5).css();
        //   return {
        //     ...base,
        //     backgroundColor: isFocused ? color : '#fff',
        //     color: '#000',
        //     fontSize: 14,
        //     padding: 5,
        //     '&:hover': {
        //       backgroundColor: color,
        //       color: '#000',
        //     },
        //   };
        // },
      }}
      {...props}
    />
  );
}

export default Select;
