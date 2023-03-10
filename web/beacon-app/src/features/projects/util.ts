import { formatDate } from '@/utils/formatDate';
export const formatProjectData = (data: any) => {
  if (!data) return [];
  return [
    {
      label: 'Project Name',
      value: data.name,
    },
    {
      label: 'Project ID',
      value: data.id,
    },
    {
      label: 'Date Created',
      value: formatDate(new Date(data.created)),
    },
  ];
};
