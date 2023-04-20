import { UserAuthResponse } from '../../auth/types/LoginService';

export interface SwitchMutation {
  switch: (orgId: string) => void;
  hasSwitchFailed: boolean;
  wasSwitchFetched: boolean;
  isSwitching: boolean;
  auth: UserAuthResponse;
  error: any;
  reset: () => void;
}

export interface SwitchDTO {
  org_id: string;
}
