import { memo } from 'react';

import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

function SetupTenant() {
  return (
    <div>
      <TenantHeader />
      <section className="mx-auto flex max-w-4xl space-x-9">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </div>
  );
}

export default memo(SetupTenant);
