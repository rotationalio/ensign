import { vi } from 'vitest';

import { APP_ROUTE } from '@/constants';

import { deleteAPIKeyRequest } from '../deleteApiKeyApi';
describe('DeleteAPIKeyService', () => {
  describe('deleteAPIKey', () => {
    let requestSpy;

    beforeEach(() => {
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

      expect(response).toStrictEqual(requestSpy.mock.results[0].value.data);
      expect(requestSpy).toHaveBeenCalledTimes(1);
      expect(requestSpy).toHaveBeenCalledWith(`${APP_ROUTE.APIKEYS}/${mockDTO.apiKey}`, {
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
      expect(response.status).toStrictEqual(requestSpy.mock.results[0].value.data.status);
    });
  });
});
