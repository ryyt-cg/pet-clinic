import {createFileRoute} from '@tanstack/react-router'
import Owners from '@/features/add-owner'

export const Route = createFileRoute('/_authenticated/add-owner/')({
  component: Owners,
})