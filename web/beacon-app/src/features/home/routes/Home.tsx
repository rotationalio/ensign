import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';

import QuickStart from '../components/QuickStart';
import QuickViewSummary from '../components/QuickViewSummary';
export default function Home() {
  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        Quick View
      </Heading>
      <QuickViewSummary />
      <Heading as="h1" className="mb-4 pt-10 text-lg font-semibold">
        Follow 3 simple steps to set up your event stream and set your data in motion.
      </Heading>
      <QuickStart />
    </AppLayout>
  );
}
