import { vi } from 'vitest';

import { projectsRequest } from '../projectListAPI';

describe('Project', () => {
  describe('Project List', () => {
    it('returns request resolved with response', async () => {
      const mockProjectResponse = {
        PromiseRejectionEvent: [
          {
            id: '1',
            name: 'project01',
          },
        ],
      };

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockProjectResponse,
        statusText: 'OK',
      });
      const request = projectsRequest(requestSpy);
      const response = await request('');
      expect(response).toBe(mockProjectResponse);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
