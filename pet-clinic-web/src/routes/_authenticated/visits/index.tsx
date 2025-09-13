import {createFileRoute} from '@tanstack/react-router'
import Visits from '@/features/visits'

export const Route = createFileRoute('/_authenticated/visits/')({
  component: Visits,
})