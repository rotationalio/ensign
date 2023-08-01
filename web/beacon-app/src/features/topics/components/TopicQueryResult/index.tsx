import ResultHeader from './ResultHeader';

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
    <div>
      <ResultHeader
        mimeType={data?.mimeType}
        eventType={data?.eventType}
        isBase64Encoded={data?.isBase64Encoded}
      />
    </div>
  );
};

export default TopicQueryResult;
