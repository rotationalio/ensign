import OtterLookingDown from '@/components/icons/otter-looking-down';

import SetupNewProject from '../components/SetupNewProject';
import AccessDocumentationStep from './AccessDocumentationStep';

export default function QuickStart() {
  return (
    <div className="relative w-full space-y-10">
      <OtterLookingDown className="2xl:right-18 hidden md:absolute md:-top-[215px] md:left-8 md:block md:scale-[.65] xl:left-3/4" />
      <SetupNewProject />
      <AccessDocumentationStep />
    </div>
  );
}
