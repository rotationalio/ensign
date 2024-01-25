/* eslint-disable unused-imports/no-unused-vars */
import { t, Trans } from '@lingui/macro';
import { Button, Modal } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { EXTERNAL_LINKS } from '@/application';
import DownloadIcon from '@/components/icons/download-icon';
import Copy from '@/components/ui/Copy';
import { MIME_TYPES } from '@/constants/mimeTypes';
import downloadFile from '@/utils/download-file';
import { formatDate } from '@/utils/formatDate';

export type ApiKeyModalProps = {
  open: boolean;
  onClose: () => void;
  projectName: string;
  keyData: any;
};

const handleDownload = (data: any, filename: string) => {
  downloadFile(data, filename, MIME_TYPES.json);
};
export default function ApiKeyModal({ open, onClose, projectName, keyData }: ApiKeyModalProps) {
  const clientInfo = JSON.stringify({
    ClientID: keyData?.client_id || '',
    ClientSecret: keyData?.client_secret || '',
  });

  const dateCreated = formatDate(new Date(keyData?.created));
  const apiKeyFileName = `APIKey-${projectName}-${keyData?.name}-${dateCreated}`;

  const onCloseHandler = () => {
    handleDownload(clientInfo, apiKeyFileName);
    onClose();
  };

  return (
    <>
      <Modal
        open={open}
        title={t`Your New API Key`}
        data-testid="keyCreated"
        onClose={onClose}
        containerClassName="download-key-modal"
      >
        <>
          <div className="flex flex-col space-y-5 px-8 pb-5 text-sm">
            <p className="mt-3">
              <Trans>
                <span className="font-bold">Your API key is ready!</span> Your API key is a unique
                code that provides access to your project.
              </Trans>
            </p>
            <p>
              <Trans>
                Your API Key contains two parts: a Client ID and Client Secret. The Client ID is the
                unique identifier for your API key. Your Client Secret is the password for your API
                key.
              </Trans>
            </p>
            <p className="font-semibold">
              <Trans>Your API Key:</Trans>
            </p>
            <div className="relative flex flex-col break-words rounded-md border bg-[#FBF8EC] p-3 text-xs">
              <div className="space-y-3">
                <div className="flex flex-col pr-5">
                  <p className="mr-1 font-semibold">
                    <Trans>Client ID:</Trans>
                  </p>
                  <p className="items-center">
                    <span className="font-mono" data-testid="clientId">
                      {keyData?.client_id}
                    </span>
                    <span className="ml-1 drop-shadow-md " data-testid="copyID">
                      <Copy text={keyData?.client_id} />
                    </span>
                  </p>
                </div>
                <div className="flex flex-col">
                  <span className="font-semibold">
                    <Trans>Client Secret:</Trans>
                  </span>
                  <p>
                    <span className="font-mono" data-testid="clientSecret">
                      {keyData?.client_secret}
                    </span>
                    <span className="ml-1 " data-testid="copySecret">
                      <Copy text={keyData?.client_secret} />
                    </span>
                  </p>
                </div>
              </div>
              <div className="absolute right-3 top-3 flex gap-2">
                <button
                  onClick={() => handleDownload(clientInfo, apiKeyFileName)}
                  data-testid="download"
                >
                  <DownloadIcon className="h-4 w-4" />
                </button>
              </div>
            </div>
            <p className="font-semibold">
              <Trans>What to do next:</Trans>
            </p>
            <ol className="ml-5 list-decimal">
              <li>
                <Trans>
                  Download and securely store the key. You'll need it to access your project via the
                  API. For security purposes, this is the only time you will see the key.
                </Trans>
              </li>
              <li>
                <Trans>Use your API key to connect your services or models to your topic.</Trans>
              </li>
            </ol>
            <div className="rounded-md bg-[#FFDDDD] p-3">
              <h2 className="mb-3 font-semibold">
                <Trans>CAUTION!</Trans>
              </h2>
              <p>
                <Trans>
                  Avoid embedding API keys directly in your code. Instead, store them in environment
                  variables or in files outside of your application's source tree. If you misplace
                  an API key or it becomes compromised, revoke it and generate a new one.
                </Trans>
              </p>
              <p className="mt-3">
                <Trans>
                  Watch our video on{' '}
                  <Link
                    to={EXTERNAL_LINKS.PROTECT_API_KEYS_VIDEO}
                    className="underline"
                    target="_blank"
                  >
                    protecting your API keys
                  </Link>
                  .
                </Trans>
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
