export interface UserTenantResponse {
  tenants: ITenant[];
  prev_page_token: string;
  next_page_token: string;
}

export interface ITenant {
  id: string;
  name: string;
  environment_type: string;
  created: string;
  modified: string;
}

export interface TenantsQuery {
  getTenants(): void;
  tenants: any;
  hasTenantsFailed: boolean;
  wasTenantsFetched: boolean;
  isFetchingTenants: boolean;
  error: any;
}
