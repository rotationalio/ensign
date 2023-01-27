import { vi } from 'vitest';
import type { RegistrationMutation } from '../CreateAccountService';
import { isAccountCreated } from '../CreateAccountService';

describe('CreateAccountService', () => {
    describe('isAccountCreated', () => {
        it('should return true if the account was created', () => {
            const mutation = { wasAccountCreated: true, user: {} } as RegistrationMutation;
            const res = isAccountCreated(mutation);
            expect(res).toBe(true);
        });
        it('should return false if the account was not created', () => {
            const mutation = { wasAccountCreated: false, user: {} } as RegistrationMutation;
            const res = isAccountCreated(mutation);
            expect(res).toBe(false);
        });
    });
});