export interface MembersResponse {
  member: MemberResponse[];
  prev_page_token: string;
  next_page_token: string;
}

export interface MemberResponse {
  id: string;
  name: string;
  role: string;
}

export interface MemberQuery {
  getMembers(): void;
  member: any;
  hasMemberFailed: boolean;
  wasMemberFetched: boolean;
  isFetchingMember: boolean;
  error: any;
}
