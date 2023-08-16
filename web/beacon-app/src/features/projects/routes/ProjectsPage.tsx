import { useEffect, useState } from 'react';

import { QuickView } from '@/components/common/QuickView';
import AppLayout from '@/components/layout/AppLayout';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import ProjectList from '../components/ProjectList';
import useFetchProjectStats from '../hooks/useFetchProjectStats';
import { getDefaultProjectStats, getProjectStatsHeaders } from '../util';

function ProjectsPage() {
  const [projectStats, setProjectStats] = useState<any>(getDefaultProjectStats());
  const { tenants } = useFetchTenants();

  // display user locale language
  console.log('user locale language', navigator.language);

  const { projectQuickView, error } = useFetchProjectStats(tenants?.tenants[0]?.id);

  useEffect(() => {
    if (projectQuickView) {
      setProjectStats(projectQuickView);
    }
  }, [projectQuickView]);

  useEffect(() => {
    if (error) {
      setProjectStats(getDefaultProjectStats());
    }
  }, [error]);

  return (
    <AppLayout>
      <QuickView data={projectStats} headers={getProjectStatsHeaders()} />
      <ProjectList />
    </AppLayout>
  );
}

export default ProjectsPage;
