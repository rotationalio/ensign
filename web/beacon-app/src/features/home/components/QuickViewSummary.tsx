import QuickViewCard from './QuickViewCard';

const summary = {
  activeProjects: 1,
  topics: 1,
  apiKeys: 0,
  dataStorage: `0.0`,
};

function QuickViewSummary() {
  return (
    <div className="grid grid-cols-2 gap-y-10 gap-x-20 lg:grid-cols-4">
      <QuickViewCard title="Active Projects" color="#ECF6FF">
        {summary.activeProjects}
      </QuickViewCard>
      <QuickViewCard title="Topics" color="#FFE9DD">
        {summary.topics}
      </QuickViewCard>
      <QuickViewCard title="API Keys" color="#ECFADC">
        {summary.apiKeys}
      </QuickViewCard>
      <QuickViewCard title="Data Storage" color="#FBF8EC">
        {summary.dataStorage} GB
        <span className="ml-4 text-xs font-normal italic">0,00%</span>
      </QuickViewCard>
    </div>
  );
}

export default QuickViewSummary;
