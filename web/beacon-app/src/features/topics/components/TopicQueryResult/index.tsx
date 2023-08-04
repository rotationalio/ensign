import PaginatedViewButtons from './PaginatedViewButtons';
import QueryResultContent from './QueryResultContent';
import ResultHeader from './ResultHeader';
import ViewingEvent from './ViewingEvent';
interface TopicQueryResultProps {
  data: any;
  isFetching?: boolean;
  error?: any;
}

const TopicQueryResult = ({ data }: TopicQueryResultProps) => {
  const totalResults = data?.results?.length;
  const results = data?.results[6];
  const isBase64Encoded = results?.is_base64_encoded;

  console.log('##rendering TopicQueryResult');

  return (
    <div className="">
      <ViewingEvent totalResults={totalResults} totalEvents={data?.total_events} />
      <ResultHeader
        mimeType={results?.mimetype}
        eventType={results?.version}
        isBase64Encoded={isBase64Encoded}
      />
      <QueryResultContent
        result={results?.data}
        mimeType={results?.mimetype}
        isBase64Encoded={isBase64Encoded}
      />
      <PaginatedViewButtons />
    </div>
  );
};

export default TopicQueryResult;
