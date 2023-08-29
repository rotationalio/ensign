import { json } from 'react-router-dom';

import { queryClient } from '@/application/config/react-query';
import { useOrgStore } from '@/store';

import { memberDetailQuery } from '../hooks/useFetchMember';

const userLoader = () => async () => {
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState || null;

  const query = memberDetailQuery(user);

  const member = queryClient.getQueryData(query.queryKey) ?? (await queryClient.fetchQuery(query));
  return json({
    member,
  });
};

export default userLoader;
