import { useEffect } from 'react';

import { useOrgStore } from '@/store';

const Indicator = () => {
  const onboarding = useOrgStore((state: any) => state.onboarding) as any;
  const { currentStep } = onboarding as any;

  useEffect(() => {
    if (currentStep && currentStep > 0) {
      const stepperItems = document.querySelectorAll('.stepper-item');
      stepperItems.forEach((item, index) => {
        if (index < currentStep - 1) {
          item.classList.remove('bg-gray-100', 'ring-white');
          item.classList.add('bg-green-500', 'ring-green-500');
        } else {
          item.classList.remove('bg-green-500', 'ring-green-500');
          item.classList.add('bg-gray-100', 'ring-white');
        }
      });
    }
  }, [currentStep]);
  return (
    <span className="stepper-item absolute -left-[4px] mt-1 flex h-2 w-2 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white"></span>
  );
};

export default Indicator;
