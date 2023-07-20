import { QuickViewData } from '@/hooks/useFetchQuickView/quickViewService';
export interface Topic {
  id: string;
  topic_name: string;
  status: string;
  created?: string;
  modified?: string;
  publishers?: number;
  subscribers?: number;
  data_storage?: QuickViewData;
}

export interface TopicsResponse {
  project_id: string;
  topics: Topic[];
  prev_page_token: string;
  next_page_token: string;
}

export interface TopicsQuery {
  getTopics: () => void;
  topics: TopicsResponse;
  hasTopicsFailed: boolean;
  wasTopicsFetched: boolean;
  isFetchingTopics: boolean;
  error: any;
}

export interface TopicQuery {
  getTopic: () => void;
  topic: Topic;
  hasTopicFailed: boolean;
  wasTopicFetched: boolean;
  isFetchingTopic: boolean;
  error: any;
}

export interface TopicQuickViewResponse {
  data: IStats[];
}

export interface TopicStatsQuery {
  getTopicStats: () => void;
  hasTopicStatsFailed: boolean;
  isFetchingTopicStats: boolean;
  topicStats: IStats[];
  wasTopicStatsFetched: boolean;
  error: any;
}
