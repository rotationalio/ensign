import { vi } from 'vitest';

import { forgotPasswordRequest } from '../ForgotPasswordApiService';

describe('ForgotPasswordApiService', () => {
  describe('forgotPasswordRequest', () => {
    it('returns request resolved with response', async () => {
      const mockAccount = {
        email: 'test@test.com',
      };

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockAccount,
        statusText: 'OK',
      });

      const request = forgotPasswordRequest(requestSpy);
      const response = await request(mockAccount);
      expect(response).toBe(mockAccount);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });

  // todo: add more failing tests
});
