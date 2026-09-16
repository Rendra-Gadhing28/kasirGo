<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { theme } from '$lib/stores/theme';
  import Sidebar from '$lib/components/layout/Sidebar.svelte';
  import Navbar from '$lib/components/layout/Navbar.svelte';
  import Toast from '$lib/components/ui/Toast.svelte';

  let { children } = $props();
  let mobileSidebarOpen = $state(false);

  onMount(async () => {
    theme.init();
    await auth.init();
  });

  const isPublicStandalone = $derived(
    page.url.pathname === '/login' ||
    page.url.pathname === '/register' ||
    page.url.pathname.startsWith('/join-member')
  );

  $effect(() => {
    if (!$auth.loading) {
      if (!$auth.user && !isPublicStandalone) {
        goto('/login');
      } else if ($auth.user && (page.url.pathname === '/login' || page.url.pathname === '/register' || page.url.pathname === '/')) {
        goto('/pos');
      }
    }
  });
</script>

<Toast />

{#if isPublicStandalone}
  <main class="min-h-screen">
    {@render children()}
  </main>
{:else if $auth.loading}
  <div class="min-h-screen flex flex-col items-center justify-center bg-[#FFFDF7] dark:bg-[#121212] gap-4">
    <div class="w-12 h-12 border-4 border-black dark:border-white border-t-[#FFE600] animate-spin"></div>
    <p class="font-black text-sm uppercase tracking-widest">Memuat KasirPro...</p>
  </div>
{:else}
  <div class="flex h-screen overflow-hidden bg-[#FFFDF7] dark:bg-[#121212]">
    <!-- Left Navigation -->
    <Sidebar
      mobileOpen={mobileSidebarOpen}
      onclose={() => (mobileSidebarOpen = false)}
    />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <Navbar ontogglemobile={() => (mobileSidebarOpen = !mobileSidebarOpen)} />
      <main class="flex-1 overflow-y-auto p-4 md:p-6">
        {@render children()}
      </main>
    </div>
  </div>
{/if}
