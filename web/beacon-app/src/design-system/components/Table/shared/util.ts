export enum STATUS {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  PENDING = 'pending',
  ERROR = 'error',
  CONFIRMED = 'confirmed',
  UNUSED = 'unused',
  REVOKED = 'revoked',
  COMPLETE = 'complete',
  INCOMPLETE = 'incomplete',
}

export type TStatus = 'Active' | 'Inactive' | 'Pending' | 'Error';

export const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};
