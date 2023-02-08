import { vi } from 'vitest';

import { createAPIKey } from '../CreateApiKey';

describe('CreateAPIKeyService', () => {
  describe('createAPIKey', () => {
    it('returns request with response', async () => {
      const mockKey = {
        client_id: '1',
        client_secret: 'secret',
        name: 'test',
        owner: 'example',
        permissions: ['write'],
        created: '23.02.01',
        modified: '23.02.01',
      };

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockKey,
        statusText: 'OK',
      });
      const request = createAPIKey(requestSpy);
      const response = await request(mockKey);
      expect(response).toBe(mockKey);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
