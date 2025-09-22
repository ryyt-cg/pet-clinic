import {
  GalleryVerticalEnd,
  LayoutDashboard,
  ListTodo,
  Package,
  Users,
  Lock,
  MessagesSquare,
  Command, UserX, FileX, ServerOff, Construction,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'satnaing',
    email: 'satnaingdev@gmail.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Pet Clinic',
      logo: Command,
      plan: 'Startup',
    },
    {
      name: 'Acme Inc',
      logo: GalleryVerticalEnd,
      plan: 'Enterprise',
    },
  ],
  navGroups: [
    {
      title: 'General',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'Find Owners',
          url: '/owners',
          icon: ListTodo,
        },
        {
          title: 'Add Owner',
          url: '/add-owner',
          icon: ListTodo,
        },

        {
          title: 'Find Visits',
          url: '/visits',
          icon: Package,
        },
        {
          title: 'Veterinarians',
          url: '/veterinarians',
          icon: MessagesSquare,
        },
        {
          title: 'About',
          url: '/about',
          icon: Users,
        },
      ],
    },
    {
      title: 'Errors',
      items: [
        {
          title: 'Unauthorized',
          url: '/errors/unauthorized',
          icon: Lock,
        },
        {
          title: 'Forbidden',
          url: '/errors/forbidden',
          icon: UserX,
        },
        {
          title: 'Not Found',
          url: '/errors/not-found',
          icon: FileX,
        },
        {
          title: 'Internal Server Error',
          url: '/errors/internal-server-error',
          icon: ServerOff,
        },
        {
          title: 'Maintenance Error',
          url: '/errors/maintenance-error',
          icon: Construction,
        },
      ],
    },
  ],
}
