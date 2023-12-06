export interface CheckboxProps {
  onClick?: () => void;
  className?: string;
  id?: string;
  label: string;
  dataCy?: string;
}

function Checkbox({ onClick, className, id, label, dataCy }: CheckboxProps) {
  return (
    <div className={className}>
      <input
        type="checkbox"
        id={id}
        className="border-2 border-gray-600"
        onChange={onClick}
        data-cy={dataCy}
      />
      <label htmlFor={id} className="ml-2">
        {label}
      </label>
    </div>
  );
}

export default Checkbox;
