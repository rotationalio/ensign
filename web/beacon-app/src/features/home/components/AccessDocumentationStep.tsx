import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

export default function AccessDocumentationStep() {
  return (
    <>
      <CardListItem
        title={t`Access Resources`}
        titleClassName="text-lg"
        className="min-h-[130px]"
        contentClassName="my-2"
      >
        <div className="mt-2 flex flex-col justify-between gap-4 px-3 xl:flex-row">
          <p className="text-md w-full sm:w-4/5">
            <Trans>
              Love seeing examples with real code? Prefer watching tutorial videos? Still learning
              the basics? We've got you covered!
            </Trans>
          </p>
          <div className="item-center place-items-center">
            <a
              href="https://ensign.rotational.dev/getting-started/"
              target="_blank"
              rel="noopener noreferrer"
              data-testid="viewDocsLink"
            >
              <Button size="medium" data-testid="viewDocs">
                <Trans>Access</Trans>
              </Button>
            </a>
          </div>
        </div>
      </CardListItem>
    </>
  );
}
