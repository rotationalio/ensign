import { UseMutateFunction } from '@tanstack/react-query';

export interface MembersResponse {
  member: MemberResponse[];
  prev_page_token: string;
  next_page_token: string;
}

export interface MemberResponse {
  id: string;
  name: string;
  role: string;
  created: string;
  modified: string;
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

export interface MembersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error: any;
}
export type NewMemberDTO = Pick<MemberResponse, 'name' | 'role'>;

export const hasMemberRequiredFields = (member: NewMemberDTO): member is Required<NewMemberDTO> => {
  return Object.values(member).every((x) => !!x);
};
