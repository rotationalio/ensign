import React from 'react';

function CheckCircleIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      className="w-3.5 h-3.5 lg:h-4 lg:w-4"
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 16 12"
      {...props}
    >
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M1 5.917 5.724 10.5 15 1.5"
      />
    </svg>
  );
}

export default CheckCircleIcon;
