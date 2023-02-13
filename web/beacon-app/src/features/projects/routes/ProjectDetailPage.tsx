import { Container, Heading, Loader } from '@rotational/beacon-core';
import { lazy, Suspense, useEffect, useState } from 'react';

import { queryCache } from '@/application/config/react-query';
import { SentryErrorBoundary } from '@/components/Error';
import { RQK } from '@/constants';
const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

export const ProjectDetailPage = () => {
  const [projectID, setProjectID] = useState<any>();

  const projects = queryCache.find(RQK.PROJECTS) as any;
  // This should get project from react-query cache
  useEffect(() => {
    if (projects) {
      setProjectID(projects[0].id as string);
    }
  }, [projects]);
  // get the first project in the list;

  return (
    <Container max={696} centered>
      <Heading as="h1" className="flex ">
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
        <TopicTable projectID={projectID} />
      </Suspense>

      <Suspense fallback={<div>Loading...</div>}>
        <APIKeysTable />
      </Suspense>
    </Container>
  );
};
