/* eslint-disable unused-imports/no-unused-vars */
import { Trans } from '@lingui/macro';
import { Button, Modal } from '@rotational/beacon-core';

import DownloadIcon from '@/components/icons/download-icon';
import Copy from '@/components/ui/Copy';
import { MIME_TYPES } from '@/constants/mimeTypes';
import downloadFile from '@/utils/download-file';

export type ApiKeyModalProps = {
  open: boolean;
  onClose: () => void;
  data: any;
};

const handleDownload = (data: any, filename: string) => {
  downloadFile(data, filename, MIME_TYPES.json);
};
export default function ApiKeyModal({ open, onClose, data }: ApiKeyModalProps) {
  const clientInfo = JSON.stringify({
    ClientID: data?.client_id || '',
    ClientSecret: data?.client_secret || '',
  });

  const onCloseHandler = () => {
    // download the api key
    //then close the modal
    handleDownload(clientInfo, 'client');
    onClose();
  };

  return (
    <>
      <Modal
        open={open}
        title="Your API Key"
        data-testid="keyCreated"
        onClose={onClose}
        containerClassName="w-[35vw]"
      >
        <>
          <div className="flex flex-col space-y-5 px-8 pb-5 text-sm">
            <p className="my-3">
              <span className="font-bold text-primary-900">Sweet!</span> you&apos;ve got a brand new
              pair of <span className="line-through">roller skates</span> API keys!
            </p>
            <div className="text-danger-500">
              <p>For security purposes, this is the only time you will see the key.</p>
              <p>Please copy and securely store the key.</p>
            </div>
            <p>
              <span className="font-semibold">Your New API Key:</span> your API key contains two
              parts: your ClientID and ClientSecret. You&apos;ll need both to sign to Ensign!
            </p>
            <div className="relative flex flex-col rounded-md border bg-[#FBF8EC] p-3 text-xs">
              <div className="w-fit space-y-3">
                <div className="flex flex-col pr-5">
                  <p className="mr-1 font-semibold">Client ID:</p>
                  <p className="items-center">
                    <span className="font-mono" data-testid="clientId">
                      {data?.client_id}
                    </span>
                    <span className="ml-1 drop-shadow-md " data-testid="copyID">
                      <Copy text={data?.client_id} />
                    </span>
                  </p>
                </div>
                <div className="flex flex-col">
                  <span className="font-semibold">Client Secret: </span>
                  <span className="font-mono" data-testid="clientSecret">
                    {data?.client_secret}
                  </span>
                  <span className="ml-1 flex " data-testid="copySecret">
                    <Copy text={data?.client_secret} />
                  </span>
                </div>
              </div>
              <div className="absolute top-3 right-3 flex gap-2">
                <button onClick={() => handleDownload(clientInfo, 'client')} data-testid="download">
                  <DownloadIcon className="h-4 w-4" />
                </button>
              </div>
            </div>
            <div className="rounded-md bg-[#FFDDDD] p-3">
              <h2 className="mb-3 font-semibold">CAUTION!</h2>
              <p>
                We don’t recommend that you embed keys directly in your code (they’re private after
                all!). Instead of embedding your API keys in your applications, store them in
                environment variables or in files outside of your application&apos;s source tree.
              </p>
              <p className="mt-3">
                If you misplace this API key or it becomes compromised, revoke it and generate a new
                one.
              </p>
            </div>
            <div className="text-center">
              <Button
                size="medium"
                className="w-full max-w-[350px] p-2"
                onClick={onCloseHandler}
                data-testid="closeKey"
              >
                <Trans>Download API Keys</Trans>
              </Button>
            </div>
          </div>
        </>
      </Modal>
    </>
  );
}
