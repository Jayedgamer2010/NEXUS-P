/// <reference types="vite/client" />

declare module '*.vue' {
    import { DefineComponent } from 'vue';
    const component: DefineComponent<{}, {}, any>;
    export default component;
}

declare module '@vueuse/sound' {
    export function useSound(url: string): { play: () => void };
}

declare module '*.yml' {
    const content: any;
    export default content;
}

declare module 'vite-plugin-eslint';

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_PTERODACTYL_URL: string
  readonly VITE_PTERODACTYL_API_KEY: string
  readonly VITE_APP_NAME: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}