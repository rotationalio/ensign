import { Heading, Loader } from '@rotational/beacon-core';
import { Suspense } from 'react';

import AppLayout from '@/components/layout/AppLayout';

const TopicDetailPage = () => {
  return (
    <AppLayout>
      <Heading as="h1">Topic Name</Heading>
      <Suspense
        fallback={
          <div>
            <Loader />
          </div>
        }
      ></Suspense>
    </AppLayout>
  );
};

export default TopicDetailPage;
