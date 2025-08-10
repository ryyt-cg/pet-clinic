import {createRootRoute, Outlet} from '@tanstack/react-router'
import {TanStackRouterDevtools} from '@tanstack/react-router-devtools'
import {Toaster} from "sonner";

export const Route = createRootRoute({
  component: () => (
      <>
        {/*<NavigationProgress />*/}
        <Outlet/>
        <Toaster duration={50000}/>
        {import.meta.env.MODE === 'development' && (
            <>
              {/*<ReactQueryDevtools buttonPosition='bottom-left' />*/}
              <TanStackRouterDevtools position='bottom-right'/>
            </>
        )}
      </>
  ),
})