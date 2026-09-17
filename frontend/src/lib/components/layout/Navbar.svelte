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
      class="neo-btn bg-white dark:bg-[#222] text-black dark:text-white px-2.5 py-1.5 border-2 text-sm flex items-center gap-1.5 font-black"
      title="Ubah Mode Gelap / Terang"
      aria-label="Ubah tema warna"
    >
      <svg class="w-4 h-4 dark:hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
      <svg class="w-4 h-4 hidden dark:inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <span class="text-xs uppercase hidden sm:inline">
        <span class="dark:hidden">Gelap</span>
        <span class="hidden dark:inline">Terang</span>
      </span>
    </button>
  </div>
</header>
