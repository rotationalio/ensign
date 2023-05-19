import OtterLookingDown from '@/components/icons/otter-looking-down';

import SetupNewProject from '../components/SetupNewProject';
import AccessDocumentationStep from './AccessDocumentationStep';

export default function QuickStart() {
  return (
    <div className="relative w-full space-y-10">
      <OtterLookingDown className="hidden md:absolute md:-top-[215px] md:-right-20 md:block md:scale-[.65] " />
      <SetupNewProject />
      <AccessDocumentationStep />
    </div>
  );
}
