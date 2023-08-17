import { Trans } from '@lingui/macro';
import { memo } from 'react';
const LeftPanel = () => {
  return (
    <>
      <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-2/6">
        <h1 className="text-center font-bold">
          <Trans>Building event-driven applications can be fast, convenient and even fun! 🎉</Trans>
        </h1>
        <p className="text-center font-bold">
          <Trans>Start today on our no-cost Starter Plan.</Trans>
        </p>
        <p>
          <Trans>
            If you have always wanted to try out eventing, but couldn&apos;t justify the high cost
            of entry or the expertise required, Ensign is for you!
          </Trans>
        </p>
        <p>
          <Trans>Want to build...</Trans>
        </p>
        <ul className="ml-5 list-disc">
          <li>
            <Trans>new prototypes without refactoring legacy database schemas?</Trans>
          </li>
          <li>
            <Trans>real-time dashboards and analytics in days rather than months?</Trans>
          </li>
          <li>
            <Trans>rich, tailored experiences so your users know how much they mean to you?</Trans>
          </li>
          <li>
            <Trans>
              MLOps pipelines that bridge the gap between the training and deployment phases?
            </Trans>
          </li>
        </ul>
        <p>
          <Trans>Let&apos;s do it hero. 💪</Trans>
        </p>
      </div>
    </>
  );
};

export default memo(LeftPanel);
