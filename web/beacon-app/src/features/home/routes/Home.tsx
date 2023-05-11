import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';

import QuickStart from '../components/QuickStart';
import QuickViewSummary from '../components/QuickViewSummary';
import WelcomeAttention from '../components/WelcomeAttention';
import { useCheckAttention } from '../hooks/useCheckAttention';
export default function Home() {
  const { hasProject, wasProjectsFetched, hasOneProjectAndIsIncomplete } = useCheckAttention();
  return (
    <AppLayout>
      {!hasProject && wasProjectsFetched && <WelcomeAttention />}
      {hasOneProjectAndIsIncomplete && <WelcomeAttention />}
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        <Trans>Quick View</Trans>
      </Heading>
      <QuickViewSummary />
      <Heading as="h1" className="mb-4 pt-10 text-lg font-semibold">
        <Trans>Get Started</Trans>
      </Heading>
      <QuickStart />
    </AppLayout>
  );
}
