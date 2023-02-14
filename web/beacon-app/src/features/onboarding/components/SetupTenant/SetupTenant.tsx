import { memo } from 'react';

import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

function SetupTenant() {
  return (
    <div>
      <TenantHeader />
      <section className="mx-auto grid max-w-4xl grid-cols-2 pb-6">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </div>
  );
}

export default memo(SetupTenant);
