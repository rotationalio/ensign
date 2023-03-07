import { Button } from '@rotational/beacon-core';
import React from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { CardListItem } from '@/components/common/CardListItem';
import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';

import { getRecentProject } from '../util';
interface ProjectDetailsStepProps {
  tenantID: string;
}
function ProjectDetailsStep({ tenantID }: ProjectDetailsStepProps) {
  const navigate = useNavigate();

  const { projects } = useFetchTenantProjects(tenantID);

  const projectDetail = getRecentProject(projects);

  const isDataAvailable = Object.keys(projectDetail || {}).length > 0;

  const redirectToProject = () => {
    navigate(`${PATH_DASHBOARD.PROJECTS}/${projects?.tenant_projects[0]?.id}`);
  };

  return (
    <>
      <CardListItem
        title="Step 1: View Project Details"
        data={projectDetail || []}
        itemKey="projectdetail"
      >
        <div className="space-y-3">
          <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
            <p className="w-full text-sm sm:w-4/5">
              View project details below. Generate your API key next to connect producers and
              consumers to Ensign and start managing your project.
            </p>
            <div className="sm:w-1/5 ">
              <Button
                className="h-[44px] w-[165px] grow text-sm"
                disabled={!isDataAvailable}
                onClick={redirectToProject}
                data-testid="manage"
              >
                Manage Project
              </Button>
            </div>
          </div>
        </div>
      </CardListItem>
    </>
  );
}

export default ProjectDetailsStep;
