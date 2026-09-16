import { writable } from 'svelte/store';

function createThemeStore() {
  const isDark = writable<boolean>(false);

  return {
    subscribe: isDark.subscribe,
    init: () => {
      if (typeof window !== 'undefined') {
        const stored = localStorage.getItem('kasirpro_theme');
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const activeDark = stored ? stored === 'dark' : prefersDark;
        isDark.set(activeDark);
        if (activeDark) {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
      }
    },
    toggle: () => {
      isDark.update((current) => {
        const next = !current;
        if (typeof window !== 'undefined') {
          localStorage.setItem('kasirpro_theme', next ? 'dark' : 'light');
          if (next) {
            document.documentElement.classList.add('dark');
          } else {
            document.documentElement.classList.remove('dark');
          }
        }
        return next;
      });
    }
  };
}

export const theme = createThemeStore();
