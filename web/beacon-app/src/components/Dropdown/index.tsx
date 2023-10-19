import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu';
import React from 'react';
interface RadixDropdownMenuProps {
  items: any;
  trigger?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (isOpen: boolean) => void;
}

// trigger is the button that opens the dropdown and will be use as children of generic dropdown
const Trigger = ({ trigger }: RadixDropdownMenuProps) => {
  return (
    <DropdownMenuPrimitive.Trigger>
      <button className="border-none focus:ring-0">{trigger}</button>
    </DropdownMenuPrimitive.Trigger>
  );
};

// this component can be used to build any dropdown
const GenericDropDown = ({ items, trigger, isOpen, onOpenChange }: RadixDropdownMenuProps) => {
  return (
    <div className="relative">
      <DropdownMenuPrimitive.Root
        open={isOpen}
        onOpenChange={onOpenChange}
      ></DropdownMenuPrimitive.Root>
    </div>
  );
};

export default GenericDropDown;
