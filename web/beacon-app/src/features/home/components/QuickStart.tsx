import SetupNewProject from '../components/SetupNewProject';
import AccessDocumentationStep from './AccessDocumentationStep';

export default function QuickStart() {
  return (
    <div className="space-y-10">
      <SetupNewProject />
      <AccessDocumentationStep />
    </div>
  );
}
