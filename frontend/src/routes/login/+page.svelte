<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { toast } from '$lib/stores/toast';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let email = $state('owner@kasirpro.id');
  let password = $state('password123');
  let loading = $state(false);
  let errorMessage = $state('');

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    loading = true;
    errorMessage = '';
    try {
      await auth.login(email, password);
      toast.success('Selamat datang kembali di KasirPro!');
      goto('/pos');
    } catch (err: any) {
      errorMessage = err.message || 'Login gagal. Periksa kembali email dan password.';
      toast.error(errorMessage);
    } finally {
      loading = false;
    }
  }

  function setDemoAccount(role: 'owner' | 'admin' | 'kasir') {
    if (role === 'owner') {
      email = 'owner@kasirpro.id';
    } else if (role === 'admin') {
      email = 'admin@kasirpro.id';
    } else {
      email = 'kasir@kasirpro.id';
    }
    password = 'password123';
  }
</script>

<div class="min-h-screen flex items-center justify-center p-4 bg-[#FFFDF7] dark:bg-[#121212]">
  <div class="w-full max-w-md bg-white dark:bg-[#1e1e1e] border-4 border-black dark:border-white shadow-[8px_8px_0px_0px_#000000] dark:shadow-[8px_8px_0px_0px_#ffffff] p-6 sm:p-8">
    
    <!-- Branding Header -->
    <div class="text-center mb-6">
      <div class="inline-flex items-center gap-1.5 bg-[#FFE600] px-3 py-1 border-3 border-black shadow-[3px_3px_0px_0px_#000000] mb-3">
        <span class="text-lg font-black tracking-tight text-black">KASIRPRO POS</span>
      </div>
      <h2 class="text-2xl font-black uppercase tracking-tight text-black dark:text-white">Masuk Sistem</h2>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400 mt-1">Sistem Kasir Pasar Modern Indonesia</p>
    </div>

    <!-- Quick Demo Logins -->
    <div class="mb-6 p-3 bg-neutral-100 dark:bg-[#282828] border-2 border-black dark:border-neutral-400">
      <p class="text-[10px] font-black uppercase tracking-wider text-neutral-600 dark:text-neutral-300 mb-2">Akun Demo Cepat:</p>
      <div class="grid grid-cols-3 gap-2">
        <button
          type="button"
          onclick={() => setDemoAccount('owner')}
          class="neo-btn bg-[#FFE600] text-black py-1.5 text-xs font-black"
        >
          Owner
        </button>
        <button
          type="button"
          onclick={() => setDemoAccount('admin')}
          class="neo-btn bg-[#B388FF] text-black py-1.5 text-xs font-black"
        >
          Admin
        </button>
        <button
          type="button"
          onclick={() => setDemoAccount('kasir')}
          class="neo-btn bg-[#00F0FF] text-black py-1.5 text-xs font-black"
        >
          Kasir
        </button>
      </div>
    </div>

    <!-- Login Form -->
    <form onsubmit={handleLogin} class="space-y-4">
      <Input
        label="Email Kasir / Toko"
        type="email"
        placeholder="nama@toko.com"
        bind:value={email}
        required
      />

      <Input
        label="Kata Sandi (Password)"
        type="password"
        placeholder="••••••••"
        bind:value={password}
        required
      />

      {#if errorMessage}
        <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 text-xs font-bold">
          {errorMessage}
        </div>
      {/if}

      <Button
        type="submit"
        variant="primary"
        size="lg"
        class="w-full mt-2"
        loading={loading}
      >
        MASUK KE SISTEM KASIR
      </Button>
    </form>

    <!-- Register Link -->
    <div class="mt-6 pt-4 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700 text-center">
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">
        Belum punya akun toko? 
        <a href="/register" class="text-black dark:text-[#FFE600] font-black underline hover:text-[#FFE600]">
          Daftar Toko Baru
        </a>
      </p>
    </div>
  </div>
</div>
