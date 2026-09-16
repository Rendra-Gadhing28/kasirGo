<script lang="ts">
  import { theme } from '$lib/stores/theme';
  import { onMount } from 'svelte';

  interface Props {
    ontogglemobile?: () => void;
  }

  let { ontogglemobile }: Props = $props();

  let currentTime = $state('');

  onMount(() => {
    const updateTime = () => {
      const now = new Date();
      currentTime = now.toLocaleTimeString('id-ID', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      }) + ' WIB';
    };
    updateTime();
    const interval = setInterval(updateTime, 1000);
    return () => clearInterval(interval);
  });
</script>

<header class="h-16 border-b-3 border-black dark:border-white bg-white dark:bg-[#1a1a1a] px-4 lg:px-6 flex items-center justify-between z-30">
  <div class="flex items-center gap-3">
    <button
      onclick={ontogglemobile}
      class="lg:hidden neo-btn bg-[#FFE600] p-2 text-black"
      aria-label="Menu"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>
    <div class="hidden sm:block">
      <span class="text-xs font-black uppercase tracking-wider text-neutral-500 dark:text-neutral-400">Outlet Aktif</span>
      <p class="text-sm font-black text-black dark:text-white">KasirPro Modern Market</p>
    </div>
  </div>

  <div class="flex items-center gap-4">
    <!-- Digital Clock Badge -->
    <div class="border-2 border-black dark:border-white bg-[#B388FF] text-black px-3 py-1 text-xs font-mono font-black shadow-[2px_2px_0px_0px_#000000] dark:shadow-[2px_2px_0px_0px_#ffffff]">
      {currentTime}
    </div>

    <!-- Dark/Light Theme Toggle -->
    <button
      onclick={() => theme.toggle()}
      class="neo-btn bg-white dark:bg-[#222] text-black dark:text-white p-2 border-2 text-sm"
      title="Ubah Mode Gelap / Terang"
      aria-label="Ubah tema warna"
    >
      <span class="dark:hidden font-bold">🌙</span>
      <span class="hidden dark:inline font-bold">☀️</span>
    </button>
  </div>
</header>
