import { t } from '@lingui/macro';

import { formatDate } from '@/utils/formatDate';

import { APIKey } from '../apiKeys/types/apiKeyService';

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
