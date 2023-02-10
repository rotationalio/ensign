export interface Topic {
  id: string;
  name: string;
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
