import AccessDocumentationStep from './AccessDocumentationStep';
import GenerateApiKeyStep from './GenerateApiKeyStep';
import ProjectDetailsStep from './ProjectDetailsStep';

export default function Steps() {
  return (
    <div className="space-y-10">
      <ProjectDetailsStep />
      <GenerateApiKeyStep />
      <AccessDocumentationStep />
    </div>
  );
}
