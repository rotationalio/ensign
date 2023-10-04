import { vi } from 'vitest';

import { APP_ROUTE } from '@/constants';

import { forgotPasswordRequest } from '../ForgotPasswordApiService';

// Mocking AxiosResponse
const mockResponse = {
  data: {},
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {},
};

// Mocking request function
const mockRequest = vi.fn();

// Mocking payload for forgotPasswordRequest
const mockPayload = {
  email: 'test@rotational.io',
};

// Mocking getValidApiResponse
vi.mock('@/application/api/ApiService', () => ({
  getValidApiResponse: (response: any) => response,
}));

describe('forgotPasswordRequest', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should send a POST request to forgot password', async () => {
    mockRequest.mockResolvedValueOnce(mockResponse);

    const apiAdapter = forgotPasswordRequest(mockRequest);

    const result = await apiAdapter(mockPayload);

    expect(mockRequest).toHaveBeenCalledWith(`${APP_ROUTE.FORGOT_PASSWORD}`, {
      method: 'POST',
      data: JSON.stringify(mockPayload),
    });
    expect(result).toEqual(mockResponse);
  });

  it('should handle API error response', async () => {
    const errorResponse = {
      ...mockResponse,
      status: 400,
      data: { error: 'Bad Request' },
    };
    mockRequest.mockRejectedValueOnce(errorResponse);

    const apiAdapter = forgotPasswordRequest(mockRequest);

    await expect(apiAdapter(mockPayload)).rejects.toEqual(errorResponse);
  });

  it('should handle non-200 status codes', async () => {
    const non200Response = {
      ...mockResponse,
      status: 201,
      data: { message: 'Created, but not the expected response' },
    };
    mockRequest.mockResolvedValueOnce(non200Response);

    const apiAdapter = forgotPasswordRequest(mockRequest);

    const result = await apiAdapter(mockPayload);
    expect(result).not.toEqual(mockResponse);
    expect(result).toEqual(non200Response);
  });

  it('should handle unexpected data in the response', async () => {
    const unexpectedDataResponse = {
      ...mockResponse,
      data: { unexpected: 'Unexpected data' },
    };
    mockRequest.mockResolvedValueOnce(unexpectedDataResponse);

    const apiAdapter = forgotPasswordRequest(mockRequest);

    const result = await apiAdapter(mockPayload);
    expect(result).not.toEqual(mockResponse);
    expect(result).toEqual(unexpectedDataResponse);
  });
});
