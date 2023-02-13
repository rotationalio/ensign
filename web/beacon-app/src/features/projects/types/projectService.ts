export interface ProjectsResponse {
  project: ProjectResponse[];
  prev_page_token: string;
  next_page_token: string;
}

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

export type ProjectQuery = {
  getProjects(): void;
  project: any;
  hasProjectFailed: boolean;
  wasProjectFetched: boolean;
  isFetchingProject: boolean;
  error: any;
};
