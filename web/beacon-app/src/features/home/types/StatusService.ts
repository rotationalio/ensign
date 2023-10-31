export interface StatusResponse {
  status: string;
  uptime: string;
  version: string;
}

export interface StatusQuery {
  getStatus(): void;
  status: any;
  hasStatusFailed: boolean;
  wasStatusFetched: boolean;
  isFetchingStatus: boolean;
  error: any;
}
