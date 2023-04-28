import { ApiAdapters } from "@/application/api/ApiAdapters";
import { Request, getValidApiResponse } from "@/application/api/ApiService";
import { APP_ROUTE } from "@/constants";
import { NewTopicDTO, Topic } from "@/features/topics/types/topicService";
import { AxiosResponse } from "axios";

export function createProjectTopic(request: Request): ApiAdapters['createProjectTopic'] {
    return async ({ projectID, name }: NewTopicDTO) => {
      const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/topics`, {
        method: 'POST',
        data: JSON.stringify({
          name,
        }),
      })) as unknown as AxiosResponse;
  
      return getValidApiResponse<Topic>(response);
    };
  }
  