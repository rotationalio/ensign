import { Trans } from '@lingui/macro';
import { Button, Heading } from '@rotational/beacon-core';
import { useState } from 'react';

import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
import RefreshIcon from '@/components/icons/refresh';
import Union from '@/components/icons/union';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import { useFetchTenantProjects } from '../hooks/useFetchTenantProjects';
import NewProjectModal from './NewProject/NewProjectModal';
import ProjectsTable from './ProjectsTable';

function ProjectList() {
  const { tenants } = useFetchTenants();

  const tenantID = tenants?.tenants[0]?.id;

  const { getProjects, isFetchingProjects, projects } = useFetchTenantProjects(tenantID);

  const [isOpenNewProjectModal, setIsOpenNewProjectModal] = useState<boolean>(false);

  const onOpenNewProjectModal = () => {
    setIsOpenNewProjectModal(true);
  };

  const onCloseNewProjectModal = () => {
    setIsOpenNewProjectModal(false);
  };

  const refreshHandler = () => {
    getProjects();
  };

  return (
    <>
      <div className="flex space-x-2 space-y-2">
        <Heading as="h1" className="mb-4 mt-6 text-lg font-semibold">
          Projects
        </Heading>
        <HelpTooltip>
          <Trans>
            A project is a collection of topics. Topics are event streams that your services,
            applications, or models can publish or subscribe to for real-time data flows. Each
            project requires an API key and secret. Secrets are known only to you. If you lose a
            secret, you lose the topic and data generated by the event stream permanently.
          </Trans>
        </HelpTooltip>
      </div>
      <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
        <div className="mt-2">
          <button disabled={isFetchingProjects} onClick={refreshHandler}>
            <RefreshIcon />
          </button>
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
            <Trans>Create Project</Trans>
          </Button>
        </div>
      </div>
      <ProjectsTable projects={projects?.tenant_projects} isLoading={isFetchingProjects} />
      <NewProjectModal
        isOpened={isOpenNewProjectModal}
        onClose={onCloseNewProjectModal}
        data-testid="newProjectModal"
      />
    </>
  );
}

export default ProjectList;
