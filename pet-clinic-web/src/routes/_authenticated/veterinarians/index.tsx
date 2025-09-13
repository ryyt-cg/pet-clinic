import {createFileRoute} from '@tanstack/react-router'
import Veterinarians from '@/features/veterinarians'

export const Route = createFileRoute('/_authenticated/veterinarians/')({
  component: Veterinarians,
})