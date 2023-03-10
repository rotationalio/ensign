import { CardListItem } from '@/components/common/CardListItem';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

import { formatProjectData } from '../util';
interface ProjectDetailsProps {
  projectID: string;
}
export const ProjectDetail = ({ projectID }: ProjectDetailsProps) => {
  const { project } = useFetchProject(projectID);
  const getFormattedProjectData = formatProjectData(project);

  return <CardListItem data={getFormattedProjectData} className="my-5" />;
};

export default ProjectDetail;
