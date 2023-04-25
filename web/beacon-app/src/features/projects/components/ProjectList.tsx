import { AriaButton as Button, Heading } from '@rotational/beacon-core';

import Union from '@/components/icons/union';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import { useFetchTenantProjects } from '../hooks/useFetchTenantProjects';
import ProjectsTable from './ProjectsTable';

function ProjectList() {
  const { tenants } = useFetchTenants();

  const tenantID = tenants?.tenants[0]?.id;

  const { projects } = useFetchTenantProjects(tenantID);

  return (
    <>
      <Heading as="h1" className="mb-4 mt-6 text-lg font-semibold">
        Projects
      </Heading>
      <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
        <div className="flex items-center gap-3"></div>
        <div>
          <Button
            className="flex items-center gap-1 bg-gray-400 text-xs text-white"
            size="small"
            isDisabled
            data-testid="create__project-btn"
          >
            <Union className="fill-white" />
            Create Project
          </Button>
        </div>
      </div>
      <ProjectsTable projects={projects?.tenant_projects} />
    </>
  );
}

export default ProjectList;
