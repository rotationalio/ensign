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
