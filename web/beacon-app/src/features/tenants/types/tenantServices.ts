export interface UserTenantResponse {
    tenant: ITenant[];
    prev_page_token: string;
    next_page_token: string;
}

export interface ITenant {
    id: string;
    name: string;
    environment_type: string;
}

export interface TenantQuery {
    getTenant(): void;
    tenants: any;
    hasTenantFailed: boolean;
    wasTenantFetched: boolean;
    isFetchingTenant: boolean;
    error: any;
}
