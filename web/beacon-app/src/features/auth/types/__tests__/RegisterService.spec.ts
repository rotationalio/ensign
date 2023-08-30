import type { NewUserAccount } from '../RegisterService';
import { hasUserMissingFields, hasUserRequiredFields } from '../RegisterService';

describe('RegisterService types checker', () => {
  describe('hasUserRequiredFields', () => {
    it('should return true if all fields are filled', () => {
      const user: NewUserAccount = {
        email: 'test@gmail.com',
        password: 'test',
        pwcheck: 'test',
      };
      const res = hasUserRequiredFields(user);
      expect(res).toBe(true);
    });
    it('should return false account is missing fields', () => {
      const user = {
        email: 'test@gmail.com',
        password: 'test',
        // pwcheck: 'test',
      } as any;
      const res = hasUserMissingFields(user);
      expect(res).toBe(false);
    });
  });
});
