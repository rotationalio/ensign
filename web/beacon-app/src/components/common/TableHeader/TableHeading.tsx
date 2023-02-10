import { Heading } from '@rotational/beacon-core';

export interface TableHeadingProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

const TableHeading = ({ children, ...rest }: TableHeadingProps) => {
  return (
    <div className="flex w-full bg-[#F7F9FB] p-2" {...rest}>
      <Heading as={'h2'} className="text-black text-sm font-semibold">
        {children}
      </Heading>
    </div>
  );
};

export default TableHeading;
