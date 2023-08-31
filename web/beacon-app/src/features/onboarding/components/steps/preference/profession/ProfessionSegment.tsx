import CheckCircleIcon from '@/components/icons/check-circle';

import Header from './Header';
const ProfessionSegment = () => {
  return (
    <div>
      <Header />
      <div className="my-5">
        <ul className="grid w-full gap-20 md:grid-cols-3">
          <li>
            <input
              id="profession_segment"
              type="radio"
              value={'work'}
              name="profession_segment"
              className="peer hidden"
              required
            />
            <label
              htmlFor="profession_segment"
              className="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white p-5 text-gray-500 hover:bg-gray-100 hover:text-gray-600 peer-checked:border-blue-600 peer-checked:text-blue-600"
            >
              <div className="block">
                <div className="w-full text-lg font-bold">Work</div>
              </div>
              <CheckCircleIcon />
            </label>
          </li>
          <li>
            <input
              id="profession_segment"
              type="radio"
              value={'education'}
              name="profession_segment"
              className="peer hidden"
            />
            <label
              htmlFor="profession_segment"
              className="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white p-5 text-gray-500 hover:bg-gray-100 hover:text-gray-600 peer-checked:border-blue-600 peer-checked:text-blue-600"
            >
              <div className="block">
                <div className="w-full text-lg font-bold">Education</div>
              </div>

              <CheckCircleIcon />
            </label>
          </li>
          <li>
            <input
              id="profession_segment"
              type="radio"
              name="profession_segment"
              className="peer hidden"
              value={'personal'}
            />
            <label
              htmlFor="profession_segment"
              className="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white p-5 text-gray-500 hover:bg-gray-100 hover:text-gray-600 peer-checked:border-blue-600 peer-checked:text-blue-600"
            >
              <div className="block">
                <div className="w-full text-lg font-bold">Personnal</div>
              </div>
              <CheckCircleIcon />
            </label>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default ProfessionSegment;
