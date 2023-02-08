import './Loader.style.css';

import { twMerge } from 'tailwind-merge';

export type LoaderProps = {
  label?: string;
  labelProps?: React.DetailedHTMLProps<
    React.HTMLAttributes<HTMLParagraphElement>,
    HTMLParagraphElement
  >;
};

function Loader(props: LoaderProps) {
  const { labelProps, label, ...rest } = props;
  return (
    <div {...rest} className="text-center">
      <div className="relative inline-block">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="68.355"
          height="67.597"
          viewBox="0 0 68.355 67.597"
          id="loading-spinner"
          data-testid="loading-spinner"
        >
          <path
            data-name="Tracé 6"
            d="M456.691,237.38A34.214,34.214,0,0,0,397,259.8h4.682c0-14.008,11.163-32.207,30.588-29.742,7.648-1.352,16.226,3.469,20.912,8.923l-10.618,12.37h22.789V228.812Zm-25.513,53.4c-7.81,0-17.631-6.877-22.317-12.329l10.929-10.207H397v22.536l9.364-10.437c6.248,6.971,14.654,13.249,24.814,13.249a33.987,33.987,0,0,0,34.178-33.8h-3.121C462.234,273.806,445.343,290.785,431.178,290.785Z"
            transform="translate(-397 -226)"
            fill="#1d65a6"
          />
        </svg>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="28.319"
          height="22.841"
          viewBox="0 0 28.319 22.841"
          style={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
          }}
        >
          <path
            id="Tracé_7"
            data-name="Tracé 7"
            d="M471.109,309.816a5.05,5.05,0,0,0-.578-2.356,8.348,8.348,0,0,0-1.87-2.8,9.025,9.025,0,0,0-2.879-1.913,7.721,7.721,0,0,0-2.446-4.3,9.006,9.006,0,0,0-6.385-2.452,9.155,9.155,0,0,0-5.471,1.719,8.221,8.221,0,0,0-3.218,4.471,6.826,6.826,0,0,0-3.965,2.344,6.354,6.354,0,0,0-1.506,4.183,6.056,6.056,0,0,0,2.072,4.627,7.107,7.107,0,0,0,5.007,1.887h4.84a9.2,9.2,0,0,0,7.389,3.606c4.981,0,9.01-3.762,9.01-8.414Zm-21.239,3.005a4.57,4.57,0,0,1-3.18-1.2,4.006,4.006,0,0,1-1.326-3.005,3.876,3.876,0,0,1,1.326-2.97,4.443,4.443,0,0,1,3.18-1.238h.644a5.693,5.693,0,0,1,1.879-4.255,6.793,6.793,0,0,1,9.113,0,6.036,6.036,0,0,1,1.351,1.887,5.282,5.282,0,0,0-.76-.036c-4.981,0-9.01,3.762-9.01,8.414a8.2,8.2,0,0,0,.373,2.4Zm12.229,3.606a6.024,6.024,0,1,1,6.436-6.01A6.236,6.236,0,0,1,462.1,316.427Zm.644-5.709,3.681,2.032-.966,1.467-4.647-2.6v-6.01h1.931Z"
            transform="translate(-442.79 -295.991)"
            fill="#1d65a6"
          />
        </svg>
      </div>
      {label && (
        <p {...labelProps} className={twMerge('mt-2 text-xs text-primary', labelProps?.className)}>
          {label}
        </p>
      )}
    </div>
  );
}

export default Loader;
