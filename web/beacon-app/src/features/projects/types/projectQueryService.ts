import { UseMutateFunction } from '@tanstack/react-query';

export interface ProjectQueryMutation {
  getProjectQuery: UseMutateFunction<ProjectQueryDTO, unknown, ProjectQueryDTO, unknown>;
  reset(): void;
  projectQuery: any;
  hasProjectQueryFailed: boolean;
  wasProjectQueryCreated: boolean;
  isCreatingProjectQuery: boolean;
  error: any;
}

export interface ProjectQueryDTO {
  projectID: string;
  query: string;
}

export interface ProjectQuery {
  results: ProjectQueryResult[];
  total_events: number;
  error: string;
}

export interface ProjectQueryResult {
  metadata: object;
  mimetype: string;
  version: string;
  is_base_64_encoded: boolean;
  data: string;
  created: string;
}
