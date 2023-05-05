import { UseMutateFunction } from '@tanstack/react-query';

import { Topic } from '@/features/topics/types/topicService';

export interface TopicMutation {
  createTopic: UseMutateFunction<Topic, unknown, NewTopicDTO, unknown>;
  reset(): void;
  topic: any;
  hasTopicFailed: boolean;
  wasTopicCreated: boolean;
  isCreatingTopic: boolean;
  error: any;
}

export interface NewTopic {
  topic_name: string;
}

export type NewTopicDTO = {
  projectID: string;
} & NewTopic;

export const isTopicCreated = (mutation: TopicMutation): mutation is Required<TopicMutation> =>
  mutation.wasTopicCreated && mutation.topic != undefined;
