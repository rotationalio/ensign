import { vi } from 'vitest';

import { APP_ROUTE } from '@/constants';

import { deleteAPIKeyRequest } from '../deleteApiKeyApi';
describe('DeleteAPIKeyService', () => {
  describe('deleteAPIKey', () => {
    let requestSpy;

    beforeEach(() => {
      // Mock the request function to return a successful response
      requestSpy = vi.fn().mockResolvedValueOnce({
        status: 204,
        data: {},
        statusText: 'OK',
      });
    });

    it('should call the request function with the correct API key', async () => {
      const mockDTO = {
        apiKey: '1',
      };

      // Act
      const deleteAPIKey = deleteAPIKeyRequest(requestSpy);
      const response = await deleteAPIKey(mockDTO.apiKey);

      expect(response).toBe(requestSpy.mock.results[0].value.data);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      expect(requestSpy).toHaveBeenCalledWith(`${APP_ROUTE.APIKEYS}/1`, {
        method: 'DELETE',
      });
    });

    it('should return a 204 status code response', async () => {
      // Arrange
      const mockDTO = {
        apiKey: '1',
      };

      // Act
      const deleteAPIKey = deleteAPIKeyRequest(requestSpy);
      const response = await deleteAPIKey(mockDTO.apiKey);

      // Assert
      expect(response.status).toBe(requestSpy.mock.results[0].value.data.status);
    });
  });
});
