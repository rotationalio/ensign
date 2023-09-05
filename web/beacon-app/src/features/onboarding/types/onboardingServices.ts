import { UseMutateFunction } from '@tanstack/react-query';

import { MemberResponse } from '@/features/members/types/memberServices';

export interface MemberUpdateMutation {
  updateMember: UseMutateFunction<MemberResponse, unknown, UpdateMemberDTO, unknown>;
  reset(): void;
  member: MemberResponse;
  hasMemberFailed: boolean;
  wasMemberUpdated: boolean;
  isUpdatingMember: boolean;
  error: any;
}

export type UpdateMemberDTO = {
  memberID: string;
  onboardingPayload: Partial<MemberResponse>;
};
