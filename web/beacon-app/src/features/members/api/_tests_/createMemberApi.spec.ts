import { vi } from 'vitest';

import { createMemberRequest } from '../createMemberApiService';

describe('Member', () => {
  describe('Create Member', () => {
    it('returns request with response', async () => {
      const mockMember = {
        id: '1',
        name: 'Kamala Khan',
        role: 'Owner',
        created: '23.02.01',
        modified: '23.02.01',
      };
      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockMember,
        statusText: 'OK',
      });
      const mockDTO = {
        name: 'Kamala Khan',
        role: 'Owner',
      };

      const request = createMemberRequest(requestSpy);
      const response = await request(mockDTO);
      expect(response).toBe(mockMember);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });

    it('returns request with error', async () => {
      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: undefined,
        statusText: 'OK',
      });
      const mockDTO = {
        name: 'Kamala Khan',
      };

      const request = createMemberRequest(requestSpy);
      const response = await request(mockDTO);
      expect(response).toBe(undefined);
      expect(requestSpy).toHaveBeenCalledTimes(1);
    });
  });
});
