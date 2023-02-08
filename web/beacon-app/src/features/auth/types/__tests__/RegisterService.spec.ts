import type { NewUserAccount } from '../RegisterService';
import { hasUserRequiredFields } from '../RegisterService';

describe('RegisterService types checker', () => {
  describe('hasUserRequiredFields', () => {
    it('should return true if all fields are filled', () => {
      const user: NewUserAccount = {
        email: 'test@gmail.com',
        name: 'test',
        password: 'test',
        organization: 'test',
        domain: 'test',
        pwcheck: 'test',
        terms_agreement: true,
        privacy_agreement: true,
      };
      const res = hasUserRequiredFields(user);
      expect(res).toBe(true);
    });
    it('should return false account is missing fields', () => {
      const user = {
        email: 'test@gmail.com',
        name: 'test',
        password: 'test',
        organization: 'test',
        domain: 'test',
        // pwcheck: 'test',
        terms_agreement: true,
      } as any;
      const res = hasUserRequiredFields(user);
      expect(res).toBe(false);
    });
  });
});
