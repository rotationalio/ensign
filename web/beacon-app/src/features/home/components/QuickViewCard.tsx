import { ReactNode } from 'react';

type QuickViewCardProps = {
  title: string;
  color: string;
  children: ReactNode;
};

function QuickViewCard({ title, color, children }: QuickViewCardProps) {
  return (
    <div
      style={{ backgroundColor: color }}
      className="flex h-20 w-full flex-col justify-between rounded-xl px-4 py-3"
    >
      <h5 className="text-xs font-semibold">{title}</h5>
      <p className="font-semibold">{children}</p>
    </div>
  );
}

export default QuickViewCard;
