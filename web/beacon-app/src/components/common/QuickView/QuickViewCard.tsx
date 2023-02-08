import { Loader } from '@rotational/beacon-core';
import { ReactNode, Suspense } from 'react';
type QuickViewCardProps = {
  title: string;
  color: string;
  children: ReactNode;
};

function QuickViewCard({ title, color, children }: QuickViewCardProps) {
  return (
    <div
      style={{ backgroundColor: color }}
      className="flex h-20 w-full flex-col justify-between rounded-xl py-3 px-4"
    >
      <Suspense
        fallback={
          <div className="justify-center text-center">
            <Loader />
          </div>
        }
      >
        <h5 className="text-xs font-semibold">{title}</h5>
        <p className="font-semibold">{children}</p>
      </Suspense>
    </div>
  );
}

export default QuickViewCard;
