import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";
import {Main} from "@/components/layout/main.tsx";
import {GetAllVets, GetVetById, GetVetWithSpecialtiesById} from "@/gateway/vet-gateway.ts";
import {GetAllOwners, GetOwnerById, GetOwnerWithPetsById} from "@/gateway/owner-gateway.ts";
import {GetPetWithVisitsById} from "@/gateway/pet-gateway.ts";
import {GetAllVisits, GetVisitById} from "@/gateway/visit-gateway.ts";
import {NotFoundError} from "@/features/errors/not-found-error.tsx";

const Veterinarians = () => {
  const allOwners = GetAllOwners()
  const ownerById = GetOwnerById(200)
  const ownerWithPetsById = GetOwnerWithPetsById(2)

  const petWithVisitsById = GetPetWithVisitsById(7)

  const allVets = GetAllVets()
  const vetById = GetVetById(2)
  const vetWithSpecialtiesById = GetVetWithSpecialtiesById(2)

  const allVisits = GetAllVisits()
  const visitById = GetVisitById(2)

  if (vetById.isPending || allVets.isPending || vetWithSpecialtiesById.isPending) {
    return <span>Loading...</span>;
  }

  if (ownerById.isPending || allOwners.isPending || ownerWithPetsById.isPending || petWithVisitsById.isPending) {
    return <span>Loading...</span>;
  }

  if (allVisits.isPending || visitById.isPending) {
    return <span>Loading...</span>;
  }

  if (ownerById.isError) {
    return <NotFoundError/>;
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
            <h1 className='text-2xl font-bold tracking-tight'>Veterinarians Page</h1>
          </div>
        </Main>
      </>
  );
};

export default Veterinarians;