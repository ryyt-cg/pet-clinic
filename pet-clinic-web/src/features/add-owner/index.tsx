import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";

const AddOwner = () => {
  return (
      <>
        <Header fixed>
          <div className='ms-auto flex items-center space-x-4'>
            <Search />
            <ThemeSwitch />
            <ConfigDrawer />
          </div>
        </Header>

        Add Owner Page
      </>
  );
};

export default AddOwner;