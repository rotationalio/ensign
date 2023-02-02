import { routes } from "@/application";
import { Route } from "react-router-dom";
import WelcomePage from "./WelcomePage";
import MainLayout from "@/components/layout/MainLayout";

export default function WelcomePageLayout () {
    return (
    <>
    <Route path={routes.welcome} element={<MainLayout />}>
        <Route
          path={routes.welcome}
          element={<WelcomePage />}
        />
    </Route>
    </>
    );
  };
  