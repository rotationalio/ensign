import { memo } from 'react';

import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

function SetupTenant() {
  return (
    <div>
      <TenantHeader />
      <section className="mx-auto max-w-4xl gap-12 md:flex">
        <div className="pb-6 sm:pb-0">
          <TenantQuickStart />
        </div>
        <CustomizeTenant />
      </section>
    </div>
  );
}

export default memo(SetupTenant);
