import { vi } from 'vitest';

import { createProjectAPI } from '../createProjectAPI';

describe('createProjectAPI', () => {
  it('returns request with response', async () => {
    const mockProject = {
      id: '1',
      name: 'project01',
    };

    const requestSpy = vi.fn().mockReturnValueOnce({
      status: 200,
      data: mockProject,
      statusText: 'OK',
    });
    const mockDTO = {
      tenantID: '1',
      name: 'project01',
    };

    const request = createProjectAPI(requestSpy);
    const response = await request(mockDTO);
    expect(response).toBe(mockProject);
    expect(requestSpy).toHaveBeenCalledTimes(1);
  });
});
