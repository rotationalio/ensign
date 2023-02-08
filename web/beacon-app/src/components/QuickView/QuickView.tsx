import QuickViewCard from './QuickViewCard';

interface QuickViewData {
  name: string;
  value: number;
}

const QUICKVIEW_LIMIT = 4; // we should have 4 items in the quick view

export interface QuickViewProps {
  data: QuickViewData[];
}
const BRAND_COLORS = [
  'var(--color-primary-600)',
  'var(--color-secondary-200)',
  'var(--color-green-light)',
  'var(--color-secondary-100)',
];

const QuickView: React.FC<QuickViewProps> = ({ data }) => {
  return (
    <div className="grid grid-cols-2 gap-y-10 gap-x-20 lg:grid-cols-4">
      {data?.length === QUICKVIEW_LIMIT &&
        data.map((item, index) => (
          <QuickViewCard key={item.name} title={item.name} color={BRAND_COLORS[index]}>
            {item.value}
          </QuickViewCard>
        ))}
    </div>
  );
};

export default QuickView;
