import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { QuickViewDTO, QuickViewResponse } from './quickViewService';

const statsRequest =
  (request: Request): ApiAdapters['getStats'] =>
  async ({ id, key }: QuickViewDTO) => {
    invariant(id, 'id is required');
    invariant(key, 'key is required');

    const response = (await request(`${APP_ROUTE.ROOT}/${key}/${id}/stats`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<QuickViewResponse>(response);
  };

export default statsRequest;
