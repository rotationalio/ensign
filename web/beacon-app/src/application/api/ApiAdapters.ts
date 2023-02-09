/* eslint-disable no-restricted-imports */
import { APIKey, NewAPIKey } from '@/features/apiKeys/types/ApiKeyServices';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type {
  NewUserAccount,
  NewUserResponseData,
  User,
} from '@/features/auth/types/RegisterService';
import type {
  ProjectDetailDTO,
  ProjectResponse,
  ProjectsResponse,
} from '@/features/projects/types/projectService';
import type { UserTenantResponse } from '@/features/tenants/types/tenantServices';
import type { QuickViewDTO } from '@/hooks/useFetchQuickView/quickViewService';
export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<NewUserResponseData>;
  authenticateUser(user: Pick<User, 'email' | 'password'>): Promise<UserAuthResponse>;
  getTenantList(): Promise<UserTenantResponse>;
  createAPIKey(key: NewAPIKey): Promise<APIKey>;
  createTenant(): Promise<any>;
  projectDetail(id: ProjectDetailDTO): Promise<ProjectResponse>;
  getStats(values: QuickViewDTO): Promise<any>;
  getProjectList(): Promise<ProjectsResponse>;
}
