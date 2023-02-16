/* eslint-disable no-restricted-imports */

import { APIKey } from '@/features/apiKeys/types/apiKeyService';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type {
  NewUserAccount,
  NewUserResponseData,
  User,
} from '@/features/auth/types/RegisterService';
import { MembersResponse } from '@/features/members/types/memberServices';
import { OrgResponse } from '@/features/organizations/types/organizationService';
import type { ProjectResponse, ProjectsResponse } from '@/features/projects/types/projectService';
import type { UserTenantResponse } from '@/features/tenants/types/tenantServices';
import type { Topics } from '@/features/topics/types/topicService';
import type { QuickViewDTO } from '@/hooks/useFetchQuickView/quickViewService';
export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<NewUserResponseData>;
  authenticateUser(user: Pick<User, 'email' | 'password'>): Promise<UserAuthResponse>;
  getTenantList(): Promise<UserTenantResponse>;
  createAPIKey(): Promise<APIKey>;
  createTenant(): Promise<any>;
  projectDetail(projectID: string): Promise<ProjectResponse>;
  getStats(values: QuickViewDTO): Promise<any>;
  getTopics(): Promise<Topics | undefined>;
  getApiKeys: () => Promise<APIKey>;
  getProjectList(): Promise<ProjectsResponse>;
  getMemberList(): Promise<MembersResponse>;
  orgDetail(orgID: string): Promise<OrgResponse>;
}
