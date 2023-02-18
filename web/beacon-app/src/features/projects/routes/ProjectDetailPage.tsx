import { Heading, Loader } from '@rotational/beacon-core';
import { lazy, Suspense, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import AppLayout from '@/components/layout/AppLayout';
const ProjectDetail = lazy(() => import('../components/ProjectDetail'));
const TopicTable = lazy(() => import('../components/TopicTable'));
const APIKeysTable = lazy(() => import('../components/APIKeysTable'));

const ProjectDetailPage = () => {
  const navigate = useNavigate();
  const param = useParams<{ id: string }>() as any;

  useEffect(() => {
    if (!param || param.id === 'undefined' || param.id === 'null') {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [param, navigate]);

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
        <ProjectDetail projectID={param?.id} />
      </Suspense>

      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <TopicTable />
      </Suspense>

      <Suspense
        fallback={
          <div className="flex justify-center">
            <Loader />
          </div>
        }
      >
        <APIKeysTable />
      </Suspense>
    </AppLayout>
  );
};

export default ProjectDetailPage;
