import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";
import {Main} from "@/components/layout/main.tsx";

const Veterinarians = () => {
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
            <h1 className='text-2xl font-bold tracking-tight'>Veterinarians Page</h1>
          </div>
        </Main>
      </>
  );
};

export default Veterinarians;