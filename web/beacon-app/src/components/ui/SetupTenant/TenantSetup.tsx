import CustomizeTenant from './CustomizeTenant';
import TenantHeader from './TenantHeader';
import TenantQuickStart from './TenantQuickStart';

export default function TenantSetup() {
  return (
    <section className="bg-hexagon bg-contain">
      <TenantHeader />
      <section className="mx-auto grid max-w-6xl grid-cols-2 pb-6">
        <TenantQuickStart />
        <CustomizeTenant />
      </section>
    </section>
  );
}
