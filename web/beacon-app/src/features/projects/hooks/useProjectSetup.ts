import { t } from '@lingui/macro';

import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';
import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';

import { useFetchTenantProjects } from './useFetchTenantProjects';

export const useProjectSetup = (projectID: string) => {
  const { tenants } = useFetchTenants();
  // fetch project list
  const { projects } = useFetchTenantProjects(tenants?.tenants[0]?.id);
  // fetch api keys
  const { apiKeys } = useFetchApiKeys(projectID);
  // fetch topics list
  const { topics } = useFetchTopics(projectID);

  const hasProject = projects?.tenant_projects?.length > 0;
  const hasApiKeys = apiKeys?.api_keys?.length > 0;
  const hasTopics = topics?.topics?.length > 0;
  const hasAlreadySetup = hasProject && hasApiKeys && hasTopics;
  const hasTenant = tenants?.tenants?.length > 0;

  // generate warning message for user to setup project
  const message = () => {
    if (!hasProject) return t`You don't have any projects yet. Please create a project first.`;
    if (!hasApiKeys) return t`You don't have any API keys yet. Please create an API key.`;
    if (!hasTopics && hasProject) return t`You don't have any topics yet. Please create a topic.`;
    if (!hasApiKeys && !hasTopics) return t`Your project needs topics and API keys.`;
    if (!hasApiKeys && hasTopics) return t`Your project needs API keys.`;
  };

  return {
    hasAlreadySetup,
    hasProject,
    hasApiKeys,
    hasTopics,
    hasTenant,
    warningMessage: message(),
  };
};

export default useProjectSetup;
