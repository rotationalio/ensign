import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

export default function TenantSetup() {
  return (
    <section>
      s
      <TenantHeader />
      <section className="mx-auto grid max-w-6xl grid-cols-2">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </section>
  );
}
