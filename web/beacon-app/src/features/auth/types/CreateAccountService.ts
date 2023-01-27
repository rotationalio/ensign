import type { NewUserAccount, User } from '@/features/auth/types/RegisterService';

export interface RegistrationMutation {
    createNewAccount(user: NewUserAccount): void;
    reset(): void;
    user: User;
    hasAccountFailed: boolean;
    wasAccountCreated: boolean;
    isCreatingAccount: boolean;
}

export const isAccountCreated = (mutation: RegistrationMutation): mutation is Required<RegistrationMutation> =>
    mutation.wasAccountCreated && mutation.user !== undefined;
