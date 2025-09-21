import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";
import {Main} from "@/components/layout/main.tsx";
import {GetPetById} from "@/gateway/pet-gateway.ts";

const Owners = () => {
  const petById = GetPetById(1)
  if (petById.isPending) {
    return <span>Loading...</span>;
  }

  return (
      <>
        <Header fixed>
          <div className='ms-auto flex items-center space-x-4'>
            <Search />
            <ThemeSwitch />
            <ConfigDrawer />
          </div>
        </Header>

        <Main>
          <div className='mb-2 flex items-center justify-between space-y-2'>
            <h1 className='text-2xl font-bold tracking-tight'>Owners Page</h1>
          </div>
          <div>
            <p>Query function status: {petById.fetchStatus}</p>
            <p>Query data status: {petById.status}</p>
            {petById.data && (
                <ul key={petById.data.id}>
                  <li>Pet ID: {petById.data.id}</li>
                  <li>Pet Name: {petById.data.name}</li>
                  <li>Pet Birth Date: {petById.data.birthdate}</li>
                </ul>
            )}
          </div>
        </Main>
      </>
  );
};

export default Owners;