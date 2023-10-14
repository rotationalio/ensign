import { t } from '@lingui/macro';

import { CardListItem } from '@/components/common/CardListItem';
import { TEMPLATE_DATA } from '@/features/home/util/utils';

interface TemplateProps {
  title: string;
  links: any[];
}

export const renderTemplate = ({ title, links }: TemplateProps) => {
  return (
    <div className="hover:bg-primary-100/40 flex flex-col gap-2 space-x-2 p-4">
      <h2 className="ml-2 flex font-semibold">{title}</h2>
      {links.map((link) => (
        <a
          href={link.url}
          target="_blank"
          rel="noopener noreferrer"
          key={link}
          className="text-[#1D65A6] hover:text-blue-400"
        >
          {link.name}
        </a>
      ))}
    </div>
  );
};

export default function Templates() {
  return (
    <>
      <CardListItem
        title={t`Templates, Code Examples & Challenges`}
        titleClassName="text-lg"
        className="min-h-[130px]"
        contentClassName="my-2"
      >
        <div>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3 ">
            {TEMPLATE_DATA.map((template) => (
              <div key={template.title}>{renderTemplate(template)}</div>
            ))}
          </div>
        </div>
      </CardListItem>
    </>
  );
}
