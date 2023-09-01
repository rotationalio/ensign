import StepCounter from '../StepCounter';
import UserPreferenceStepForm from './form';

const UserPreferenceStepStep = () => {
  const submitFormHandler = (values: any) => {
    console.log('[] values', values);
  };
  return (
    <>
      <StepCounter />

      <div className="flex flex-col justify-center ">
        <UserPreferenceStepForm onSubmit={submitFormHandler} />
      </div>
    </>
  );
};

export default UserPreferenceStepStep;
