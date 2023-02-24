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

export interface MembersQuery {
  getMembers(): void;
  members: any;
  hasMembersFailed: boolean;
  wasMembersFetched: boolean;
  isFetchingMembers: boolean;
  error: any;
}
