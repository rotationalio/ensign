import { UseMutateFunction } from '@tanstack/react-query';

import { MEMBER_ROLE, MEMBER_STATUS } from '@/constants/rolesAndStatus';
export interface MembersResponse {
  member: MemberResponse[];
  prev_page_token: string;
  next_page_token: string;
}

export interface MemberResponse {
  id: string;
  name: string;
  role: MemberRole;
  workspace: string;
  profession_segment: string;
  developer_segment: string;
  organization: string;
  email: string;
  status: MemberStatus; // TODO: remove this once the new enpoint is ready
  onboarding_status: string;
  created: string;
  picture: string;
  invited?: boolean;
  date_added?: string;
  last_activity: string;
}

export interface MemberQuery {
  getMember(): void;
  member: any;
  hasMemberFailed: boolean;
  wasMemberFetched: boolean;
  isFetchingMember: boolean;
  error: any;
}

export interface MemberMutation {
  createMember: UseMutateFunction<MemberResponse, unknown, NewMemberDTO, unknown>;
  reset(): void;
  member: any;
  hasMemberFailed: boolean;
  wasMemberCreated: boolean;
  isCreatingMember: boolean;
  error: any;
}

export interface MemberDeleteMutation {
  deleteMember: UseMutateFunction<unknown, unknown, void, unknown>;
  reset(): void;
  member: any;
  hasMemberFailed: boolean;
  wasMemberDeleted: boolean;
  isDeletingMember: boolean;
  error: any;
}

export interface MembersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error: any;
}
export type NewMemberDTO = Pick<MemberResponse, 'email' | 'role'>;
export type DeleteMemberDTO = Pick<MemberResponse, 'id'>;

export const hasMemberRequiredFields = (member: NewMemberDTO): member is Required<NewMemberDTO> => {
  return Object.values(member).every((x) => !!x);
};

export type MemberRole = keyof typeof MEMBER_ROLE;
export type MemberStatus = keyof typeof MEMBER_STATUS;
