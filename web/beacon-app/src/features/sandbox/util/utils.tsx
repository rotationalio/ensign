export const ACCOUNT_TYPE = {
  SANDBOX: 'sandbox',
};

export const isSandboxAccount = (account: string) => {
  return account === ACCOUNT_TYPE.SANDBOX;
};
