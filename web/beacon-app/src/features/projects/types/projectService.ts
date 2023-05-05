export interface ProjectsResponse {
  project: ProjectResponse[];
  prev_page_token: string;
  next_page_token: string;
}

export type ProjectResponse = {
  id: string;
  name: string;
  created: string;
  modified?: string;
  description?: string;
  status?: string;
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

export type ProjectsQuery = {
  getProjects(): void;
  projects: any;
  hasProjectsFailed: boolean;
  wasProjectsFetched: boolean;
  isFetchingProjects: boolean;
  error: any;
};

export interface ProjectQuickViewData {
  name: string;
  value: number;
  units?: string;
  percent?: number;
}

export interface ProjectQuickViewResponse {
  data: ProjectQuickViewData[];
}

export interface ProjectQuickViewQuery {
  getProjectQuickView: () => void;
  hasProjectQuickViewFailed: boolean;
  isFetchingProjectQuickView: boolean;
  projectQuickView: any;
  wasProjectQuickViewFetched: boolean;
  error: any;
}
