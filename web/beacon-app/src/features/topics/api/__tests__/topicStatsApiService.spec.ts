import { vi } from 'vitest';

import { APP_ROUTE } from '@/constants';

import { getTopicsStatsMockData } from '../../__mocks__';
import type { TopicQuickViewResponse } from '../../types/topicService';
import topicStatsApiRequest from '../topicStatsApiService';
vi.mock('invariant');
// vi.mock('axios');
vi.mock('invariant');
describe('Topics Stats API Service ', () => {
  describe('topicsStatsRequest', () => {
    it('returns request resolved with response', async () => {
      const mockDTO = {
        topicID: '1',
      } as any;

      const mockResponse = {
        data: getTopicsStatsMockData(),
      } as TopicQuickViewResponse;

      const requestSpy = vi.fn().mockReturnValueOnce({
        status: 200,
        data: mockResponse,
        statusText: 'OK',
      });

      const request = topicStatsApiRequest(requestSpy);
      const response = await request(mockDTO.topicID);
      expect(response).toBe(mockResponse);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      // should return request payload
      expect(requestSpy).toHaveBeenCalledWith(`${APP_ROUTE.TOPICS}/${mockDTO.topicID}/stats`, {
        method: 'GET',
      });
    });
  });
});
