import { json } from 'react-router-dom';

import { queryClient } from '@/application/config/react-query';
import { useOrgStore } from '@/store';

import { memberDetailQuery } from '../hooks/useFetchMember';

const userLoader = () => async () => {
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState || null;

  if (!user) {
    return json({
      member: null,
    });
  }

  try {
    const query = memberDetailQuery(user);

    const member =
      queryClient.getQueryData(query?.queryKey) ?? (await queryClient.fetchQuery(query));
    return json({
      member,
    });
  } catch (error: any) {
    if (error?.response?.status === 401) {
      window.location.href = '/';
    }
  }
};

export default userLoader;
