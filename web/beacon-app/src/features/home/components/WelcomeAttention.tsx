import { Trans } from '@lingui/macro';

import RotationalNotifcationImage from '@/assets/images/rotational-ipn.png'; // todo: ensure to use a better image , this one seems creepy
import { Image } from '@/components/ui/Image';

const WelcomeAttention = () => {
  return (
    <>
      <div
        className="px-auto mb-8 mt-4 flex flex-row items-center justify-between space-x-10 rounded-md border border-black/30 p-2 px-5 text-justify"
        data-cy="projWelcome"
      >
        <div className="flex flex-col space-y-10 ">
          <p className="text-md">
            <Trans>
              <span className="font-bold"> Welcome to Ensign </span>, your all-in-one platform for
              real-time data management. Ensign is a flexible database meets streaming engine for
              data teams to build and deploy real-time models, data products, and services.
            </Trans>
          </p>
          <p>
            <Trans>Ready to dive in? Learn how to use Ensign or start your first project.</Trans>
          </p>
        </div>

        <Image src={RotationalNotifcationImage} alt="rebecca preview image" className="w-1/4" />
      </div>
    </>
  );
};

export default WelcomeAttention;
