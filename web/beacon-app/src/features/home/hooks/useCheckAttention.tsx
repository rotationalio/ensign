import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';
import { ProjectStatus } from '@/features/projects/types/Project';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
export function useCheckAttention() {
  const { tenants } = useFetchTenants();
  const { projects, wasProjectsFetched } = useFetchTenantProjects(tenants?.tenants[0]?.id);

  const hasProject = projects?.tenant_projects?.length > 0;
  const hasOneProjectAndIsIncomplete =
    projects?.tenant_projects?.length === 1 &&
    projects?.tenant_projects[0]?.status === ProjectStatus.INCOMPLETE;

  return {
    hasProject,
    wasProjectsFetched,
    hasOneProjectAndIsIncomplete,
  };
}
