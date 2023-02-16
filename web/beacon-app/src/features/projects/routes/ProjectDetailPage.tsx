import { Heading, Loader } from '@rotational/beacon-core';
import { lazy, Suspense } from 'react';
import { useParams } from 'react-router-dom';

import { SentryErrorBoundary } from '@/components/Error';
import AppLayout from '@/components/layout/AppLayout';
const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const projectID = useParams<{ id: string }>() as string;

  if (!projectID) return <div>Project not found</div>;

  return (
    <AppLayout>
      <Heading as="h1" className="mb-4 text-lg font-semibold">
        Project Detail Page
      </Heading>
      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <SentryErrorBoundary fallback={<div>Something went wrong</div>}>
          <ProjectDetail projectID={projectID} />
        </SentryErrorBoundary>
      </Suspense>

      <Suspense fallback={<div>Loading...</div>}>
        <TopicTable />
      </Suspense>

      <Suspense fallback={<div>Loading...</div>}>
        <APIKeysTable />
      </Suspense>
    </AppLayout>
  );
};

export default ProjectDetailPage;
