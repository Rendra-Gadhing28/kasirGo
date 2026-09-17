<script lang="ts">
  import { page } from '$app/state';
  import { auth } from '$lib/stores/auth';
  import Badge from '$lib/components/ui/Badge.svelte';

  interface Props {
    mobileOpen?: boolean;
    onclose?: () => void;
  }

  let { mobileOpen = false, onclose }: Props = $props();

  const navItems = [
    { label: 'Kasir (POS)', href: '/pos', icon: 'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z' },
    { label: 'Dashboard', href: '/dashboard', icon: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z' },
    { label: 'Produk', href: '/products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' },
    { label: 'Kategori', href: '/categories', icon: 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z' },
    { label: 'Pelanggan', href: '/customers', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
    { label: 'Supplier', href: '/suppliers', icon: 'M8 4H6a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-2m-4-1v8m0 0l3-3m-3 3L9 8m-5 5h2.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293h3.172a1 1 0 00.707-.293l2.414-2.414a1 1 0 01.707-.293H20' },
    { label: 'Riwayat Transaksi', href: '/transactions', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01' },
  ];
</script>

<!-- Mobile backdrop -->
{#if mobileOpen}
  <button
    type="button"
    class="fixed inset-0 bg-black/60 z-40 lg:hidden cursor-default"
    onclick={onclose}
    aria-label="Tutup menu samping"
  ></button>
{/if}

<aside
  class="fixed lg:static top-0 bottom-0 left-0 z-40 w-64 bg-white dark:bg-[#1a1a1a] border-r-3 border-black dark:border-white flex flex-col justify-between transition-transform duration-200 lg:translate-x-0 {mobileOpen ? 'translate-x-0' : '-translate-x-full'}"
>
  <div>
    <!-- Logo Header -->
    <div class="p-5 border-b-3 border-black dark:border-white flex items-center justify-between bg-[#FFE600] text-black">
      <div>
        <h1 class="text-xl font-black tracking-tighter uppercase flex items-center gap-1.5">
          <span>KASIR</span><span class="bg-black text-white px-1.5 py-0.5 border-2 border-black text-xs font-black">PRO</span>
        </h1>
        <p class="text-[10px] font-extrabold uppercase tracking-widest text-neutral-800">Indonesian POS</p>
      </div>
      <button
        onclick={onclose}
        class="lg:hidden w-7 h-7 flex items-center justify-center bg-black text-white border-2 border-black font-black text-sm"
      >
        ×
      </button>
    </div>

    <!-- Nav links -->
    <nav class="p-3 space-y-1.5">
      {#each navItems as item}
        {@const isActive = page.url.pathname === item.href}
        <a
          href={item.href}
          onclick={onclose}
          class="flex items-center gap-3 px-3.5 py-2.5 font-black text-sm border-2 transition-all {isActive
            ? 'bg-[#00F0FF] text-black border-black dark:border-white shadow-[3px_3px_0px_0px_#000000] dark:shadow-[3px_3px_0px_0px_#ffffff] translate-x-1'
            : 'border-transparent text-neutral-800 dark:text-neutral-200 hover:border-black dark:hover:border-white hover:bg-neutral-100 dark:hover:bg-neutral-800 hover:text-black dark:hover:text-white'}"
        >
          <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d={item.icon} />
          </svg>
          <span>{item.label}</span>
        </a>
      {/each}
    </nav>
  </div>

  <!-- Bottom User Profile & Logout -->
  <div class="p-4 border-t-3 border-black dark:border-white bg-neutral-50 dark:bg-[#151515]">
    {#if $auth.user}
      <div class="mb-3">
        <div class="flex items-center justify-between gap-1">
          <p class="text-sm font-black truncate max-w-[130px] text-black dark:text-white">{$auth.user.name}</p>
          <Badge variant={$auth.user.role === 'owner' ? 'yellow' : $auth.user.role === 'admin' ? 'purple' : 'cyan'}>
            {$auth.user.role}
          </Badge>
        </div>
        <p class="text-xs text-neutral-500 dark:text-neutral-400 font-bold truncate">{$auth.user.email}</p>
      </div>
    {/if}

    <button
      onclick={() => auth.logout()}
      class="w-full neo-btn bg-[#FF5252] text-white py-2 text-xs font-black uppercase tracking-wider"
    >
      Keluar (Logout)
    </button>
  </div>
</aside>
