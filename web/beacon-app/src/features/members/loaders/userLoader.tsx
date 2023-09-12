import { useOrgStore } from '@/store';

import { useFetchMember } from '../hooks/useFetchMember';

const useFetchCurrentMember = () => {
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState;

  const { member, isFetchingMember, wasMemberFetched } = useFetchMember(user);

  // console.log('[] member', member);

  return {
    member,
    isMemberLoading: isFetchingMember,
    wasMemberFetched,
  };
};

export default useFetchCurrentMember;
