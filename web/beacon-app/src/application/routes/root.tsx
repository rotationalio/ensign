import { createBrowserRouter, createRoutesFromElements, Route } from 'react-router-dom';

import ErrorPage from '@/components/ErrorPage';
// import routers from features
// should we import all routes files in features folder automatically with a glob pattern?

const Root = () => {
  return (
    <div>
      <h3>Root Router</h3>
    </div>
  );
};

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      path="/"
      element={<Root />}
      //   loader={rootSkeletonLoader}
      //   action={rootAction}
      errorElement={<ErrorPage />}
    >
      {/* <Route errorElement={<ErrorPage />}>
        <Route index element={<Index />} />
        <Route
          path="projects"
          element={<Projects />}
          loader={projectSkeletonLoader}
          action={projectAction}
        />
        <Route
          path="project/:id"
          element={<ProjectId />}
          loader={projectSkeletonLoader}
          action={projectGetAction}
        /> */}
      {/*         
      </Route> */}
    </Route>
  )
);

export default router;
