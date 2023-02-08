import QuickViewCard from './QuickViewCard';

interface QuickViewData {
  name: string;
  value: number;
}

const QUICKVIEW_CARD_LENGTH = 4; // we should have 4 items in the quick view

export interface QuickViewProps {
  data: QuickViewData[];
}
const BRAND_COLORS = ['#ECF6FF', '#FFE9DD', '#ECFADC', '#FBF8EC'];

const QuickView: React.FC<QuickViewProps> = ({ data }) => {
  return (
    <div className="grid grid-cols-2 gap-y-10 gap-x-20 lg:grid-cols-4">
      {data?.length === QUICKVIEW_CARD_LENGTH &&
        data.map((item, index) => (
          <QuickViewCard key={item.name} title={item.name} color={BRAND_COLORS[index]}>
            {item.value}
          </QuickViewCard>
        ))}
    </div>
  );
};

export default QuickView;
