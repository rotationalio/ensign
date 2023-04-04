import React from 'react';

function Union(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="14"
      height="14"
      viewBox="0 0 14 14"
      fill="none"
      {...props}
    >
      <path
        d="M7.625 1.375C7.625 1.02982 7.34518 0.75 7 0.75C6.65482 0.75 6.375 1.02982 6.375 1.375V6.375H1.375C1.02982 6.375 0.75 6.65482 0.75 7C0.75 7.34518 1.02982 7.625 1.375 7.625H6.375V12.625C6.375 12.9702 6.65482 13.25 7 13.25C7.34518 13.25 7.625 12.9702 7.625 12.625V7.625H12.625C12.9702 7.625 13.25 7.34518 13.25 7C13.25 6.65482 12.9702 6.375 12.625 6.375H7.625V1.375Z"
        fill="inherit"
      />
    </svg>
  );
}

export default Union;
