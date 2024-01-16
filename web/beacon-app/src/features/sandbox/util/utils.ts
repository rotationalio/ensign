import { ACCOUNT_TYPE } from '../types/accountType';

export const isSandboxAccount = (account: string) => {
  return account === ACCOUNT_TYPE.SANDBOX;
};
