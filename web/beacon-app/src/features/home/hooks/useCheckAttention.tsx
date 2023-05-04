import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

export function useCheckAttention() {
  const { tenants } = useFetchTenants();
  const { projects, wasProjectsFetched } = useFetchTenantProjects(tenants?.tenants[0]?.id);

  const hasProject = projects?.tenant_projects?.length > 0;

  return {
    hasProject,
    wasProjectsFetched,
  };
}
