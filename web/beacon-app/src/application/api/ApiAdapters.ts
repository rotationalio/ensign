import { APIKey, NewAPIKey } from '@/features/apiKeys/types/ApiKeyServices';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type { NewUserAccount, User } from '@/features/auth/types/RegisterService';
import type { UserTenantResponse } from '@/features/tenants/types/tenantServices';

export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<User>;
  authenticateUser(user: Pick<User, 'username' | 'password'>): Promise<UserAuthResponse>;
  getTenantList(): Promise<UserTenantResponse>;
  createAPIKey(key: NewAPIKey): Promise<APIKey>;
}
