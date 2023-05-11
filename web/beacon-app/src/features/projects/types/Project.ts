import { MemberResponse } from '@/features/members/types/memberServices';
export type Project = {
  created: string;
  id: string;
  modified: string;
  name: string;
  tenant_id: string;
  description?: string;
  status?: string;
  owner: Partial<MemberResponse>;
};

export enum ProjectStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  DELETED = 'DELETED',
  INCOMPLETE = 'Incomplete',
  COMPLETE = 'Complete',
}
