import { t, Trans } from '@lingui/macro';
import { Button, Loader } from '@rotational/beacon-core';
import { Suspense } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { CardListItem } from '@/components/common/CardListItem';
import { useFetchTenantProjects } from '@/features/projects/hooks/useFetchTenantProjects';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useOrgStore } from '@/store';

import { getRecentProject } from '../util';
// interface ProjectDetailsStepProps {
//   tenantID: string;
// }
function ProjectDetailsStep() {
  const navigate = useNavigate();
  const orgDataState = useOrgStore.getState() as any;
  const { tenants } = useFetchTenants();

  const tenantID = tenants?.tenants[0]?.id;

  const { projects, wasProjectsFetched, isFetchingProjects } = useFetchTenantProjects(tenantID);

  if (wasProjectsFetched) {
    // set the projectID in the store
    orgDataState.setProjectID(projects?.tenant_projects[0]?.id);
  }

  const projectDetail = getRecentProject(projects);

  const isDataAvailable = Object.keys(projectDetail || {}).length > 0;

  const redirectToProject = () => {
    navigate(`${PATH_DASHBOARD.PROJECTS}/${projects?.tenant_projects[0]?.id}`);
  };

  return (
    <>
      <Suspense fallback={<Loader size="sm" />}>
        {isFetchingProjects && (
          <div className="flex justify-center">
            <Loader />
          </div>
        )}
        {wasProjectsFetched && projects && (
          <CardListItem
            title={t`Step 1: View Project Details`}
            data={projectDetail}
            itemKey="projectdetail"
          >
            <div className="space-y-3">
              <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
                <p className="w-full text-sm sm:w-4/5">
                  <Trans>
                    View project details below. Generate your API key next to connect producers and
                    consumers to Ensign and start managing your project.
                  </Trans>
                </p>
                <div className="sm:w-1/5 ">
                  <Button
                    className="h-[44px] w-[165px] grow text-sm"
                    disabled={!isDataAvailable}
                    onClick={redirectToProject}
                    data-testid="manage"
                    variant="primary"
                  >
                    <Trans>Manage Project</Trans>
                  </Button>
                </div>
              </div>
            </div>
          </CardListItem>
        )}
      </Suspense>
    </>
  );
}

export default ProjectDetailsStep;
