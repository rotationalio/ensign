// import { Project } from '../types';
import { t } from '@lingui/macro';

import { useOrgStore } from '@/store';
export const getRecentProject = (projects: any) => {
  const org = useOrgStore.getState() as any;
  let p = [] as any;
  if (projects && projects?.tenant_projects?.length) {
    const recent = projects?.tenant_projects[0]; // TODO: get most recent project instead of first

    const { name, id } = recent; // The project object response from the API
    org.setProjectID(id);
    org.setProject({
      name,
      id,
    });
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

export const getDefaultHomeStats = () => {
  return [
    {
      name: t`Active Projects`,
      value: 0,
    },
    {
      name: t`Topics`,
      value: 0,
    },
    {
      name: t`API Keys`,
      value: 0,
    },
    {
      name: t`Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};
export const getHomeStatsHeaders = () => {
  return [t`Active Projects`, t`Topics`, t`API Keys`, t`Storage`];
};
