import { UseMutateFunction } from '@tanstack/react-query';

import { ProjectResponse } from './projectService';

export interface ProjectMutation {
  createNewProject: UseMutateFunction<ProjectResponse, unknown, NewProjectDTO, unknown>;
  reset(): void;
  project: ProjectResponse;
  hasProjectFailed: boolean;
  wasProjectCreated: boolean;
  isCreatingProject: boolean;
  error: any;
}

export type NewProjectDTO = Pick<ProjectResponse, 'name' | 'description'> & {
  tenantID: string;
};

export const isProjectCreated = (
  mutation: ProjectMutation
): mutation is Required<ProjectMutation> =>
  mutation.wasProjectCreated && mutation.project != undefined;
