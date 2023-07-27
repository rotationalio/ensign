import { t } from '@lingui/macro';

import { Topic } from '@/features/topics/types/topicService';
import { formatDate } from '@/utils/formatDate';

import { APIKey } from '../apiKeys/types/apiKeyService';
import type { Project } from './types/Project';

export const formatProjectData = (data: any) => {
  if (!data) return [];
  return [
    {
      label: 'Project Name',
      value: data?.name,
    },
    {
      label: 'Project ID',
      value: data?.id,
    },
    {
      label: 'Date Created',
      value: formatDate(new Date(data?.created)),
    },
  ];
};

type APIKeyActions = {
  handleOpenRevokeAPIKeyModal: (key: APIKey) => void;
};

export const getApiKeys = (apiKeys: any, actions?: APIKeyActions) => {
  if (!apiKeys?.api_keys || apiKeys?.api_keys.length === 0) return [];
  return Object.keys(apiKeys?.api_keys).map((key) => {
    const { id, name, client_id, permissions, status, last_used, created } = apiKeys.api_keys[key];
    return {
      id,
      name,
      client_id,
      permissions,
      status,
      last_used,
      created,
      actions: [
        {
          label: t`Revoke API Key`,
          onClick: () => actions?.handleOpenRevokeAPIKeyModal(apiKeys.api_keys[key]),
        },
      ],
    };
  }) as any;
};

export const getTopics = (topics: any) => {
  if (!topics?.topics || topics?.topics.length === 0) return [];
  return Object.keys(topics?.topics).map((topic) => {
    const { id, topic_name, status, created, modified } = topics.topics[topic];
    return { id, topic_name, status, created, modified };
  }) as any;
};

export const getNormalizedDataStorage = (value?: number, units?: string) => {
  if (!value) {
    return '0GB';
  }
  return String(value) + units;
};

export const getDefaultProjectStats = () => {
  return [
    {
      name: t`Active Projects`,
      value: 0,
    },
    {
      name: t`Topics`,
      value: 0,
    },
    {
      name: t`API Keys`,
      value: 0,
    },
    {
      name: t`Data Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};

export const getProjectStatsHeaders = () => {
  return [t`Active Projects`, t`Topics`, t`API Keys`, t`Data Storage`];
};

export const getInitialColumns = () => {
  const initialColumns = [
    { Header: t`Topic Name`, accessor: 'topic_name' },
    { Header: t`Status`, accessor: 'status' },
    {
      Header: t`Publishers`,
      accessor: (t: Topic) => {
        const publishers = t?.publishers;
        return publishers || '---';
      },
    },
    {
      Header: t`Subscribers`,
      accessor: (t: Topic) => {
        const subscribers = t?.subscribers;
        return subscribers || '---';
      },
    },
    {
      Header: t`Data Storage`,
      accessor: (t: Topic) => {
        const value = t?.data_storage?.value;
        const units = t?.data_storage?.units;
        return getNormalizedDataStorage(value, units);
      },
    },
    {
      Header: t`Date Created`,
      accessor: (date: any) => {
        return formatDate(new Date(date?.created));
      },
    },
  ];

  return initialColumns;
};

export const getFormattedProjectData = (project: Project) => {
  const { description, status, created, owner } = project || {};
  return [
    {
      label: t`Project Status`,
      value: status || '---',
    },
    {
      label: t`Description`,
      value: description || '---',
    },
    {
      label: t`Owner`,
      value: owner?.name || '---',
    },
    {
      label: t`Date Created`,
      value: formatDate(new Date(created)),
    },
  ];
};
