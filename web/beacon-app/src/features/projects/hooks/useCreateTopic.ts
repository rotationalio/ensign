import { TopicMutation } from "@/features/topics/types/topicService";
import { useMutation } from "@tanstack/react-query";
import { createProjectTopic } from "../api/createTopicApiService";
import axiosInstance from "@/application/api/ApiService";

export function useCreateTopic(): TopicMutation {
    const mutation = useMutation(createProjectTopic(axiosInstance), {
        retry: 0,
       /*  TODO: Add on success */
    })
    return {
        createTopic: mutation.mutate,
        reset: mutation.reset,
        topic: mutation.data as TopicMutation['topic'],
        hasTopicFailed: mutation.isError,
        wasTopicCreated: mutation.isSuccess,
        isCreatingTopic: mutation.isLoading,
        error: mutation.error,
    }
}