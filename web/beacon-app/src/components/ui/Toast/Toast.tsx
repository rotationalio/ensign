import * as ToastPrimitive from '@radix-ui/react-toast';
import * as React from 'react';
import { twMerge } from 'tailwind-merge';

import { ErrorIcon, InfoIcon, SuccessIcon, WarningIcon, XIcon } from '@/components/icons/Icons';

import { ToastWithRadixProps } from './Toast.types';
import { getToastBgVariantStyle, getToastColorVariantStyle } from './util';
export const Toast = React.forwardRef<HTMLLIElement, ToastWithRadixProps>((props, ref) => {
  const {
    children,
    variant,
    title,
    description,
    duration = 5000,
    isOpen = false,
    onClose,
    hasIcon = false,
    placement = 'right',
  } = props;

  const [isToastOpen, setIsToastOpen] = React.useState(isOpen);

  React.useEffect(() => {
    let timeoutId: any;
    if (duration && isToastOpen) {
      timeoutId = setTimeout(() => {
        setIsToastOpen(false);
      }, duration);
    }

    return () => {
      clearTimeout(timeoutId);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <ToastPrimitive.Provider swipeDirection={placement} duration={duration}>
        <ToastPrimitive.Root
          open={isToastOpen}
          ref={ref}
          duration={duration}
          onOpenChange={onClose}
          className={twMerge(
            'shadow-lg fixed inset-x-4 bottom-4 z-50 w-auto rounded-lg md:top-4 md:right-4 md:left-auto md:bottom-auto md:w-full md:max-w-sm',
            'bg-white',
            variant && getToastBgVariantStyle(variant),
            'radix-state-open:animate-toast-slide-in-bottom md:radix-state-open:animate-toast-slide-in-right',
            'radix-state-closed:animate-toast-hide',
            'radix-swipe-direction-right:radix-swipe-end:animate-toast-swipe-out-x',
            'radix-swipe-direction-right:translate-x-radix-toast-swipe-move-x',
            'radix-swipe-direction-down:radix-swipe-end:animate-toast-swipe-out-y',
            'radix-swipe-direction-down:translate-y-radix-toast-swipe-move-y',
            'radix-swipe-cancel:translate-x-0 radix-swipe-cancel:duration-200 radix-swipe-cancel:ease-[ease]',
            'focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75'
          )}
        >
          <div className="flex ">
            {hasIcon && (
              <div className="flex flex-shrink-0 items-center py-4 pl-5">
                {
                  {
                    success: <SuccessIcon className="h-5 w-5 text-white" />,
                    danger: <ErrorIcon className="h-5 w-5 text-white" />,
                    warning: <WarningIcon className="h-5 w-5 text-white" />,
                    primary: <InfoIcon className="h-5 w-5 text-white" />,
                    info: <InfoIcon className="h-5 w-5 text-white" />,
                    secondary: <></>,
                    default: <></>,
                  }[variant || 'default']
                }
              </div>
            )}
            <div className="w-0 flex-1 items-center py-4 pl-5">
              <div className="radix w-full">
                <ToastPrimitive.Title
                  className={twMerge(
                    'text-sm font-medium text-gray-900',
                    variant && getToastColorVariantStyle(variant)
                  )}
                >
                  {title}
                </ToastPrimitive.Title>
                <ToastPrimitive.Description
                  className={twMerge(
                    'mt-1 text-sm text-gray-700 ',
                    variant && getToastColorVariantStyle(variant)
                  )}
                >
                  {description}
                  {children}
                </ToastPrimitive.Description>
              </div>
            </div>
            <div className="flex">
              <div className="flex flex-col space-y-1 px-2">
                <div className="flex h-0 flex-1">
                  {onClose && (
                    <ToastPrimitive.Close
                      className={twMerge(
                        'border-transparent  hover:bg-gray-50 flex w-full items-center justify-center px-3 py-2 text-sm font-medium  text-gray-700  focus:z-10 focus:outline-none focus-visible:ring focus-visible:ring-purple-500 focus-visible:ring-opacity-75',
                        variant && getToastColorVariantStyle(variant)
                      )}
                    >
                      <XIcon
                        className={twMerge(
                          'h-5 w-5',
                          variant && getToastColorVariantStyle(variant)
                        )}
                      />
                    </ToastPrimitive.Close>
                  )}
                </div>
              </div>
            </div>
          </div>
        </ToastPrimitive.Root>

        <ToastPrimitive.Viewport />
      </ToastPrimitive.Provider>
    </>
  );
});

Toast.displayName = 'Toast';

export default Toast;
