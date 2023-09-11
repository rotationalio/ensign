import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';

// import { useOrgStore } from '@/store';
import { getOnboardingStepsData } from '../../../shared/utils';
import StepCounter from '../StepCounter';
import UserPreferenceStepForm from './form';

const UserPreferenceStep = () => {
  const navigate = useNavigate();
  // const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();
  const { wasMemberUpdated, isUpdatingMember, error, updateMember, hasMemberFailed } =
    useUpdateMember();
  const hasError = error && error.response.status === 400;

  const submitFormHandler = (values: any) => {
    const requestPayload = {
      memberID: member?.id,
      payload: {
        ...getOnboardingStepsData(member),
        developer_segment: values?.developer_segment?.map((item: any) => item.value),
        profession_segment: values?.profession_segment,
      },
    };
    updateMember(requestPayload);
  };

  useEffect(() => {
    if (wasMemberUpdated || !hasMemberFailed) {
      navigate(PATH_DASHBOARD.HOME);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [wasMemberUpdated, member]);

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
