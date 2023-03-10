import { vi } from 'vitest';

import statsRequest from '../quickViewApiService';
vi.mock('invariant');
// vi.mock('axios');

describe('QuickViewApiService', () => {
  describe('statsRequest', () => {
    it('returns request resolved with response', async () => {
      const mockResponse = {
        name: 'test',
        value: 1,
        units: 'test',
        percent: 1,
      };

      const mockTenantID = '1';

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: [mockResponse],
        statusText: 'OK',
      });

      const request = statsRequest(requestSpy);
      const response = await request(mockTenantID);
      expect(response).toStrictEqual([mockResponse]);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      // should return request payload
      expect(requestSpy).toHaveBeenCalledWith(`/tenant/${mockTenantID}/stats`, {
        method: 'GET',
      });
    });
  });
});
