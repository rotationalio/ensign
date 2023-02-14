import { AriaButton as Button, Modal } from '@rotational/beacon-core';

import CopyIcon from '@/components/icons/copy-icon';
import DownloadIcon from '@/components/icons/download-icon';

export type ApiKeyModalProps = {
  open: boolean;
};

export default function ApiKeyModal({ open }: ApiKeyModalProps) {
  return (
    <>
      <Modal open={open} title="Your API Key" size="large">
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
          <div className="relative flex flex-col gap-2 rounded-xl border bg-[#FBF8EC] p-3">
            <p>
              <span className="font-semibold">Client ID:</span> DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa
            </p>
            <p>
              <span className="font-semibold">Client Secret</span>{' '}
              wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS
            </p>
            <div className="absolute top-3 right-3 flex gap-2">
              <CopyIcon className="h-5 w-5" />
              <DownloadIcon className="h-5 w-5" />
            </div>
          </div>
          <div className="rounded-xl bg-[#FFDDDD] p-3">
            <h2 className="mb-3 font-semibold">CAUTION!</h2>
            <p>
              We don’t recommend that you embed keys directly in your code (they’re private after
              all!). Instead of embedding your API keys in your applications, store them in
              environment variables or in files outside of your application&apos;s source tree. If
              you misplace this API key or it becomes compromised, revoke it and generate a new one.
            </p>
          </div>
          <div className="text-center">
            <Button size="medium" className="w-full max-w-[350px]">
              I read the above and <br />
              definitely saved this key
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
}
