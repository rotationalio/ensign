import { t } from '@lingui/macro';
import invariant from 'invariant';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import useUserLoader from '@/features/members/loaders/userLoader';
import { isOnboardedMember } from '@/features/members/utils';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

const useDashOnboarding = () => {
  const { tenants, wasTenantsFetched } = useFetchTenants();
  const hasTenants =
    tenants?.tenants && Array.isArray(tenants?.tenants) && tenants?.tenants?.length > 0;
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
      invariant(
        hasTenants,
        t`Something went wrong. Please contact us at support@rotational.io for assistance.`
      );
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
