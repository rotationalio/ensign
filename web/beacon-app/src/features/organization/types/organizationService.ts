export interface OrgResponse {
  id: string;
  name: string;
  domain: string;
  created: string;
  modified: string;
}

export interface OrgDetailQuery {
  getOrgDetail(): void;
  org: any;
  hasOrgFailed: boolean;
  wasOrgFetched: boolean;
  isFetchingOrg: boolean;
  error: any;
}

export interface OrgListQuery {
  getOrgList(): void;
  organizations: any;
  hasOrgListFailed: boolean;
  wasOrgListFetched: boolean;
  isFetchingOrgList: boolean;
  error: any;
}

export interface OrgListResponse {
  organizations: OrgResponse[];
}
