import { vi } from 'vitest';

import { tenantRequest } from '../tenantListAPI';

describe('Tenant', () => {
  describe('Tenant Services', () => {
    it('returns request resolved with response', async () => {
      const mockTenantResponse = {
         tenant : [
          {
            id: "1",
            name: "Test Tenant",
            environment_type: "Prod",
          }
         ]
      }

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockTenantResponse,
        statusText: 'OK',
      });
      const request = tenantRequest(requestSpy);
      const response = await request();
      expect(response).toBe(mockTenantResponse);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
