import { Container, Heading, Loader } from '@rotational/beacon-core';
import { lazy, Suspense } from 'react';
import { useParams } from 'react-router-dom';

// import { SentryErrorBoundary } from '@/components/Error';

const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const projectID = useParams<{ id: string }>() as string;
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
        {/* <SentryErrorBoundary fallback={<div>Something went wrong</div>}> */}
        <ProjectDetail projectID={projectID} />
        {/* </SentryErrorBoundary> */}
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

export default ProjectDetailPage;
