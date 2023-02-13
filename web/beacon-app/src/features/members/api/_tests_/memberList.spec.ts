import { vi } from 'vitest';

import { memberRequest } from '../memberListAPI';

describe('Member', () => {
  describe('Member List', () => {
    it('returns member list request resolved with response', async () => {
      const mockMembersResponse = {
        PromiseRejectionEvent: [
          {
            id: '1',
            name: 'Kamala Khan',
            role: 'Owner',
          },
        ],
      };

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockMembersResponse,
        statusText: 'OK',
      });
      const request = memberRequest(requestSpy);
      const response = await request();
      expect(response).toBe(mockMembersResponse);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
