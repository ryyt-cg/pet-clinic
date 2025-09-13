import {createFileRoute} from '@tanstack/react-router'
import Owners from '@/features/owners'

export const Route = createFileRoute('/_authenticated/owners/')({
  component: Owners,
})