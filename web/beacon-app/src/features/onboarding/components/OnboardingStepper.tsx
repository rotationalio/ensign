const OnboardingStepper = () => {
  return (
    <>
      <ol className="relative border-l border-gray-200 text-gray-500 dark:border-gray-700 dark:text-gray-400">
        <li className="mb-10 ml-6">
          <span className="absolute -left-4 flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white dark:bg-green-900 dark:ring-gray-900"></span>
          <h3 className="font-medium leading-tight">Step 1 of 4</h3>
          <p className="text-sm">Your Team Name</p>
        </li>
        <li className="mb-10 ml-6">
          <span className="absolute -left-4 flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white dark:bg-gray-700 dark:ring-gray-900"></span>
          <h3 className="font-medium leading-tight">Step 2 of 4</h3>
          <p className="text-sm">Your Workspace URL</p>
        </li>
        <li className="mb-10 ml-6">
          <span className="absolute -left-4 flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white dark:bg-gray-700 dark:ring-gray-900"></span>
          <h3 className="font-medium leading-tight">Step 3 of 4</h3>
          <p className="text-sm">Your Name</p>
        </li>
        <li className="ml-6">
          <span className="absolute -left-4 flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white dark:bg-gray-700 dark:ring-gray-900"></span>
          <h3 className="font-medium leading-tight">Step 4 of 4</h3>
          <p className="text-sm">Your Goals</p>
        </li>
      </ol>
    </>
  );
};

export default OnboardingStepper;
