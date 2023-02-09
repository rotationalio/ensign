import { useQuery } from "@tanstack/react-query";
import { ProjectDetailQuery } from "../types/projectService";
import { RQK } from "@/constants";
import axiosInstance from "@/application/api/ApiService";
import { projectRequest } from "../api/projectListAPI";

export function UseFetchProjects(): ProjectDetailQuery {
    const query = useQuery([RQK.PROJECT_LIST], projectRequest(axiosInstance), {
        refetchOnWindowFocus: false,
        refetchOnMount: true,
        // set stale time to 15 minutes
        // TODO: Change stale time
        staleTime: 1000 * 60 * 15,
  });

  return {
    getProject: query.refetch,
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data,
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}