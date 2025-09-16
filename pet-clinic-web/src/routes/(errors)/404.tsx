import { createFileRoute } from '@tanstack/react-router'
import { NotFoundError } from '@/features/errors/not-found-error.tsx'

export const Route = createFileRoute('/(errors)/404')({
  component: NotFoundError,
})
