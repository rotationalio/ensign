import { describe, expect, it, vi } from 'vitest';

import { updateMemberAPI } from '../updateMemberAPI';

describe('onboardingStepAPI', () => {
  it('returns request with response without developer segment', async () => {
    const mockOnboardingMember = {
      id: '1',
      organization: 'test',
      workspace: 'test-workspace',
      name: 'Clint Barton',
      profession_segment: 'Education',
    };

    const requestSpy = vi.fn().mockResolvedValueOnce({
      status: 200,
      data: mockOnboardingMember,
      statusText: 'OK',
    });
    const mockDTO = {
      memberID: '1',
      organization: 'test',
      workspace: 'test-workspace',
      name: 'Clint Barton',
      profession_segment: 'Education',
    };

    const request = updateMemberAPI(requestSpy);
    const response = await request(mockDTO);
    expect(response).toBe(mockOnboardingMember);
    expect(requestSpy).toHaveBeenCalledTimes(1);
  });

  it('returns request with response with developer segment', async () => {
    const mockOnboardingMember = {
      id: '1',
      organization: 'test',
      workspace: 'test-workspace',
      name: 'Clint Barton',
      profession_segment: 'Work',
    };

    const requestSpy = vi.fn().mockResolvedValueOnce({
      status: 200,
      data: mockOnboardingMember,
      statusText: 'OK',
    });
    const mockDTO = {
      memberID: '1',
      organization: 'test',
      workspace: 'test-workspace',
      name: 'Clint Barton',
      profession_segment: 'Work',
      developer_segment: 'DevOps',
    };

    const request = updateMemberAPI(requestSpy);
    const response = await request(mockDTO);
    expect(response).toBe(mockOnboardingMember);
    expect(requestSpy).toHaveBeenCalledTimes(1);
  });
});
