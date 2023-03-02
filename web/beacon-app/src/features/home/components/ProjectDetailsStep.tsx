import { Button } from '@rotational/beacon-core';
import { useNavigate } from 'react-router-dom';

import { CardListItem } from '@/components/common/CardListItem';
import { useFetchProjects } from '@/features/projects/hooks/useFetchProjects';

import { getRecentProject } from '../util';
function ProjectDetailsStep() {
  const navigate = useNavigate();
  const { projects, wasProjectsFetched, getProjects } = useFetchProjects();

  if (!projects) {
    getProjects();
  }

  if (wasProjectsFetched || projects) {
    console.log('projects', projects);
  }

  const isDataAvailable = projects?.projects.length > 0;
  console.log('isDataAvailable', isDataAvailable);

  const redirectToProject = () => {
    navigate(`/projects/${projects?.projects[0].id}`);
  };

  return (
    <>
      <CardListItem
        title="Step 1: View Project Details"
        data={getRecentProject(projects?.projects[0])}
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
