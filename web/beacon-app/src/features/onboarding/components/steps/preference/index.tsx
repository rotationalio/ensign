import { useEffect } from 'react';

import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import UserPreferenceStepForm from './form';

const UserPreferenceStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();
  const { wasMemberUpdated, isUpdatingMember, reset, error, updateMember } = useUpdateMember();
  const hasError = error && error.response.status === 400;

  const submitFormHandler = (values: any) => {
    console.log('[] values', values);
    const { organization, name, workspace } = member;
    const requestPayload = {
      memberID: member?.id,
      payload: {
        organization,
        name,
        workspace,
        developer_segment: values.developer_segment.map((item: any) => item.value).join(','),
        profession_segment: values.profession_segment,
      },
    };

    console.log(requestPayload);
    updateMember(requestPayload);
  };

  // move to next step if member was updated
  useEffect(() => {
    if (wasMemberUpdated) {
      increaseStep();
      reset();
    }
  }, [wasMemberUpdated, increaseStep, reset]);

  return (
    <>
      <StepCounter />

      <div className="flex flex-col justify-center ">
        <UserPreferenceStepForm
          onSubmit={submitFormHandler}
          isSubmitting={isUpdatingMember}
          error={error}
          initialValues={{
            developer_segment: member.developer_segment?.split(',').map((item: any) => ({
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
