import invariant from 'invariant';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import useUserLoader from '@/features/members/loaders/userLoader';
import { isOnboardedMember } from '@/features/members/utils';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import ErrorMessage from '@/utils/error-message';

const useDashOnboarding = () => {
  const { tenants, wasTenantsFetched } = useFetchTenants();
  const hasTenants = tenants?.tenants?.length > 0;
  const navigate = useNavigate();
  const { member: loaderData, wasMemberFetched, isMemberLoading } = useUserLoader();
  const isOnboarded = isOnboardedMember(loaderData?.onboarding_status);

  useEffect(() => {
    if (!isOnboarded) {
      navigate(PATH_DASHBOARD.ONBOARDING);
    }
  }, [isOnboarded, navigate]);

  useEffect(() => {
    if (wasTenantsFetched) {
      invariant(hasTenants, ErrorMessage.somethingWentWrong);
    }
  }, [hasTenants, wasTenantsFetched]);

  return {
    isOnboarded,
    isMemberLoading,
    loaderData,
    wasMemberFetched,
    wasTenantsFetched,
    hasTenants,
    tenants,
  };
};

export default useDashOnboarding;
