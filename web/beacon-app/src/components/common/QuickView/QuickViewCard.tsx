import { Loader } from '@rotational/beacon-core';
import { ReactNode, Suspense } from 'react';
type QuickViewCardProps = {
  title: string;
  color: string;
  children: ReactNode;
};

function QuickViewCard({ title, color, children, ...props }: QuickViewCardProps) {
  return (
    <div
      style={{ backgroundColor: color }}
      className="flex h-[100px] w-full flex-col rounded-2xl px-8 py-4"
      {...props}
    >
      <Suspense
        fallback={
          <div className="justify-center text-center">
            <Loader />
          </div>
        }
      >
        <h5 className="pb-3 pt-1 font-semibold">{title}</h5>
        <p className="text-2xl font-semibold">{children}</p>
      </Suspense>
    </div>
  );
}

export default QuickViewCard;
