import { projectQueryMockData } from '../../__mocks__';
import PaginatedViewButtons from './PaginatedViewButtons';
import QueryResultContent from './QueryResultContent';
import ResultHeader from './ResultHeader';
import ViewingEvent from './ViewingEvent';
interface TopicQueryResultProps {
  result: any;
  isFetching?: boolean;
}

const TopicQueryResult = ({ result, isFetching = false }: TopicQueryResultProps) => {
  const { data, error } = result;
  if (isFetching) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="">
      <ViewingEvent data={data || projectQueryMockData} />
      <ResultHeader
        mimeType={data?.mimeType}
        eventType={data?.eventType}
        isBase64Encoded={data?.isBase64Encoded}
      />
      <QueryResultContent result={data?.results} mimeType={data?.mimeType} />
      <PaginatedViewButtons />
    </div>
  );
};

export default TopicQueryResult;
