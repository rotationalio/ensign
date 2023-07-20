import { QuickView } from '@/components/common/QuickView';
import AppLayout from '@/components/layout/AppLayout';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import ProjectList from '../components/ProjectList';
import useFetchProjectStats from '../hooks/useFetchProjectStats';
import { getDefaultProjectStats, getProjectStatsHeaders } from '../util';

function ProjectsPage() {
  const { tenants } = useFetchTenants();

  const { projectQuickView } = useFetchProjectStats(tenants?.tenants[0]?.id);

  const getProjectStats = () => {
    if (!projectQuickView) return getDefaultProjectStats();
    return projectQuickView;
  };

  return (
    <AppLayout>
      <QuickView data={getProjectStats()} headers={getProjectStatsHeaders()} />
      <ProjectList />
    </AppLayout>
  );
}

export default ProjectsPage;
