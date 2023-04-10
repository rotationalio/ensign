import ReactSelect, { Props } from 'react-select';

function Select(props: Props) {
  return (
    <ReactSelect
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
      }}
      {...props}
    />
  );
}

export default Select;
