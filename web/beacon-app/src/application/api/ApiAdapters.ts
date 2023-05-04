/* eslint-disable no-restricted-imports */

import { APIKey } from '@/features/apiKeys/types/apiKeyService';
import { APIKeyDTO } from '@/features/apiKeys/types/createApiKeyService';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';
import type {
  NewUserAccount,
  NewUserResponseData,
  User,
} from '@/features/auth/types/RegisterService';
import {
  MemberResponse,
  MembersResponse,
  NewMemberDTO,
} from '@/features/members/types/memberServices';
import { OrgListResponse, OrgResponse } from '@/features/organization/types/organizationService';
import { NewProjectDTO } from '@/features/projects/types/createProjectService';
import type { ProjectResponse, ProjectsResponse } from '@/features/projects/types/projectService';
import type { UserTenantResponse } from '@/features/tenants/types/tenantServices';
import type { NewTopicDTO, Topic, Topics } from '@/features/topics/types/topicService';
export interface ApiAdapters {
  createNewAccount(user: NewUserAccount): Promise<NewUserResponseData>;
  authenticateUser(
    user: Pick<User, 'email' | 'password' | 'invite_token'>
  ): Promise<UserAuthResponse>;
  getTenantList(): Promise<UserTenantResponse>;
  createProjectAPIKey(payload: APIKeyDTO): Promise<APIKey>;
  createTenant(): Promise<any>;
  projectDetail(projectID: string): Promise<ProjectResponse>;
  getStats(tenantID: string): Promise<any>;
  getTopics(): Promise<Topics | undefined>;
  createProjectTopic(payload: NewTopicDTO): Promise<Topic>;
  getApiKeys: (projectID: string) => Promise<APIKey>;
  getProjectList(tenantID: string): Promise<ProjectsResponse>;
  getMemberList(): Promise<MembersResponse>;
  getMemberDetail(memberID: string): Promise<MembersResponse>;
  orgDetail(orgID: string): Promise<OrgResponse>;
  checkToken(token: string): Promise<any>;
  getPermissions(): Promise<any>;
  getInviteTeamMember(token: string): Promise<any>;
  createMember(member: NewMemberDTO): Promise<MemberResponse>;
  updateMemberRole(memberId: string, data: any): Promise<any>;
  deleteMember(memberId: string): Promise<any>;
  getOrganizationList(): Promise<OrgListResponse>;
  switchOrganization(orgID: string): Promise<UserAuthResponse>;
  createNewProject(payload: NewProjectDTO): Promise<ProjectResponse>;
  getProjectStats(tenantID: string): Promise<any>;
}
