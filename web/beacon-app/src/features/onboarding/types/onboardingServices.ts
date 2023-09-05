import { UseMutateFunction } from '@tanstack/react-query';

import { MemberResponse } from '@/features/members/types/memberServices';

export interface OnboardingMemberUpdateMutation {
  updateMember: UseMutateFunction<MemberResponse, unknown, UpdateMemberOnboardingDTO, unknown>;
  reset(): void;
  member: MemberResponse;
  hasMemberFailed: boolean;
  wasMemberUpdated: boolean;
  isUpdatingMember: boolean;
  error: any;
}

export type UpdateMemberOnboardingDTO = {
  memberID: string;
  onboardingPayload: Pick<
    MemberResponse,
    'organization' | 'workspace' | 'name' | 'profession_segment' | 'developer_segment'
  >;
};
