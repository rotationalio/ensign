import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';

import type { QuickViewResponse } from './quickViewService';

const statsRequest =
  (request: Request): ApiAdapters['getStats'] =>
  async (tenantID: string) => {
    invariant(tenantID, 'id is required');
    const link = `/tenant/${tenantID}/stats`;
    try {
      const response = (await request(`${link}`, {
        method: 'GET',
      })) as any;

      return getValidApiResponse<QuickViewResponse>(response);
    } catch (e: any) {
      getValidApiError(e);
    }
  };

export default statsRequest;
