import { vi } from 'vitest';

import statsRequest from '../quickViewApiService';
import { QuickViewDTO } from '../quickViewService';
vi.mock('invariant');
// vi.mock('axios');

describe('QuickViewApiService', () => {
  describe('statsRequest', () => {
    it('returns request resolved with response', async () => {
      const mockStats = {
        id: '1',
        key: 'tenant',
      } satisfies QuickViewDTO;

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockStats,
        statusText: 'OK',
      });

      const request = statsRequest(requestSpy);
      const response = await request(mockStats);
      expect(response).toBe(mockStats);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      // should return request payload
      expect(requestSpy).toHaveBeenCalledWith(`/${mockStats.key}/${mockStats.id}/stats`, {
        method: 'GET',
      });
    });
    it('throws error when required fields are missing', async () => {
      const mockStats = {
        id: '1',
      } as any;

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockStats,
        statusText: 'OK',
      });
      const request = statsRequest(requestSpy);

      expect(requestSpy).toHaveBeenCalledTimes(0);
      expect(Promise.reject(request(mockStats))).rejects.toThrowError();
      // should return request payload
      expect(requestSpy).toHaveBeenCalledWith(`/${mockStats.key}/${mockStats.id}/stats`, {
        method: 'GET',
      });
    });
  });
});
