export interface CheckboxProps {
  onClick: () => void;
  value: string;
}

function Checkbox({ value, onClick }: CheckboxProps) {
  return (
    <>
      <label>
        <input type="checkbox" onChange={onClick} />
        <span className="ml-2">{value}</span>
      </label>
    </>
  );
}

export default Checkbox;
