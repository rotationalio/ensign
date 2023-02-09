import { ApiAdapters } from "@/application/api/ApiAdapters";
import { getValidApiResponse, Request } from "@/application/api/ApiService";
import { APP_ROUTE } from "@/constants";
import { UserProjectResponse } from "../types/projectService";

export function projectRequest(request: Request): ApiAdapters['getProjectList'] {
    return async () => {
        const response = (await request(`${APP_ROUTE.PROJECTS}`, {
            method: 'GET',
        })) as any;
        return getValidApiResponse<UserProjectResponse>(response);
    }
}