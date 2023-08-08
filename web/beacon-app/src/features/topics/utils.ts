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

export const getProjectQueryMetaData = (results: any) => {
  console.log('results', results);
  if (!results) {
    return [];
  }
  return Object.keys(results).map((key) => {
    const { metadata } = results[key];
    // Return the metadata as a string
    const metadataString = JSON.stringify(metadata);
    return { metadataString };
  });
};

export const getQueryPaginationCounter = (count: number, total: any) => {
  if (total > 0) {
    return ` ${count} results of ${total} total`;
  }
  return '0 results of 0 total';
};

// TODO:  implement event pagination later to have the right count of events
export const getEventsPaginationCounter = (count: number, total: number) => {
  if (total > 0) {
    return ` ${count} of ${total} `;
  }
  return '0 of 0';
};
