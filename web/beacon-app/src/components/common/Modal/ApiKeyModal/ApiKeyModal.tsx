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
  downloadFile(data, filename, MIME_TYPES.txt);
};
export default function ApiKeyModal({ open, onClose, data }: ApiKeyModalProps) {
  return (
    <>
      <Modal
        open={open}
        title="Your API Key"
        containerClassName="h-[90vh] overflow-scroll max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
      >
        <>
          <button onClick={onClose} className="bg-transparent absolute top-4 right-4 border-none">
            <CloseIcon className="h-4 w-4" />
          </button>
          <div className="gap-3 space-y-5 px-8 text-sm">
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
            <div className="relative flex flex-col gap-2 rounded-xl border bg-[#FBF8EC] p-3 text-xs">
              <p className="flex">
                <span className="mr-1 font-semibold">Client ID:</span> {data?.client_id}
                <span className="ml-2 flex space-x-1">
                  <Copy text={data?.client_id} />
                  <button onClick={() => handleDownload(data?.client_id, 'client_id')}>
                    <DownloadIcon className="h-4 w-4" />
                  </button>
                </span>
              </p>
              <p className="flex items-center">
                <span>
                  <span className="font-semibold">Client Secret:</span> {data?.client_secret}
                </span>
                <span className="ml-2 flex space-x-1">
                  <Copy text={data?.client_secret} />
                  <button onClick={() => handleDownload(data?.client_secret, 'client_secret')}>
                    <DownloadIcon className="h-4 w-4" />
                  </button>
                </span>
              </p>
              <div className="absolute top-3 right-3 flex gap-2"></div>
            </div>
            <div className="rounded-xl bg-[#FFDDDD] p-3">
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
              <Button size="medium" className="w-full max-w-[350px] p-2" onClick={onClose}>
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
