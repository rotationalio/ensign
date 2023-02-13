export default function TenantHeader() {
  return (
    <main className="mx-auto mt-16 mb-12 max-w-7xl">
      <h1 className="text-center text-5xl font-bold">Tenant Setup</h1>
      <p className="mt-10 text-2xl">
        Tenants are collections of settings. You can think of a tenant as your environment. You can
        create one or more tenants. For example, you can create separate tenants for your
        development, staging, and production environments.
      </p>
      <p className="mt-10 text-2xl">
        You can start by selecting the <span className="font-bold">Quick Start</span> tenant on the{' '}
        <span className="font-bold">Starter Plan</span>. If you&#39;d like to customize your tenant
        based on cloud provider(s) and regions, you&#39;ll have to select a{' '}
        <span className="font-bold">paid plan</span>. You can change, add, or remove tenants later.
      </p>
    </main>
  );
}
