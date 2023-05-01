import { t, Trans } from '@lingui/macro';
import { Button, Card } from '@rotational/beacon-core';

import { EXTRENAL_LINKS } from '@/application';
import { CardListItem } from '@/components/common/CardListItem';

export default function AccessDocumentationStep() {
  return (
    <>
      <CardListItem title={t`Step 3: View Documentation`}>
        <Card.Body>
          <div className="mt-5 flex flex-col gap-8 px-3 xl:flex-row">
            <p className="w-full text-sm sm:w-4/5">
              <Trans>
                Love seeing examples with real code? Prefer watching tutorial videos? Still learning
                the basics? We’ve got you covered!
              </Trans>
            </p>
            <div className="sm:w-1/5">
              <a
                href={EXTRENAL_LINKS.DOCUMENTATION}
                target="_blank"
                rel="noopener noreferrer"
                data-testid="viewDocsLink"
              >
                <Button className="h-[44px] w-[165px]" color="primary" data-testid="viewDocs">
                  <Trans>View Docs</Trans>
                </Button>
              </a>
            </div>
          </div>
        </Card.Body>
      </CardListItem>
    </>
  );
}
