import { Heading } from '@rotational/beacon-core';
import cn from 'classnames';
import { Link } from 'react-router-dom';

import { appConfig } from '@/application/config';
import ExternalIcon from '@/components/icons/external-icon';
import { footerItems } from '@/constants/dashLayout';
import OnboardingStepper from '@/features/onboarding/components/OnboardingStepper';
type SidebarProps = {
  className?: string;
};

function OnboardingSideBar({ className }: SidebarProps) {
  const { version: appVersion, revision: gitRevision } = appConfig;

  return (
    <>
      <aside
        className={cn(
          `fixed top-0 left-0 flex h-screen flex-col bg-[#1D65A6] pt-5 pb-10 text-white md:w-[250px]`,
          className
        )}
      >
        <div className="flex h-full flex-col">
          <div>
            <Heading as="h1" className="ml-8 space-y-3 text-lg font-bold">
              Ensign
            </Heading>
          </div>
          <div className="ml-8 mt-12 space-y-3">
            <OnboardingStepper />
          </div>
          <div className="grow"></div>
          <div className="ml-8 space-y-3">
            <ul className="space-y-1 text-xs text-white">
              {footerItems.map((item) => (
                <li key={`${item.name}`}>
                  <Link to={item.href} target="_blank" className="flex">
                    {item.name}{' '}
                    {item.isExternal && <ExternalIcon className="ml-1 h-3 w-3 text-white" />}
                  </Link>
                </li>
              ))}
            </ul>
            <p className="w-full text-xs text-white">
              {appVersion && <span>App Version {appVersion} </span>}
              {gitRevision && <span>& Git Revision {gitRevision} </span>}
            </p>
          </div>
        </div>
      </aside>
    </>
  );
}

export default OnboardingSideBar;
