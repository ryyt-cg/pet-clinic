import {createRootRoute, Outlet} from '@tanstack/react-router'
import { NavigationProgress } from '@/components/navigation-progress'
import {TanStackRouterDevtools} from '@tanstack/react-router-devtools'
import {Toaster} from "sonner";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";
import {NotFoundError} from "@/features/errors/not-found-error.tsx";
import {GeneralError} from "@/features/errors/general-error.tsx";

export const Route = createRootRoute({
  component: () => (
      <>
        <NavigationProgress />
        <Outlet />
        <Toaster duration={5000} />
        {import.meta.env.MODE === 'development' && (
            <>
              <ReactQueryDevtools buttonPosition='bottom-left' />
              <TanStackRouterDevtools position='bottom-right'/>
            </>
        )}
      </>
  ),
  notFoundComponent: NotFoundError,
  errorComponent: GeneralError,
})