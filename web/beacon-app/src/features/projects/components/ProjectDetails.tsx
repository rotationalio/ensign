import { useCallback, useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

import { ProjectDetailDTO } from '../types/projectService';
interface ProjectDetailsProps {
  id: ProjectDetailDTO;
}
export const ProjectDetails = ({ id }: ProjectDetailsProps) => {
  const [items, setItems] = useState<any>([]);
  const { project, isFetchingProject, wasProjectFetched, error } = useFetchProject(id);
  const projectDetails = useCallback(
    () =>
      wasProjectFetched
        ? project.map((project: any) => {
            return {
              label: project.id,
              value: project.name,
            };
          })
        : [],
    [project, wasProjectFetched]
  );

  if (wasProjectFetched && !isFetchingProject && !error) {
    setItems(projectDetails);
  }

  return <CardListItem data={items} />;
};

export default ProjectDetails;
