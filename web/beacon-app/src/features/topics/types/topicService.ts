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

export interface Topics {
  project_id: string;
  topics: Topic[];
  prev_page_token: string;
  next_page_token: string;
}

export interface TopicsQuery {
  getTopics: () => void;
  topics: any;
  hasTopicsFailed: boolean;
  wasTopicsFetched: boolean;
  isFetchingTopics: boolean;
  error: any;
}
