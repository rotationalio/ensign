import { t } from '@lingui/macro';

import { formatDate } from '@/utils/formatDate';

import type { TopicEvents } from '../topics/types/topicEventsService';
import type { Topic } from '../topics/types/topicService';

export const getDefaultTopicStats = () => {
  return [
    {
      name: t`Online Publishers`,
      value: 0,
      units: '',
    },
    {
      name: t`Online Subscribers`,
      value: 0,
      units: '',
    },
    {
      name: t`Total Events`,
      value: 0,
    },
    {
      name: t`Data Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};

export const getTopicStatsHeaders = () => {
  return [t`Online Publishers`, t`Online Subscribers`, t`Total Events`, t`Data Storage`];
};

export const getFormattedTopicData = (topic: Topic) => {
  return [
    {
      label: t`Topic ID`,
      value: topic?.id,
    },
    {
      label: t`Status`,
      value: topic?.status,
    },
    {
      label: t`Created`,
      value: formatDate(new Date(topic?.created as string)),
    },
    {
      label: t`Modified`,
      value: formatDate(new Date(topic?.modified as string)),
    },
  ];
};

export const getEventDetailColumns = () => {
  const initialColumns = [
    {
      Header: t`Event Type`,
      accessor: 'type',
    },
    {
      Header: t`Version`,
      accessor: 'version',
    },
    {
      Header: t`MIME Type`,
      accessor: 'mimetype',
    },
    {
      Header: t`# of Events`,
      accessor: (event: TopicEvents) => {
        return event?.events?.value;
      },
    },
    {
      Header: t`% of Events`,
      accessor: (event: TopicEvents) => {
        return `${event?.events?.percent}%`;
      },
    },
    {
      Header: t`Storage Volume`,
      accessor: (event: TopicEvents) => {
        return `${event?.storage?.value} ${event?.storage?.units}`;
      },
    },
    {
      Header: t`% of Volume`,
      accessor: (event: TopicEvents) => {
        return `${event?.storage?.percent}%`;
      },
    },
  ];

  return initialColumns;
};

export const getFormattedEventsDetailData = (events: TopicEvents) => {
  return [
    {
      label: t`Event Type`,
      value: events?.type,
    },
    {
      label: t`Version`,
      value: events?.version,
    },
    {
      label: t`MIME Type`,
      value: events?.mimetype,
    },
    {
      label: t`# of Events`,
      value: events?.events?.value,
    },
    {
      label: t`% of Events`,
      value: `${events?.events?.percent}%`,
    },
    {
      label: t`Storage Volume`,
      value: `${events?.storage?.value} ${events?.storage?.units}`,
    },
    {
      label: t`% of Volume`,
      value: `${events?.storage?.percent}%`,
    },
  ];
};
