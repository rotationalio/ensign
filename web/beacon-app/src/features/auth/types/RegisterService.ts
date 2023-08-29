export interface User {
  user_id: string;
  name: string;
  pwcheck: string;
  organization: string;
  domain?: string;
  terms_agreement?: boolean;
  privacy_agreement: boolean;
  email: string;
  invite_token?: string;
  password: string;
}

export interface NewUserResponseData {
  user_id: string;
  org_id: string;
  email: string;
  message: string;
  role: string;
  created: string;
}

export type NewUserAccount = Omit<
  User,
  'user_id' | 'name' | 'organization' | 'domain' | 'terms_agreement' | 'privacy_agreement'
>;

export type NewInvitedUserAccount = Omit<User, 'user_id' | 'organization' | 'domain'>;

export const hasUserRequiredFields = (account: NewUserAccount | NewInvitedUserAccount): boolean => {
  if (!account.email || !account.password || !account.pwcheck) {
    return false;
  }
  return true;
}