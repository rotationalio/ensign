import { QuickView } from '@/components/common/QuickView';
import AppLayout from '@/components/layout/AppLayout';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import ProjectList from '../components/ProjectList';
import useFetchProjectStats from '../hooks/useFetchProjectStats';

function ProjectsPage() {
  const { tenants } = useFetchTenants();

  const { projectQuickView } = useFetchProjectStats(tenants?.tenants[0]?.id);

  return (
    <AppLayout>
      <QuickView data={projectQuickView} />
      <ProjectList />
    </AppLayout>
  );
}

export default ProjectsPage;
