import { useCallback, useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { useFetchProject } from '@/features/projects/hooks/useFetchProject';

interface ProjectDetailsProps {
  projectID: string;
}
export const ProjectDetail = ({ projectID }: ProjectDetailsProps) => {
  const [items, setItems] = useState<any>([]);
  console.log('[ProjectDetail] projectID', projectID);
  const { project, isFetchingProject, wasProjectFetched, error } = useFetchProject(projectID);
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

export default ProjectDetail;
