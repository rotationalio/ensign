import invariant from 'invariant';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { isOnboardedMember } from '@/features/members/utils';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import ErrorMessage from '@/utils/error-message';
import { cleanCookiesOnDashboard } from '@/utils/misc';
const useDashOnboarding = () => {
  const { tenants, wasTenantsFetched } = useFetchTenants();
  const hasTenants = tenants?.tenants?.length > 0;
  const navigate = useNavigate();
  const { profile: loaderData, wasProfileFetched, isFetchingProfile } = useFetchProfile();
  const isOnboarded = isOnboardedMember(loaderData?.onboarding_status);

  useEffect(() => {
    if (!isOnboarded) {
      navigate(PATH_DASHBOARD.ONBOARDING);
    }
  }, [isOnboarded, navigate]);

  useEffect(() => {
    if (wasTenantsFetched) {
      invariant(hasTenants, ErrorMessage.somethingWentWrong);
      cleanCookiesOnDashboard();
    }
  }, [hasTenants, wasTenantsFetched]);

  return {
    isOnboarded,
    isFetchingProfile,
    loaderData,
    wasProfileFetched,
    wasTenantsFetched,
    hasTenants,
    tenants,
  };
};

export default useDashOnboarding;
