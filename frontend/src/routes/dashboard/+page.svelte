<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { formatRupiah } from '$lib/utils/formatters';
  import Button from '$lib/components/ui/Button.svelte';

  let summary = $state<any>(null);
  let loading = $state(true);

  async function loadDashboard() {
    loading = true;
    try {
      const res = await api.get('/dashboard/summary');
      summary = res.data;
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat analitik dashboard');
    } finally {
      loading = false;
    }
  }

  onMount(loadDashboard);
</script>

<div class="space-y-6">
  <!-- Top Banner -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Dashboard & Analisis Penjualan</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Ringkasan performa toko, omset harian, dan tren kasir</p>
    </div>
    <div class="flex gap-2">
      <Button variant="accent" onclick={loadDashboard}>
        Refresh Data
      </Button>
      <a href="/pos" class="neo-btn bg-[#FFE600] text-black px-4 py-2.5 text-sm font-black">
        BUKA KASIR (POS)
      </a>
    </div>
  </div>

  {#if loading}
    <div class="p-12 text-center font-black">Memuat ringkasan performa toko...</div>
  {:else if summary}
    <!-- 4 Key Metrics Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Omset Hari Ini -->
      <div class="neo-box p-4 bg-[#FFE600] text-black flex flex-col justify-between">
        <div>
          <span class="text-[11px] font-black uppercase tracking-wider text-black opacity-80">Omset Penjualan Hari Ini</span>
          <div class="text-2xl font-black font-mono mt-1 text-black">
            {formatRupiah(summary.today_revenue)}
          </div>
        </div>
        <div class="mt-4 pt-2 border-t-2 border-black text-xs font-extrabold flex justify-between text-black">
          <span>{summary.today_transaction_count} Transaksi</span>
          <span>Sukses Terbayar</span>
        </div>
      </div>

      <!-- Total Produk -->
      <div class="neo-box p-4 bg-[#00F0FF] text-black flex flex-col justify-between">
        <div>
          <span class="text-[11px] font-black uppercase tracking-wider text-black opacity-80">Total Katalog Produk</span>
          <div class="text-2xl font-black font-mono mt-1 text-black">
            {summary.total_products} SKU
          </div>
        </div>
        <div class="mt-4 pt-2 border-t-2 border-black text-xs font-extrabold flex justify-between text-black">
          <span>Stok Siap Jual</span>
          <a href="/products" class="underline text-black hover:opacity-80">Kelola</a>
        </div>
      </div>

      <!-- Low stock warning -->
      <div class="neo-box p-4 bg-[#FF5252] text-white flex flex-col justify-between">
        <div>
          <span class="text-[11px] font-black uppercase tracking-wider text-white opacity-90">Perhatian Stok Menipis</span>
          <div class="text-2xl font-black font-mono mt-1 text-white">
            {summary.low_stock_count} Barang
          </div>
        </div>
        <div class="mt-4 pt-2 border-t-2 border-white text-xs font-extrabold flex justify-between text-white">
          <span>Stok &le; 10 unit</span>
          <a href="/products" class="underline text-white hover:opacity-80">Restock</a>
        </div>
      </div>

      <!-- Customers -->
      <div class="neo-box p-4 bg-[#B388FF] text-black flex flex-col justify-between">
        <div>
          <span class="text-[11px] font-black uppercase tracking-wider text-black opacity-80">Pelanggan Terdaftar</span>
          <div class="text-2xl font-black font-mono mt-1 text-black">
            {summary.total_customers} Member
          </div>
        </div>
        <div class="mt-4 pt-2 border-t-2 border-black text-xs font-extrabold flex justify-between text-black">
          <span>Program Loyalitas</span>
          <a href="/customers" class="underline text-black hover:opacity-80">Lihat</a>
        </div>
      </div>
    </div>

    <!-- 7-Day Trend Chart -->
    <div class="neo-box p-6 bg-white dark:bg-[#1a1a1a]">
      <div class="flex items-center justify-between mb-6 pb-2 border-b-2 border-black dark:border-white">
        <div>
          <h3 class="text-base font-black uppercase tracking-tight">Tren Penjualan 7 Hari Terakhir</h3>
          <p class="text-xs text-neutral-500 font-bold">Grafik pendapatan omset harian (Rupiah)</p>
        </div>
      </div>

      {#if summary.sales_trend_7days && summary.sales_trend_7days.length > 0}
        {@const maxVal = Math.max(...summary.sales_trend_7days.map((d: any) => d.total), 10000)}
        <div class="grid grid-cols-7 gap-2 sm:gap-4 items-end h-56 pt-6 px-2">
          {#each summary.sales_trend_7days as day}
            {@const heightPct = Math.max(10, Math.round((day.total / maxVal) * 100))}
            <div class="flex flex-col items-center gap-2 h-full justify-end group">
              <!-- Tooltip value -->
              <span class="text-[10px] sm:text-xs font-black font-mono text-center truncate w-full group-hover:scale-110 transition-transform">
                {formatRupiah(day.total)}
              </span>

              <!-- Bar -->
              <div
                class="w-full bg-[#00E676] dark:bg-[#33eb91] border-2 border-black dark:border-white shadow-[3px_3px_0px_0px_#000000] dark:shadow-[3px_3px_0px_0px_#ffffff] transition-all group-hover:bg-[#FFE600]"
                style="height: {heightPct}%;"
              ></div>

              <!-- Day Label -->
              <span class="text-[11px] font-black uppercase tracking-tighter text-neutral-600 dark:text-neutral-300">
                {day.date}
              </span>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Category & Payment Breakdown -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Payment Methods breakdown -->
      <div class="neo-box p-5 bg-white dark:bg-[#1a1a1a]">
        <h3 class="text-sm font-black uppercase tracking-tight mb-4 pb-2 border-b-2 border-black dark:border-white">
          Metode Pembayaran Kasir
        </h3>
        {#if summary.payment_methods && summary.payment_methods.length > 0}
          <div class="space-y-3">
            {#each summary.payment_methods as pm}
              <div class="p-3 border-2 border-black dark:border-neutral-500 bg-neutral-50 dark:bg-[#252525] flex items-center justify-between">
                <div>
                  <span class="text-xs font-black uppercase tracking-wider px-2 py-0.5 bg-black text-white dark:bg-white dark:text-black mr-2">
                    {pm.method}
                  </span>
                  <span class="text-xs font-bold text-neutral-500 dark:text-neutral-400">{pm.count}x Transaksi</span>
                </div>
                <div class="text-sm font-mono font-black text-black dark:text-white">{formatRupiah(pm.total)}</div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-xs text-neutral-500 font-bold">Belum ada transaksi terbayar.</p>
        {/if}
      </div>

      <!-- Category Sales breakdown -->
      <div class="neo-box p-5 bg-white dark:bg-[#1a1a1a]">
        <h3 class="text-sm font-black uppercase tracking-tight mb-4 pb-2 border-b-2 border-black dark:border-white">
          Kontribusi Kategori Terlaris
        </h3>
        {#if summary.category_sales && summary.category_sales.length > 0}
          <div class="space-y-3">
            {#each summary.category_sales as cs}
              <div class="p-3 border-2 border-black dark:border-neutral-500 bg-neutral-50 dark:bg-[#252525] flex items-center justify-between">
                <span class="text-xs font-black uppercase tracking-wider text-black dark:text-white">{cs.category}</span>
                <span class="text-sm font-mono font-black text-[#00E676] dark:text-[#33eb91]">{formatRupiah(cs.total)}</span>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-xs text-neutral-500 font-bold">Belum ada penjualan per kategori.</p>
        {/if}
      </div>
    </div>
  {/if}
</div>
