import { t } from '@lingui/macro';
import { useEffect } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import { getOnboardingStepsData, hasCompletedOnboarding } from '../../../shared/utils';
import StepCounter from '../StepCounter';
import UserPreferenceStepForm from './form';

const UserPreferenceStep = () => {
  const navigate = useNavigate();
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();
  const { wasMemberUpdated, isUpdatingMember, error, updateMember, reset } = useUpdateMember();
  const hasError = error && error.response.status === 400;

  const submitFormHandler = (values: any) => {
    // console.log('[] values', values);
    const requestPayload = {
      memberID: member?.id,
      payload: {
        ...getOnboardingStepsData(member),
        developer_segment: values?.developer_segment?.map((item: any) => item.value),
        profession_segment: values?.profession_segment,
      },
    };

    // console.log(requestPayload);
    updateMember(requestPayload);
  };

  // move to next step if member was updated
  useEffect(() => {
    if (wasMemberUpdated && hasCompletedOnboarding(member)) {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [wasMemberUpdated, increaseStep, navigate, member]);

  // if it missing other info show toast
  useEffect(() => {
    if (wasMemberUpdated && !hasCompletedOnboarding(member)) {
      reset();
      toast.error(t`Please complete all required fields to continue.`);
    }
  }, [wasMemberUpdated, increaseStep, navigate, member, reset]);

  return (
    <>
      <StepCounter />

      <div className="flex flex-col justify-center ">
        <UserPreferenceStepForm
          onSubmit={submitFormHandler}
          isSubmitting={isUpdatingMember}
          error={error}
          initialValues={{
            developer_segment: member?.developer_segment?.map((item: any) => ({
              label: item,
              value: item,
            })),
            profession_segment: member.profession_segment,
          }}
          hasError={hasError}
        />
      </div>
    </>
  );
};

export default UserPreferenceStep;
