import { json } from 'react-router-dom';

import { queryClient } from '@/application/config/react-query';
import { useOrgStore } from '@/store';
import { clearCookies } from '@/utils/cookies';

import { memberDetailQuery } from '../hooks/useFetchMember';

const userLoader = () => async () => {
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState || null;
  // console.log('userLoader', user);

  if (!user) {
    return null;
  } else {
    try {
      const query = memberDetailQuery(user);

      const member =
        queryClient.getQueryData(query?.queryKey) ?? (await queryClient.fetchQuery(query));
      return json({
        member,
      });
    } catch (error: any) {
      if (error?.response?.status === 401) {
        clearCookies();
        // routeLoader.clearCacheAll(); // TODO: clear cache all routes when logout
        window.location.href = '/';
      }
    }
  }
};

export default userLoader;
