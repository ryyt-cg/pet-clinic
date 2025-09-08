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
└── types.ts
```

## Migration Steps

