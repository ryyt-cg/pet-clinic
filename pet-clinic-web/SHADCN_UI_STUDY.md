# Shadcn UI Study
## Purpose:
Documenting steps to migrate and customize Shadcn UI layout from [shadcn-admin](https://github.com/satnaing/shadcn-admin.git) to this project.

## Understanding Shadcn Admin
Directory components contains various UI components and layouts. Key files include:
- `src/components/layout/`: Main layout components (header, sidebar, footer).
- `src/components/ui/`: Shadcn UI components (buttons, modals, forms).
- `src/context/`: Context providers for state management.

Directory src/context in shadcn-admin contains:
``` 
├── direction-provider.tsx
├── font-provider.tsx
├── layout-provider.tsx
├── search-provider.tsx
└── theme-provider.tsx
```


Directory src/components/layout in shadcn-admin contains:
```
├── app-sidebar.tsx
├── app-title.tsx
├── authenticated-layout.tsx
├── data
│   └── sidebar-data.ts
├── header.tsx
├── main.tsx
├── nav-group.tsx
├── nav-user.tsx
├── team-switcher.tsx
├── top-nav.tsx
└── pet-io.ts
```

## Migration Steps
Add Shadcn UI components and context providers to your project. Key steps include:
```bash
pnpm dlx shadcn@latest add alert
pnpm dlx shadcn@latest add alert-dialog
pnpm dlx shadcn@latest add avatar
pnpm dlx shadcn@latest add button
pnpm dlx shadcn@latest add command 
pnpm dlx shadcn@latest add dialog
pnpm dlx shadcn@latest add drawer
pnpm dlx shadcn@latest add dropdown-menu
pnpm dlx shadcn@latest add form
pnpm dlx shadcn@latest add input
pnpm dlx shadcn@latest add label
pnpm dlx shadcn@latest add popover
pnpm dlx shadcn@latest add scroll-area
pnpm dlx shadcn@latest add separator
pnpm dlx shadcn@latest add select
pnpm dlx shadcn@latest add sonner
pnpm dlx shadcn@latest add tooltip


# install react components
pnpm add react-top-loading-bar
pnpm add @radix-ui/react-icons

pnpm install @hookform/resolvers
pnpm add @faker-js/faker --save-dev
pnpm add @radix-ui/react-checkbox


# install tanstack components
pnpm add @tanstack/react-router
pnpm add -D @tanstack/react-router-devtools
pnpm add -D @tanstack/router-plugin

pnpm add @tanstack/react-query
pnmp add #tanstack/react-query-devtools
pnpm add @tanstack/eslint-plugin-query

pnpm add @tanstack/react-form
pnpm add @tanstack/react-table
```
