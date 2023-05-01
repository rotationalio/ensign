import { json } from 'react-router-dom';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { tenantsRequest } from '@/features/tenants/api/tenantListAPI';

import { projectStatsApiRequest } from '../api/projectStatsApiService';

export async function projectPageLoader() {
  const tenants = await queryClient.fetchQuery({
    queryFn: tenantsRequest(axiosInstance),
  });

  const projectStats = await projectStatsApiRequest(axiosInstance)(tenants.tenants[0].id);

  return json({
    projectStats,
  });
}
