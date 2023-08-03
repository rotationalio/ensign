export interface TopicEvents {
  type: string;
  version: string;
  mimetype: string;
  events: {
    value: number;
    percent: number;
  };
  storage: {
    value: number;
    units: string;
    percent: number;
  };
}

export interface TopicEventsQuery {
  getTopicEvents: () => void;
  hasTopicEventsFailed: boolean;
  isFetchingTopicEvents: boolean;
  topicEvents: TopicEvents[];
  wasTopicEventsFetched: boolean;
  error: any;
}
