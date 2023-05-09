/* eslint-disable unused-imports/no-unused-vars */
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

  return (
    <>
      <Modal
        open={open}
        title="Your API Key"
        containerClassName="overflow-scroll max-h-[90vh] max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="keyCreated"
        onClose={onClose}
      >
        <>
          <div className="gap-3 space-y-5 px-8 pb-5 text-sm">
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
            <div className="relative flex flex-col gap-2 rounded-md border bg-[#FBF8EC] p-3 text-xs">
              <div className="space-y-3">
                <p className="flex flex-col pr-5">
                  <span className="mr-1 font-semibold">Client ID:</span>
                  <span className="flex items-center">
                    <span className="font-mono" data-testid="clientId">
                      {data?.client_id}
                    </span>
                    <span className="ml-1 drop-shadow-md" data-testid="copyID">
                      <Copy text={data?.client_id} />
                    </span>
                  </span>
                </p>
                <p className="flex flex-col pr-5">
                  <span>
                    <span className="font-semibold">Client Secret: </span>
                    <span className="font-mono" data-testid="clientSecret">
                      {data?.client_secret}
                    </span>
                    <span className="ml-1" data-testid="copySecret">
                      <Copy text={data?.client_secret} />
                    </span>
                  </span>
                </p>
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
                onClick={onClose}
                data-testid="closeKey"
              >
                I read the above and <br />
                definitely saved this key
              </Button>
            </div>
          </div>
        </>
      </Modal>
    </>
  );
}
