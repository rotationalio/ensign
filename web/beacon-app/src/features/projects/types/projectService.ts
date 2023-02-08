export type ProjectResponse = {
  id: string;
  name: string;
};

export type ProjectDetailQuery = {
  // getProject: (id: ProjectDetailDTO) => Promise<ProjectResponse | undefined>;
  project: any;
  hasProjectFailed: boolean;
  wasProjectFetched: boolean;
  isFetchingProject: boolean;
  error: any;
};

export type ProjectDetailDTO = Pick<ProjectResponse, 'id'>;
