import AppLayout from '@/components/layout/AppLayout';

import Step from '../components/Step';
import Stepper from '../components/Stepper';

const OnboardingPage = () => {
  return (
    <AppLayout>
      <Stepper />
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
