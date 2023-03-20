import { Loader } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

import { formatProjectData } from '../util';
interface ProjectDetailProps {
  projectID: string;
}

export const ProjectDetail = ({ projectID }: ProjectDetailProps) => {
  const { project, isFetchingProject, wasProjectFetched } = useFetchProject(projectID);
  const getFormattedProjectData = formatProjectData(project);
  return (
    <>
      {isFetchingProject && (
        <div className="flex justify-center">
          <Loader />
        </div>
      )}
      {wasProjectFetched && project && (
        <CardListItem data={getFormattedProjectData} className="my-5" />
      )}
    </>
  );
};

export default ProjectDetail;
