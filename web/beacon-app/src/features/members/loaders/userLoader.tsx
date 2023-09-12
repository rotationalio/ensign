import { useOrgStore } from '@/store';

import { useFetchMember } from '../hooks/useFetchMember';

const useUserLoader = () => {
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState;

  const { member } = useFetchMember(user);

  // console.log('[] member', member);

  return {
    member,
  };
};

export default useUserLoader;
