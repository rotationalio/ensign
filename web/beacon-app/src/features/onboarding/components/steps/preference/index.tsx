import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useUpdateProfile } from '@/features/members/hooks/useUpdateProfile';

// import { useOrgStore } from '@/store';
import { getOnboardingStepsData } from '../../../shared/utils';
import StepCounter from '../StepCounter';
import UserPreferenceStepForm from './form';

const UserPreferenceStep = () => {
  const navigate = useNavigate();
  // const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { profile } = useFetchProfile();
  const { wasProfileUpdated, isUpdatingProfile, error, updateProfile, hasProfileFailed } =
    useUpdateProfile();
  const hasError = error && error.response.status === 400;

  const submitFormHandler = (values: any) => {
    const requestPayload = {
      memberID: profile?.id,
      payload: {
        ...getOnboardingStepsData(profile),
        developer_segment: values?.developer_segment?.map((item: any) => item.value),
        profession_segment: values?.profession_segment,
      },
    };
    updateProfile(requestPayload);
  };

  useEffect(() => {
    if (wasProfileUpdated || !hasProfileFailed) {
      navigate(PATH_DASHBOARD.HOME);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [wasProfileUpdated, profile]);

  return (
    <>
      <StepCounter />

      <div className="flex flex-col justify-center ">
        <UserPreferenceStepForm
          onSubmit={submitFormHandler}
          isSubmitting={isUpdatingProfile}
          error={error}
          initialValues={{
            developer_segment: profile?.developer_segment?.map((item: any) => ({
              label: item,
              value: item,
            })),
            profession_segment: profile.profession_segment,
          }}
          hasError={hasError}
        />
      </div>
    </>
  );
};

export default UserPreferenceStep;
