import { useEffect, useState } from 'react';

// import { ProjectQueryResponse } from '@/features/projects/types/projectQueryService';
const usePaginateTopicQuery = (data: any) => {
  const [count, setCount] = useState<number>(0);
  const [result, setResult] = useState<any>([]);
  const [isNextClickDisabled, setIsNextClickDisabled] = useState<boolean>(true);
  const [isPrevClickDisabled, setIsPrevClickDisabled] = useState<boolean>(true);
  //   console.log('[] count page', count);
  //   console.log('[] count result length', data?.results?.length);
  const handleNextClick = () => {
    setCount(count + 1);
  };

  const handlePrevClick = () => {
    setCount(count - 1);
  };

  // define the default result when the page loads

  useEffect(() => {
    if (data?.results.length > 0) {
      setResult(data?.results[0]);
    }
  }, [data]);

  // increment the counteer when result is not empty

  useEffect(() => {
    if (data?.results?.length > 0 && count === 0) {
      setCount(count + 1);
    }
  }, [data, count]);

  useEffect(() => {
    if (data?.results.length > 0 && count > 0) {
      setResult(data?.results[count - 1]);
    }
  }, [data, count]);

  useEffect(() => {
    if (count === 0 || count === 1) {
      setIsPrevClickDisabled(true);
    } else {
      setIsPrevClickDisabled(false);
    }
  }, [count, result]);

  useEffect(() => {
    if (count === 0 || (data?.results?.length > 0 && count === data?.results?.length)) {
      setIsNextClickDisabled(true);
    } else {
      setIsNextClickDisabled(false);
    }
  }, [data, count]);

  return {
    result,
    isNextClickDisabled,
    isPrevClickDisabled,
    handleNextClick,
    handlePrevClick,
    counter: count,
  };
};

export default usePaginateTopicQuery;
