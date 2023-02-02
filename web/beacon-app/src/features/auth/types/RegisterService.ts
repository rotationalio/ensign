export interface User {
  user_id: string;
  name: string;
  pwcheck: string;
  organization: string;
  domain: string;
  terms_agreement: boolean;
  privacy_agreement: boolean;
  email: string;
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


export type NewUserAccount = Omit<User, 'user_id'>;

export const hasUserRequiredFields = (account: NewUserAccount): account is Required<NewUserAccount> => {
  return (
    account.name !== undefined && account.name !== '' && account.email !== undefined && account.email !== '' && account.password !== undefined && account.password !== '' && account.pwcheck !== undefined && account.pwcheck !== '' && account.organization !== undefined && account.organization !== '' && account.domain !== undefined && account.domain !== '' && account.terms_agreement !== undefined && account.terms_agreement !== false && account.privacy_agreement !== undefined && account.privacy_agreement !== false
  );
};

