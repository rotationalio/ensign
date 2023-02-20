// import { Project } from '../types';
import { useOrgStore } from '@/store';

export const getRecentProject = (projects: any | undefined) => {
  const org = useOrgStore.getState() as any;
  if (projects && projects.length > 0) {
    const recent = projects[0]; // TODO: get most recent project instead of first
    const { name, id } = recent; // The project object response from the API
    org.setProjectID(id);
    return [
      {
        label: 'Project Name',
        value: name,
      },
      {
        label: 'Project ID',
        value: id,
      },
    ];
  }
  return [];
};
