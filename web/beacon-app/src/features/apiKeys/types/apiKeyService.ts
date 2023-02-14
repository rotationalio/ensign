export interface APIKey {
  id: string;
  client_id: string;
  client_secret: string;
  name: string;
  owner: string;
  permissions: string[];
  created?: string;
  modified?: string;
}

export interface APIKeysQuery {
  getApiKeys: () => void;
  apiKeys: any;
  hasApiKeysFailed: boolean;
  wasApiKeysFetched: boolean;
  isFetchingApiKeys: boolean;
  error: any;
}

export type NewAPIKey = Omit<APIKey, 'id'>;
