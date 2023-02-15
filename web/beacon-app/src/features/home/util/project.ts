// import { Project } from '../types';

export const getRecentProject = (projects: any | undefined) => {
  if (projects && projects.length > 0) {
    const recent = projects[0]; // TODO: get most recent project instead of first
    const { name, id } = recent; // The project object response from the API
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
