import { useState } from 'react';

const useProjectActive = (projectID: string) => {
  const projectKEY = 'isActiveProject-' + projectID;
  const getIsProjectActive = localStorage.getItem(projectKEY);
  const [isActive, setIsActive] = useState(getIsProjectActive === 'true');

  return {
    isActive,
    setIsActive,
  };
};

export default useProjectActive;
