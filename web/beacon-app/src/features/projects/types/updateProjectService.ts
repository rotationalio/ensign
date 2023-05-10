import { UseMutateFunction } from '@tanstack/react-query';

import { ProjectResponse } from './projectService';

export interface ProjectUpdateMutation {
  updateProject: UseMutateFunction<ProjectResponse, unknown, UpdateProjectDTO, unknown>;
  reset(): void;
  project: ProjectResponse;
  hasProjectFailed: boolean;
  wasProjectCreated: boolean;
  isCreatingProject: boolean;
  error: any;
}

export type UpdateProjectDTO = {
  projectID: string;
  payload: Partial<Omit<ProjectResponse, 'id' | 'created' | 'modified'>>;
};

export const isProjectUpdated = (
  mutation: ProjectUpdateMutation
): mutation is Required<ProjectUpdateMutation> =>
  mutation.wasProjectCreated && mutation.project != undefined;
