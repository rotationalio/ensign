import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

import AccessDocumentationStep from './AccessDocumentationStep';
import GenerateApiKeyStep from './GenerateApiKeyStep';
import ProjectDetailsStep from './ProjectDetailsStep';

export default function QuickStart() {
  const { tenants } = useFetchTenants();

  return (
    <div className="space-y-10">
      <ProjectDetailsStep tenantID={tenants?.tenants[0]?.id} />
      <GenerateApiKeyStep />
      <AccessDocumentationStep />
    </div>
  );
}
