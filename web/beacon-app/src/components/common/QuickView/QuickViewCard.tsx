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
      className="flex h-[100px] w-full flex-col justify-between rounded-xl py-4 px-6"
    >
      <Suspense
        fallback={
          <div className="justify-center text-center">
            <Loader />
          </div>
        }
      >
        <h5 className="text-sm font-semibold" data-cy={title.toLowerCase()}>
          {title}
        </h5>
        <p className="text-md font-semibold" data-cy={title.toLowerCase() + '_value'}>
          {children}
        </p>
      </Suspense>
    </div>
  );
}

export default QuickViewCard;
