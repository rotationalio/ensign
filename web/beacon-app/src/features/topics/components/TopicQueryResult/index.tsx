import usePaginateTopicQuery from '../../hooks/usePaginateTopicQuery';
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

  const {
    result,
    isNextClickDisabled,
    isPrevClickDisabled,
    handleNextClick,
    handlePrevClick,
    counter,
  } = usePaginateTopicQuery(data);

  return (
    <div className="">
      <ViewingEvent
        totalResults={totalResults}
        totalEvents={data?.total_events}
        counter={counter}
      />
      <ResultHeader
        mimeType={result?.mimetype}
        eventType={result?.version}
        isBase64Encoded={result?.is_base64_encoded}
      />
      <QueryResultContent
        result={result?.data}
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
