import { writable } from 'svelte/store';

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'info';
  message: string;
}

function createToastStore() {
  const { subscribe, update } = writable<ToastMessage[]>([]);

  return {
    subscribe,
    show: (message: string, type: 'success' | 'error' | 'info' = 'info') => {
      const id = Math.random().toString(36).substring(2, 9);
      update((toasts) => [...toasts, { id, type, message }]);
      setTimeout(() => {
        update((toasts) => toasts.filter((t) => t.id !== id));
      }, 3500);
    },
    success: (msg: string) => {
      toast.show(msg, 'success');
    },
    error: (msg: string) => {
      toast.show(msg, 'error');
    },
    remove: (id: string) => {
      update((toasts) => toasts.filter((t) => t.id !== id));
    }
  };
}

export const toast = createToastStore();
