interface IconProps {
  className?: string;
}
export function SortIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 320 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M41 288h238c21.4 0 32.1 25.9 17 41L177 448c-9.4 9.4-24.6 9.4-33.9 0L24 329c-15.1-15.1-4.4-41 17-41zm255-105L177 64c-9.4-9.4-24.6-9.4-33.9 0L24 183c-15.1 15.1-4.4 41 17 41h238c21.4 0 32.1-25.9 17-41z"></path>
    </svg>
  );
}

export function SortUpIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 320 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M279 224H41c-21.4 0-32.1-25.9-17-41L143 64c9.4-9.4 24.6-9.4 33.9 0l119 119c15.2 15.1 4.5 41-16.9 41z"></path>
    </svg>
  );
}

export function SortDownIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 320 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M41 288h238c21.4 0 32.1 25.9 17 41L177 448c-9.4 9.4-24.6 9.4-33.9 0L24 329c-15.1-15.1-4.4-41 17-41z"></path>
    </svg>
  );
}

export function ThreeDotIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 320 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M256 256c0 17.7-14.3 32-32 32s-32-14.3-32-32 14.3-32 32-32 32 14.3 32 32zm-32-192c-17.7 0-32 14.3-32 32s14.3 32 32 32 32-14.3 32-32-14.3-32-32-32zm0 384c-17.7 0-32 14.3-32 32s14.3 32 32 32 32-14.3 32-32-14.3-32-32-32z"></path>
    </svg>
  );
}

export function HThreeDotIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      height={16}
      width={16}
      viewBox="0 0 24 24"
    >
      <circle cx="5" cy="12" r="2"></circle>
      <circle cx="12" cy="12" r="2"></circle>
      <circle cx="19" cy="12" r="2"></circle>
    </svg>
  );
}

export function NoDataIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 512 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M256 8C119.033 8 8 119.033 8 256s111.033 248 248 248 248-111.033 248-248S392.967 8 256 8zm0 464c-110.28 0-200-89.72-200-200S145.72 56 256 56s200 89.72 200 200-89.72 200-200 200z"></path>
      <path d="M256 128c-44.183 0-80 35.817-80 80s35.817 80 80 80 80-35.817 80-80-35.817-80-80-80zm0 128c-26.51 0-48-21.49-48-48s21.49-48 48-48 48 21.49 48 48-21.49 48-48 48z"></path>

      <path d="M256 320c-8.837 0-16-7.163-16-16v-96c0-8.837 7.163-16 16-16s16 7.163 16 16v96c0 8.837-7.163 16-16 16z"></path>
    </svg>
  );
}

// checkmark svg icon

export function SuccessIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 512 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M256 8C119.033 8 8 119.033 8 256s111.033 248 248 248 248-111.033 248-248S392.967 8 256 8zm0 464c-110.28 0-200-89.72-200-200S145.72 56 256 56s200 89.72 200 200-89.72 200-200 200z"></path>

      <path d="M352 176.5L207.5 336 160 288.5l48-48 47.5 47.5L352 176.5z"></path>
    </svg>
  );
}

// error svg icon

export function ErrorIcon({ className }: IconProps) {
  return (
    <svg
      className={className}
      stroke="currentColor"
      fill="currentColor"
      strokeWidth="0"
      viewBox="0 0 512 512"
      height="1em"
      width="1em"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M256 8C119.033 8 8 119.033 8 256s111.033 248 248 248 248-111.033 248-248S392.967 8 256 8zm0 464c-110.28 0-200-89.72-200-200S145.72 56 256 56s200 89.72 200 200-89.72 200-200 200z"></path>
      <path d="M352 176.5L207.5 336 160 288.5l48-48 47.5 47.5L352 176.5z"></path>
    </svg>
  );
}

// info svg icon

export function InfoIcon({ className }: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      role="img"
      className={className}
    >
      <title>Info</title>
      <desc>An icon representing the letter &lsquo;i&lsquo; in a circle</desc>
      <circle cx="12" cy="12" r="10"></circle>
      <line x1="12" y1="16" x2="12" y2="12"></line>
      <line x1="12" y1="8" x2="12.01" y2="8"></line>
    </svg>
  );
}

// warning svg icon

export function WarningIcon({ className }: IconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      role="img"
      className={className}
    >
      <title>Alert</title>
      <desc>An icon representing an exclamation mark in an octogone</desc>
      <polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"></polygon>
      <line x1="12" y1="8" x2="12" y2="12"></line>
      <line x1="12" y1="16" x2="12.01" y2="16"></line>
    </svg>
  );
}

export const XIcon = ({ className }: IconProps) => (
  <svg
    width="22"
    height="22"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    role="img"
    className={className}
  >
    <title>X</title>
    <desc>An icon representing an X</desc>
    <line x1="18" y1="6" x2="6" y2="18"></line>
    <line x1="6" y1="6" x2="18" y2="18"></line>
  </svg>
);

export const ChevronDownIcon = ({ className }: IconProps) => (
  <svg
    className={className}
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 320 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="currentColor"
      d="M10.5 192.5L160 342l149.5-149.5c4.7-4.7 12.3-4.7 17 0l35.3 35.3c4.7 4.7 4.7 12.3 0 17L176 434.3c-4.7 4.7-12.3 4.7-17 0L10.5 209.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
  </svg>
);

export const ChevronRightIcon = ({ className }: IconProps) => (
  <svg
    className={className}
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 320 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="currentColor"
      d="M10.5 192.5L160 342l149.5-149.5c4.7-4.7 12.3-4.7 17 0l35.3 35.3c4.7 4.7 4.7 12.3 0 17L176 434.3c-4.7 4.7-12.3 4.7-17 0L10.5 209.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
  </svg>
);

export const ChevronLeftIcon = ({ className }: IconProps) => (
  <svg
    className={className}
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 320 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="currentColor"
      d="M309.5 209.5L160 59l-149.5 149.5c-4.7 4.7-12.3 4.7-17 0L58.2 174.2c-4.7-4.7-4.7-12.3 0-17L176 77.7c4.7-4.7 12.3-4.7 17 0L309.5 192.5c4.7 4.7 4.7 12.3 0 17z"
    ></path>
  </svg>
);

export const ChevronDoubleRightIcon = ({ className }: IconProps) => (
  <svg
    className={className}
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 448 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="currentColor"
      d="M10.5 192.5L160 342l149.5-149.5c4.7-4.7 12.3-4.7 17 0l35.3 35.3c4.7 4.7 4.7 12.3 0 17L176 434.3c-4.7 4.7-12.3 4.7-17 0L10.5 209.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
    <path
      fill="currentColor"
      d="M224.5 192.5L374 342l149.5-149.5c4.7-4.7 12.3-4.7 17 0l35.3 35.3c4.7 4.7 4.7 12.3 0 17L390 434.3c-4.7 4.7-12.3 4.7-17 0L224.5 209.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
  </svg>
);

export const ChevronDoubleLeftIcon = ({ className }: IconProps) => (
  <svg
    className={className}
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 448 512"
    height="1em"
    width="1em"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="currentColor"
      d="M224.5 209.5L374 59l149.5 149.5c4.7 4.7 4.7 12.3 0 17L390 337.8c-4.7 4.7-12.3 4.7-17 0L224.5 222.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
    <path
      fill="currentColor"
      d="M10.5 209.5L160 59l149.5 149.5c4.7 4.7 4.7 12.3 0 17L176 337.8c-4.7 4.7-12.3 4.7-17 0L10.5 222.5c-4.7-4.7-4.7-12.3 0-17z"
    ></path>
  </svg>
);

export const StatusColorIcon = (props: React.SVGProps<SVGSVGElement>) => {
  return (
    <svg
      width="6"
      height="6"
      viewBox="0 0 6 6"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M6 3C6 4.65685 4.65685 6 3 6C1.34315 6 0 4.65685 0 3C0 1.34315 1.34315 0 3 0C4.65685 0 6 1.34315 6 3Z"
        fill={props.fill}
      />
    </svg>
  );
};
