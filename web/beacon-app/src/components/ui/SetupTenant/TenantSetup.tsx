import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

export default function TenantSetup() {
  return (
    <section>
      <TenantHeader />
      <section className="grid grid-cols-2">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </section>
  );
}
