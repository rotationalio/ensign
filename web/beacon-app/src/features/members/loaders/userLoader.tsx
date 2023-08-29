import { json } from 'react-router-dom';

import { useOrgStore } from '@/store';

import { useFetchMember } from '../hooks/useFetchMember';

const useUserLoader = () => {
  const orgDataState = useOrgStore.getState() as any;

  const { member } = useFetchMember(orgDataState?.user);

  return json({ userProfile: member });
};

export default useUserLoader;
