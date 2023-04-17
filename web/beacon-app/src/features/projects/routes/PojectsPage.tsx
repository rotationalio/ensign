import { Heading } from '@rotational/beacon-core';

import { QuickView } from '@/components/common/QuickView';
import AppLayout from '@/components/layout/AppLayout';

import ProjectList from '../components/ProjectList';

function PojectsPage() {
  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 mt-6 text-lg font-semibold">
        Quick View
      </Heading>
      <QuickView />
      <ProjectList />
    </AppLayout>
  );
}

export default PojectsPage;
