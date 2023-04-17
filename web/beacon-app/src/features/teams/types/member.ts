import { MEMBER_ROLE, MEMBER_STATUS } from '@/constants/rolesAndStatus';

export type Member = {
  id: string;
  email: string;
  name: string;
  role: MemberRole;
  status: MemberStatus;
  created: string;
  modified: string;
};

export type MemberRole = keyof typeof MEMBER_ROLE;
export type MemberStatus = keyof typeof MEMBER_STATUS;
