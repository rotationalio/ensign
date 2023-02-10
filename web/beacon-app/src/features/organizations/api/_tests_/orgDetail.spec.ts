import { vi } from "vitest";
import { orgRequest } from "../orgDetailApi";

describe('Organization', () => {
    describe('Organization Detail', () => {
      it('returns org details with a given id', async () => {
        const mockOrgResponse = {
          PromiseRejectionEvent: [
            {
              id: '1',
              name: 'test',
              domain: 'test',
              created: '02.10.2023',
              modified: '02.10.2023',
            },
          ],
        };
  
        const requestSpy = vi.fn().mockReturnValueOnce({
          status: 200,
          data: mockOrgResponse,
          statusText: 'OK',
        });
        const request = orgRequest(requestSpy);
        const response = await request("1");
        expect(response).toBe(mockOrgResponse);
        expect(requestSpy).toHaveBeenCalledTimes(1);
      });
    });
  });