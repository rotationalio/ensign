import { CardListItem } from '@/components/common/CardListItem';

interface ProjectDetailsProps {
  project: { label: string; value: any }[];
}
export const ProjectDetail = ({ project }: ProjectDetailsProps) => {
  return <CardListItem data={project} className="my-5" />;
};

export default ProjectDetail;
