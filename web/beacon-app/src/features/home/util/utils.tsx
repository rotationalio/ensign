// import { Project } from '../types';
import { useOrgStore } from '@/store';

export const getRecentProject = (projects: any) => {
  const org = useOrgStore.getState() as any;
  let p = [] as any;
  if (projects && projects?.tenant_projects?.length) {
    const recent = projects?.tenant_projects[0]; // TODO: get most recent project instead of first

    const { name, id } = recent; // The project object response from the API
    org.setProjectID(id);
    p = [
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
  return p;
};
