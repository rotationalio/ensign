import { Button, Menu, useMenu } from '@rotational/beacon-core';
import { Fragment } from 'react';

import SettingIcon from '@/components/icons/setting';

export type SettingsDataProps = {
  name: string;
  onClick: () => void;
};

interface SettingsProps {
  key?: any;
  data: SettingsDataProps[];
}
const SettingsButton = ({ data, key }: SettingsProps) => {
  const k = key || Math.random().toString(36).substring(7);
  const { isOpen, close, open, anchorEl } = useMenu({ id: k });

  return (
    <>
      <div>
        <Button
          variant="ghost"
          size="custom"
          className="flex-end bg-inherit hover:bg-transparent border-none"
          onClick={open}
          data-cy="detailActions"
        >
          <SettingIcon />
        </Button>
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          {data.map((item: SettingsDataProps, idx: any) => (
            <Fragment key={idx}>
              <Menu.Item onClick={item.onClick} data-testid="cancelButton">
                {item.name}
              </Menu.Item>
            </Fragment>
          ))}
        </Menu>
      </div>
    </>
  );
};

export default SettingsButton;
