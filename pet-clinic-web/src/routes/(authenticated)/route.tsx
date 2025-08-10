import {createFileRoute, Outlet} from "@tanstack/react-router";
// import Layout from "@/components/layout/layout.tsx";

export const Route = createFileRoute('/(authenticated)')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
      <>
        {/*<Layout>*/}
        <Outlet/>
        {/*</Layout>*/}
      </>
  );
}