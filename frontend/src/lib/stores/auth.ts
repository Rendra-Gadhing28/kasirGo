import { writable, derived } from 'svelte/store';
import { api } from '$lib/api/client';
import type { User } from '$lib/types';

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
}

const initialState: AuthState = {
  user: null,
  token: null,
  loading: true
};

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,
    
    init: async () => {
      update((s) => ({ ...s, loading: true }));
      if (typeof window !== 'undefined') {
        const storedToken = localStorage.getItem('kasirpro_token');
        const storedUser = localStorage.getItem('kasirpro_user');
        if (storedToken && storedUser) {
          try {
            const parsedUser = JSON.parse(storedUser);
            update((s) => ({ ...s, token: storedToken, user: parsedUser }));
            
            // Verify session with server
            const res = await api.get<User>('/auth/me');
            if (res.data) {
              update((s) => ({ ...s, user: res.data! }));
              localStorage.setItem('kasirpro_user', JSON.stringify(res.data));
            }
          } catch (e) {
            localStorage.removeItem('kasirpro_token');
            localStorage.removeItem('kasirpro_user');
            update((s) => ({ ...s, token: null, user: null }));
          }
        }
      }
      update((s) => ({ ...s, loading: false }));
    },

    login: async (email: string, password: string) => {
      const res = await api.post('/auth/login', { email, password });
      if (res.data) {
        const userData = res.data.user;
        const accessToken = res.data.access_token;
        update((s) => ({ ...s, user: userData, token: accessToken, loading: false }));
        if (typeof window !== 'undefined') {
          localStorage.setItem('kasirpro_token', accessToken);
          localStorage.setItem('kasirpro_user', JSON.stringify(userData));
        }
      }
      return res;
    },

    register: async (data: { store_name: string; name: string; email: string; password: string; phone?: string }) => {
      const res = await api.post('/auth/register', data);
      if (res.data) {
        const userData = res.data.user;
        const accessToken = res.data.access_token;
        update((s) => ({ ...s, user: userData, token: accessToken, loading: false }));
        if (typeof window !== 'undefined') {
          localStorage.setItem('kasirpro_token', accessToken);
          localStorage.setItem('kasirpro_user', JSON.stringify(userData));
        }
      }
      return res;
    },

    registerStaff: async (data: { outlet_code: string; name: string; email: string; password: string; role?: string }) => {
      const res = await api.post('/auth/register-staff', data);
      if (res.data) {
        const userData = res.data.user;
        const accessToken = res.data.access_token;
        update((s) => ({ ...s, user: userData, token: accessToken, loading: false }));
        if (typeof window !== 'undefined') {
          localStorage.setItem('kasirpro_token', accessToken);
          localStorage.setItem('kasirpro_user', JSON.stringify(userData));
        }
      }
      return res;
    },

    logout: async () => {
      try {
        await api.post('/auth/logout');
      } catch (e) {
        // ignore error
      }
      update((s) => ({ ...s, user: null, token: null, loading: false }));
      if (typeof window !== 'undefined') {
        localStorage.removeItem('kasirpro_token');
        localStorage.removeItem('kasirpro_user');
      }
    }
  };
}

export const auth = createAuthStore();
export const isAuthenticated = derived(auth, ($auth) => $auth.user !== null);
