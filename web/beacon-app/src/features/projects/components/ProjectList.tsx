import { Button, Heading } from '@rotational/beacon-core';
import { useState } from 'react';

import RefreshIcon from '@/components/icons/refresh';
import Union from '@/components/icons/union';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import { useFetchTenantProjects } from '../hooks/useFetchTenantProjects';
import NewProjectModal from './NewProject/NewProjectModal';
import ProjectsTable from './ProjectsTable';
import ProjectTooltip from './ProjectTooltip';

function ProjectList() {
  const { tenants } = useFetchTenants();

  const tenantID = tenants?.tenants[0]?.id;

  const { projects } = useFetchTenantProjects(tenantID);

  const [isOpenNewProjectModal, setIsOpenNewProjectModal] = useState<boolean>(false);

  const onOpenNewProjectModal = () => {
    setIsOpenNewProjectModal(true);
  };

  const onCloseNewProjectModal = () => {
    setIsOpenNewProjectModal(false);
  };

  return (
    <>
      <div className="flex space-x-2 space-y-2">
        <Heading as="h1" className="mb-4 mt-6 text-lg font-semibold">
          Projects
        </Heading>
        <ProjectTooltip />
      </div>
      <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
        <div className="mt-2">
          <RefreshIcon />
        </div>
        <div className="flex items-center gap-3"></div>
        <div>
          <Button
            className="flex items-center gap-1"
            size="small"
            data-testid="create-project-btn"
            onClick={onOpenNewProjectModal}
          >
            <Union className="fill-white" />
            Create Project
          </Button>
        </div>
      </div>
      <ProjectsTable projects={projects?.tenant_projects} />
      <NewProjectModal isOpened={isOpenNewProjectModal} onClose={onCloseNewProjectModal} />
    </>
  );
}

export default ProjectList;
