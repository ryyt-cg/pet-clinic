import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import { DirectionProvider } from './context/direction-provider'
import '@/index.css'
import App from "@/app.tsx";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <DirectionProvider>
        <App/>
      </DirectionProvider>
    </StrictMode>,
)
