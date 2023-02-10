export interface OrgResponse {
    id: string;
    name: string;
    domain: string;
    created: string;
    modified: string;
}

export interface OrgDetailQuery {
    org: any;
    hasOrgFailed: boolean;
    wasOrgFetched: boolean;
    isFetchingOrg: boolean;
    error: any;
}

export type OrgDetailDTO = Pick<OrgResponse, 'id'>;