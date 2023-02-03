import { APIKey, NewAPIKey } from '@/features/apiKeys/types/ApiKeyServices';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type { NewUserAccount, User, NewUserResponseData } from '@/features/auth/types/RegisterService';
import type { UserTenantResponse } from '@/features/tenants/types/tenantServices';

export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<NewUserResponseData>;
  authenticateUser(user: Pick<User, 'email' | 'password'>): Promise<UserAuthResponse>;
  getTenantList(): Promise<UserTenantResponse>;
  createAPIKey(key: NewAPIKey): Promise<APIKey>;
}
