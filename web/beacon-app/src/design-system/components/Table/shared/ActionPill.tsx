import useMenu from '../../../hooks/useMenu';
import mergeClassnames from '../../../utils/mergeClassnames';
import Button from '../../Button/Button';
import { HThreeDotIcon } from '../../Icon/Icons';
import Menu from '../../Menu/Menu';

export type ActionProps = {
  label: string;
  onClick: () => void;
};

export type ActionPillProps = {
  actions: ActionProps[];
  className?: string;
};

export function ActionPill({ actions, className }: ActionPillProps) {
  const { anchorEl, isOpen, open, close } = useMenu({
    id: 'wrapped-menu',
  });

  //console.log('actions: ', actions);

  return (
    <div className={mergeClassnames('relative', className)}>
      <Button
        variant="ghost"
        onClick={open}
        aria-controls={isOpen ? 'wrapped-menu' : undefined}
        aria-expanded={isOpen || undefined}
        aria-haspopup="menu"
        size="custom"
        className="h-[16px] w-[16px] border-none p-0 focus:border-none focus:outline-none focus:ring-0 focus:ring-offset-0"
      >
        <HThreeDotIcon />
      </Button>
      {actions.length > 0 && (
        <>
          <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
            {actions.map((action) => (
              <Menu.Item key={action.label} onClick={action.onClick}>
                {action.label}
              </Menu.Item>
            ))}
          </Menu>
        </>
      )}
    </div>
  );
}
