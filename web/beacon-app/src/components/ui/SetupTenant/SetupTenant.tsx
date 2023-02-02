import { LandingFooter } from '@/components/auth/LandingFooter';
import { LandingHeader } from '@/components/auth/LandingHeader';
import { memo } from 'react';

import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

import { t } from '@lingui/macro';

function SetupTenant() {
  return (
    <div className="bg-hexagon bg-contain">
      <LandingHeader />
      <TenantHeader />
      <section className="mx-auto grid max-w-4xl grid-cols-2 pb-6">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
      <LandingFooter />
    </div>
  );
}

export default memo(SetupTenant)