import HostedDataIcon from '@/components/icons/hosted-data-icon';
import AccessDashboard from '@/components/ui/AccessDashboard/AccessDashboard';

function SetupTenantComplete() {
  return (
    <section className="mx-auto grid max-w-4xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <HostedDataIcon className="mx-auto mt-6" />
      <div>
        <h2 className="mt-8 font-bold">
          Set up Your Tenant <span className="font-normal">(required)</span>
        </h2>
        <p className="mt-8">
          Your tenant is your team&apos;s control panel for all projects and topics. Specify
          preferences around encryption, privacy, and locality (e.g. for GDPR, CCPA, etc).
        </p>
      </div>
      <div className="mx-auto mt-36">
        <AccessDashboard />
      </div>
    </section>
  );
}

export default SetupTenantComplete;
