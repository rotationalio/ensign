import { vi } from 'vitest';

import { createAccountRequest } from '../RegisterApiService';

describe('AccountApiService', () => {
  describe('createAccountRequest', () => {
    it('returns request resolved with response', async () => {
      const mockAccount = {
        name: 'name',
        organization: 'organization',
        domain: 'domain',
        terms_agreement: true,
        privacy_agreement: true,
        email: 'test@rotational.io',
        password: 'password',
        pwcheck: 'password',

      };

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockAccount,
        statusText: 'OK',
      });
      const request = createAccountRequest(requestSpy);
      const response = await request(mockAccount);
      expect(response).toBe(mockAccount);
      expect(requestSpy).toHaveBeenCalledTimes(1);

      // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/register', {
      //     method: "POST",
      //     headers: {
      //         "Content-Type": "application/json",
      //     },
      //     body: JSON.stringify(mockAccount),
      // });

      // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/register');
      // expect(requestSpy).toHaveBeenCalledWith(expect.stringContaining('localhost:8088'));
    });
  });
});
