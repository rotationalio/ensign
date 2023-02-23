import { Heading } from '@rotational/beacon-core';

import AppLayout from '@/components/layout/AppLayout';
import { useFetchProjects } from '@/features/projects/hooks/useFetchProjects';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { getRecentTenant } from '@/utils/formatData';

import QuickViewSummary from '../components/QuickViewSummary';
import Steps from '../components/Steps';

export default function Home() {
  const { tenants, getTenants } = useFetchTenants();
  const recentTenant = getRecentTenant(tenants);
  const { projects, getProjects } = useFetchProjects(recentTenant?.id);

  // fetch data and catching them in the state
  if (!projects) {
    getProjects();
  }
  if (!tenants) {
    getTenants();
  }

  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        Quick view
      </Heading>
      <QuickViewSummary />
      <Heading as="h1" className="mb-4 pt-10 text-lg font-semibold">
        Follow 3 simple steps to set up your event stream and set your data in motion.
      </Heading>
      <Steps />
    </AppLayout>
  );
}
