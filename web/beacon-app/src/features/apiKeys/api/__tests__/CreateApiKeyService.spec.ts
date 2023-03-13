import { vi } from 'vitest';

import { createProjectAPIKey } from '../createApiKey';

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
      const mockDTO = {
        projectID: '1',
        name: 'test',
        permissions: ['write'],
      };

      const request = createProjectAPIKey(requestSpy);
      const response = await request(mockDTO);
      expect(response).toBe(mockKey);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
