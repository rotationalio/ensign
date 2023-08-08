import { UseMutateFunction } from '@tanstack/react-query';

export interface ProjectQueryMutation {
  getProjectQuery: UseMutateFunction<ProjectQueryResponse, unknown, ProjectQueryDTO, unknown>;
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

export interface MetaDataResponse {
  key: string;
  value: string;
}
export interface ProjectQueryResult {
  metadata: MetaDataResponse[];
  mimetype: string;
  created: string;
  is_base64_encoded: boolean;
  data: any;
  version: string;
}

export interface ProjectQueryResponse {
  results: ProjectQueryResult[];
  total_events: number;
}
