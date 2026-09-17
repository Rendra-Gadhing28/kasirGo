<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah, formatDateTime } from '$lib/utils/formatters';
  import type { Transaction } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import ThermalReceipt from '$lib/components/pos/ThermalReceipt.svelte';

  let transactions = $state<Transaction[]>([]);
  let loading = $state(true);

  // Filters
  let search = $state('');
  let status = $state('');
  let startDate = $state('');
  let endDate = $state('');
  let page = $state(1);
  let totalPages = $state(1);
  let totalCount = $state(0);

  // Selected for Receipt
  let selectedTrans = $state<Transaction | null>(null);
  let receiptModalOpen = $state(false);

  // Void modal
  let voidModalOpen = $state(false);
  let voidTransId = $state<number | null>(null);
  let voidReason = $state('');
  let voidLoading = $state(false);

  async function loadTransactions() {
    loading = true;
    try {
      let query = `/transactions?page=${page}&limit=20`;
      if (search) query += `&search=${encodeURIComponent(search)}`;
      if (status) query += `&status=${status}`;
      if (startDate) query += `&start_date=${startDate}`;
      if (endDate) query += `&end_date=${endDate}`;

      const res = await api.get<Transaction[]>(query);
      transactions = res.data || [];
      if (res.meta) {
        totalPages = res.meta.total_pages || 1;
        totalCount = res.meta.total || 0;
      }
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat riwayat transaksi');
    } finally {
      loading = false;
    }
  }

  function viewReceipt(t: Transaction) {
    selectedTrans = t;
    receiptModalOpen = true;
  }

  function openVoidModal(id: number) {
    voidTransId = id;
    voidReason = 'Kesalahan input kasir / retur pembeli';
    voidModalOpen = true;
  }

  async function handleVoid() {
    if (!voidTransId) return;
    voidLoading = true;
    try {
      await api.post(`/transactions/${voidTransId}/void`, {
        reason: voidReason
      });
      toast.success('Transaksi berhasil dibatalkan dan stok dikembalikan!');
      voidModalOpen = false;
      await loadTransactions();
    } catch (e: any) {
      toast.error(e.message || 'Gagal membatalkan transaksi');
    } finally {
      voidLoading = false;
    }
  }

  onMount(loadTransactions);
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Riwayat & Laporan Transaksi</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Total {totalCount} transaksi tercatat di database</p>
    </div>
    <div class="flex gap-2">
      <Button variant="accent" onclick={loadTransactions}>
        Refresh
      </Button>
    </div>
  </div>

  <!-- Search & Filter Controls -->
  <div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
    <div>
      <Input
        placeholder="Cari No. Faktur (INV-...)"
        bind:value={search}
        onkeydown={(e) => { if (e.key === 'Enter') { page = 1; loadTransactions(); } }}
      />
    </div>
    <div>
      <select
        bind:value={status}
        onchange={() => { page = 1; loadTransactions(); }}
        class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
      >
        <option value="">Semua Status Bayar</option>
        <option value="paid">Lunas (Paid)</option>
        <option value="pending">Menunggu (Pending)</option>
        <option value="cancelled">Dibatalkan (Void)</option>
      </select>
    </div>
    <div>
      <input
        type="date"
        bind:value={startDate}
        onchange={() => { page = 1; loadTransactions(); }}
        class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
      />
    </div>
    <div class="flex gap-2">
      <input
        type="date"
        bind:value={endDate}
        onchange={() => { page = 1; loadTransactions(); }}
        class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
      />
      <Button variant="primary" onclick={() => { page = 1; loadTransactions(); }}>
        Filter
      </Button>
    </div>
  </div>

  <!-- Transactions Table -->
  <div class="neo-box overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-[#FFE600] text-black border-b-3 border-black dark:border-white font-black text-xs uppercase tracking-wider">
            <th class="p-3.5">No Faktur</th>
            <th class="p-3.5">Waktu</th>
            <th class="p-3.5">Kasir</th>
            <th class="p-3.5">Pelanggan</th>
            <th class="p-3.5">Metode</th>
            <th class="p-3.5 text-center">Status</th>
            <th class="p-3.5 text-right">Total</th>
            <th class="p-3.5 text-center w-36">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y-2 divide-black dark:divide-neutral-700 text-sm font-bold text-black dark:text-white">
          {#if loading}
            <tr><td colspan="8" class="p-8 text-center text-neutral-500 font-bold">Memuat transaksi...</td></tr>
          {:else if transactions.length === 0}
            <tr><td colspan="8" class="p-8 text-center text-neutral-500 font-bold">Tidak ada riwayat transaksi ditemukan.</td></tr>
          {:else}
            {#each transactions as t (t.id)}
              <tr class="hover:bg-yellow-50 dark:hover:bg-[#252525] transition-colors">
                <td class="p-3.5 font-mono font-black text-xs">
                  {t.invoice_no}
                </td>
                <td class="p-3.5 text-xs text-neutral-600 dark:text-neutral-400 whitespace-nowrap">
                  {formatDateTime(t.created_at)}
                </td>
                <td class="p-3.5 text-xs">{t.user?.name || '-'}</td>
                <td class="p-3.5 text-xs">{t.customer?.name || 'Walk-in'}</td>
                <td class="p-3.5 uppercase font-mono text-xs">{t.payment_method}</td>
                <td class="p-3.5 text-center">
                  {#if t.payment_status === 'paid'}
                    <Badge variant="green">Lunas</Badge>
                  {:else if t.payment_status === 'pending'}
                    <Badge variant="yellow">Pending</Badge>
                  {:else}
                    <Badge variant="red">Batal</Badge>
                  {/if}
                </td>
                <td class="p-3.5 text-right font-black font-mono text-sm whitespace-nowrap">
                  {formatRupiah(t.total_amount)}
                </td>
                <td class="p-3.5 text-center">
                  <div class="inline-flex gap-1.5">
                    <button
                      onclick={() => viewReceipt(t)}
                      class="px-2 py-1 text-xs font-black bg-[#00F0FF] text-black border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                    >
                      Struk
                    </button>
                    {#if t.payment_status === 'paid' && $auth.user?.role !== 'kasir'}
                      <button
                        onclick={() => openVoidModal(t.id)}
                        class="px-2 py-1 text-xs font-black bg-[#FF5252] text-white border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                      >
                        Void
                      </button>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    {#if totalPages > 1}
      <div class="p-4 border-t-3 border-black dark:border-white bg-neutral-50 dark:bg-[#1f1f1f] flex items-center justify-between">
        <span class="text-xs font-bold text-neutral-600 dark:text-neutral-400">
          Halaman {page} dari {totalPages}
        </span>
        <div class="flex gap-2">
          <Button
            size="sm"
            variant="outline"
            disabled={page <= 1}
            onclick={() => { page--; loadTransactions(); }}
          >
            Sebelumnya
          </Button>
          <Button
            size="sm"
            variant="outline"
            disabled={page >= totalPages}
            onclick={() => { page++; loadTransactions(); }}
          >
            Berikutnya
          </Button>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- Modal Struk -->
<Modal
  open={receiptModalOpen}
  title="STRUK TRANSAKSI"
  maxWidth="max-w-md"
  onclose={() => (receiptModalOpen = false)}
>
  {#if selectedTrans}
    <ThermalReceipt
      transaction={selectedTrans}
      onclose={() => (receiptModalOpen = false)}
    />
  {/if}
</Modal>

<!-- Modal Void / Retur -->
<Modal
  open={voidModalOpen}
  title="BATALKAN TRANSAKSI (VOID)"
  maxWidth="max-w-md"
  onclose={() => (voidModalOpen = false)}
>
  <div class="space-y-4">
    <p class="text-sm font-bold text-red-600">
      Peringatan: Membatalkan transaksi akan mengembalikan stok seluruh barang ke inventaris dan memotong poin loyalitas terkait.
    </p>

    <div class="flex flex-col gap-1.5">
      <label for="void-reason" class="text-xs font-black uppercase">Alasan Pembatalan / Void</label>
      <textarea
        id="void-reason"
        rows="3"
        bind:value={voidReason}
        placeholder="Tulis alasan..."
        class="neo-input"
        required
      ></textarea>
    </div>

    <div class="flex justify-end gap-2 pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700">
      <Button variant="outline" onclick={() => (voidModalOpen = false)}>Batal</Button>
      <Button variant="danger" loading={voidLoading} onclick={handleVoid}>
        Konfirmasi Pembatalan & Kembalikan Stok
      </Button>
    </div>
  </div>
</Modal>
