import { useLoaderData } from 'react-router-dom';

import { QuickView } from '@/components/common/QuickView';
import AppLayout from '@/components/layout/AppLayout';

import ProjectList from '../components/ProjectList';

function ProjectsPage() {
  const loaderData = useLoaderData() as any;
  return (
    <AppLayout>
      <QuickView data={loaderData?.projectStats} />
      <ProjectList />
    </AppLayout>
  );
}

export default ProjectsPage;
