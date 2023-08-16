import { useEffect, useState } from 'react';

// import { ProjectQueryResponse } from '@/features/projects/types/projectQueryService';
const usePaginateTopicQuery = (data: any, onReset: boolean) => {
  const [count, setCount] = useState<number>(0);
  const [result, setResult] = useState<any>([]);
  const [isNextClickDisabled, setIsNextClickDisabled] = useState<boolean>(true);
  const [isPrevClickDisabled, setIsPrevClickDisabled] = useState<boolean>(true);

  const handleNextClick = () => {
    setCount(count + 1);
  };

  const handlePrevClick = () => {
    setCount(count - 1);
  };

  // define the default result when the data is not empty
  useEffect(() => {
    if (data?.results?.length > 0 && count === 0) {
      setResult(data?.results[0]);
    }
  }, [data, count]);

  // set default counter when result is available

  useEffect(() => {
    if (data?.results?.length > 0 && count === 0) {
      setCount(count + 1);
    }
  }, [data, count]);

  useEffect(() => {
    if (data?.results?.length > 0 && count > 0) {
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

  // clear the result when the data is empty
  useEffect(() => {
    if (onReset) {
      setResult([]);
      setCount(0);
    }
  }, [onReset]);

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
