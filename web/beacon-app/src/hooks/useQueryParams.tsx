import queryString from 'query-string';
import { useLocation } from 'react-router-dom';

const useQueryParams = () => {
  const location = useLocation();
  const params = queryString.parse(location.search); // Parse the query string

  return params;
};

export default useQueryParams;
