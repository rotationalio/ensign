import { Loader } from '@rotational/beacon-core';
import { memo } from 'react';

import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

function SetupTenant() {
  /* if(loading) {
    return <Loader />
  } */
  return (
    <section className="bg-hexagon bg-contain">
      <TenantHeader />
      <section className="mx-auto grid max-w-4xl grid-cols-2 pb-6">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </section>
  );
}

export default memo(SetupTenant)