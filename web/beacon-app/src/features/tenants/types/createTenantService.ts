export interface TenantMutation {
  createTenant(): void;
  tenant: any;
  hasTenantFailed: boolean;
  wasTenantFetched: boolean;
  isFetchingTenant: boolean;
  error: any;
}
