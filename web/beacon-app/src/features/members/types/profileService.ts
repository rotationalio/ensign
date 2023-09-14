import { UseMutateFunction } from '@tanstack/react-query';

import { MemberResponse, NewMemberDTO, UpdateMemberDTO } from './memberServices';

export interface ProfileMutation {
  profileMutation: UseMutateFunction<MemberResponse, unknown, NewMemberDTO, unknown>;
  reset(): void;
  profile: any;
  hasProfileFailed: boolean;
  wasProfileCreated: boolean;
  isCreatingProfile: boolean;
  error: any;
}

export interface ProfileQuery {
  getProfile(): void;
  profile: any;
  hasProfileFailed: boolean;
  wasProfileFetched: boolean;
  isFetchingProfile: boolean;
  error: any;
}

export interface ProfileUpdateMutation {
  updateProfile: UseMutateFunction<MemberResponse, unknown, UpdateMemberDTO, unknown>;
  reset(): void;
  profile: MemberResponse;
  hasProfileFailed: boolean;
  wasProfileUpdated: boolean;
  isUpdatingProfile: boolean;
  error: any;
}
