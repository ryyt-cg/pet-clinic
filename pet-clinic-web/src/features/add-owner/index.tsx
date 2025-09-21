import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";
import {Main} from "@/components/layout/main.tsx";
import {GetAllPets} from "@/gateway/pet-gateway.ts";

const AddOwner = () => {
  const allPets = GetAllPets()
  if (allPets.isPending) {
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
            <h1 className='text-2xl font-bold tracking-tight'>Add Owner</h1>
          </div>
          <div>
            <p>Query function status: {allPets.fetchStatus}</p>
            <p>Query data status: {allPets.status}</p>
            {allPets.data?.pets.map((pet) => (
                <ul key={pet.id}>
                  <li>Pet ID: {pet.id}</li>
                  <li>Pet Name: {pet.name}</li>
                  <li>Pet Birth Date: {pet.birthdate}</li>
                </ul>
            ))}
          </div>
        </Main>

      </>
  );
};

export default AddOwner;