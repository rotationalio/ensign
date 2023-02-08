import Page from '@/components/common/Page';

import QuickViewSummary from '../components/QuickViewSummary';
import Steps from '../components/Steps';

export default function Home() {
  return (
    <Page>
      <h1 className="mb-4 text-sm font-semibold">Quick view</h1>
      <QuickViewSummary />
      <p className="mb-6 mt-8 text-sm font-semibold">
        Follow 3 simple steps to set up your event stream and set your data in motion.
      </p>
      <Steps />
    </Page>
  );
}
