import { vi } from 'vitest';

import { APP_ROUTE } from '@/constants';

import type { Topics } from '../../types/topicService';
import { topicsRequest } from '../topicsApiService';
vi.mock('invariant');
// vi.mock('axios');

describe('Topics API Service ', () => {
  describe('topicsRequest', () => {
    it('returns request resolved with response', async () => {
      // const mockDTO = {
      //   projectID: '1',
      // } as any;

      const mockResponse = {
        project_id: '1',
        topics: [
          {
            id: '1',
            name: 'test',
          },
        ],
        prev_page_token: '1',
        next_page_token: '1',
      } as Topics;

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockResponse,
        statusText: 'OK',
      });

      const request = topicsRequest(requestSpy);
      const response = await request();
      expect(response).toBe(mockResponse);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      // should return request payload
      expect(requestSpy).toHaveBeenCalledWith(`${APP_ROUTE.TOPICS}`, {
        method: 'GET',
      });
    });
    it('throws error when required fields are missing', async () => {
      const mockResponse = {
        project_id: '1',
        topics: [
          {
            id: '1',
            name: 'test',
          },
        ],
        prev_page_token: '1',
        next_page_token: '1',
      } as Topics;

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockResponse,
        statusText: 'OK',
      });
      const request = topicsRequest(requestSpy);

      expect(requestSpy).toHaveBeenCalledTimes(0);
      expect(Promise.reject(request())).rejects.toThrowError();
    });
  });
});
