import { UseMutateFunction } from "@tanstack/react-query";

export interface Topic {
  id: string;
  name: string;
  state: string;
  created?: string;
  modified?: string;
}

export interface Topics {
  project_id: string;
  topics: Topic[];
  prev_page_token: string;
  next_page_token: string;
}

export interface TopicsQuery {
  getTopics: () => Promise<Topics | undefined | unknown>;
  topics: any;
  hasTopicsFailed: boolean;
  wasTopicsFetched: boolean;
  isFetchingTopics: boolean;
  error: any;
}

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
  name: string;
}

export type NewTopicDTO = {
  projectID: string;
} & NewTopic;

export const isTopicCreated = (mutation: TopicMutation): mutation is Required<TopicMutation> =>
  mutation.wasTopicCreated && mutation.topic != undefined;