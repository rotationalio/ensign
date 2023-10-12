import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';

import QuickStart from '../components/QuickStart';
import QuickViewSummary from '../components/QuickViewSummary';
import WelcomeAttention from '../components/WelcomeAttention';

export default function Home() {
  return (
    <AppLayout>
      <WelcomeAttention />
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        <Trans>Quick View</Trans>
      </Heading>
      <QuickViewSummary />
      <QuickStart />
    </AppLayout>
  );
}
