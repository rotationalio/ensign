/* eslint-disable unused-imports/no-unused-vars */
import { Button, Modal } from '@rotational/beacon-core';

import { Close as CloseIcon } from '@/components/icons/close';
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
    client_secret: data?.client_secret || '',
    client_id: data?.client_id || '',
  });

  return (
    <>
      <Modal
        open={open}
        title="Your API Key"
        containerClassName="overflow-scroll max-h-[90vh] max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="keyCreated"
      >
        <>
          <button onClick={onClose} className="bg-transparent absolute top-4 right-4 border-none">
            <CloseIcon className="h-4 w-4" />
          </button>
          <div className="gap-3 space-y-5 px-8 pb-5 text-sm">
            <p className="my-3">
              <span className="font-bold text-primary-900">Sweet!</span> you&apos; got a brand new
              pair of <span className="line-through">roller skates</span> API keys!
            </p>
            <p className="text-danger-500">
              For security purposes, this is the only time you will see the key. Please copy and
              securely store the key.
            </p>
            <p>
              <span className="font-semibold">Your API Key:</span> your API key contains two parts:
              your ClientID and ClientSecret. You&apos;ll need both to sign to Ensign!
            </p>
            <div className="relative flex flex-col gap-2 rounded-md border bg-[#FBF8EC] p-3 text-xs">
              <div className="space-y-2">
                <p className="flex items-center pr-5">
                  <span className="mr-1 font-semibold">Client ID:</span>
                  <span className="flex items-center">
                    <span
                      className="flex items-center rounded-md bg-white p-1"
                      data-testid="clientId"
                    >
                      {data?.client_id}
                    </span>
                    <span className="ml-1 drop-shadow-md" data-testid="copyID">
                      <Copy text={data?.client_id} />
                    </span>
                  </span>
                </p>
                <p className="flex items-center pr-5">
                  <span>
                    <span className="font-semibold">Client Secret: </span>
                    <span
                      className="rounded-md bg-white p-1 leading-relaxed"
                      data-testid="clientSecret"
                    >
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
                environment variables or in files outside of your application&apos;s source tree. If
                you misplace this API key or it becomes compromised, revoke it and generate a new
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
