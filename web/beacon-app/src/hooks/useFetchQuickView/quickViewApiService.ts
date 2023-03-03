import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';

import type { QuickViewDTO, QuickViewResponse } from './quickViewService';

const statsRequest =
  (request: Request): ApiAdapters['getStats'] =>
  async ({ id, key }: QuickViewDTO) => {
    invariant(id, 'id is required');
    invariant(key, 'key is required');
    const link = `/stats/${key}/${id}`;
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
