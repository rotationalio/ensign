import usePaginateTopicQuery from '../../hooks/usePaginateTopicQuery';
import PaginatedViewButtons from './PaginatedViewButtons';
import QueryResultContent from './QueryResultContent';
import ResultHeader from './ResultHeader';
import ViewingEvent from './ViewingEvent';
interface TopicQueryResultProps {
  data: any;
  isFetching?: boolean;
  error?: any;
  onReset?: boolean;
  hasInvalidQuery: boolean;
}

const TopicQueryResult = ({ data, onReset, hasInvalidQuery }: TopicQueryResultProps) => {
  const {
    result,
    isNextClickDisabled,
    isPrevClickDisabled,
    handleNextClick,
    handlePrevClick,
    counter,
  } = usePaginateTopicQuery(data, onReset || false);
  return (
    <div className="">
      <ViewingEvent
        totalResults={data?.results?.length}
        totalEvents={data?.total_events}
        counter={counter}
        metadataResult={result?.metadata}
      />
      <ResultHeader
        mimeType={result?.mimetype}
        eventType={result?.version}
        isBase64Encoded={result?.is_base64_encoded}
      />
      <QueryResultContent
        result={result?.data}
        hasInvalidQuery={hasInvalidQuery}
        mimeType={result?.mimetype}
        isBase64Encoded={result?.is_base64_encoded}
      />
      <PaginatedViewButtons
        onClickNext={handleNextClick}
        onClickPrevious={handlePrevClick}
        isNextDisabled={isNextClickDisabled}
        isPreviousDisabled={isPrevClickDisabled}
      />
    </div>
  );
};

export default TopicQueryResult;
