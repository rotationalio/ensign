import copy from 'copy-to-clipboard';
import { useEffect, useState } from 'react';

import Check from '@/components/icons/check';
import CopyIcon from '@/components/icons/copy-icon';

export const useCopy = (text: string): [boolean, () => void] => {
  const [isCopied, setIsCopied] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      setIsCopied(false);
    }, 5000);
  }, [isCopied]);

  const handleCopy = () => {
    setIsCopied(copy(text));
  };

  return [isCopied, handleCopy];
};

export default function Copy({ text }: { text: string }) {
  const [isCopied, setIsCopied] = useCopy(text);

  return isCopied ? (
    <button>
      <Check className="h-4 w-4" />
    </button>
  ) : (
    <button onClick={setIsCopied}>
      <CopyIcon className="h-4 w-4" />
    </button>
  );
}
