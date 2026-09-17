<script lang="ts">
  import { cart, cartTotals } from '$lib/stores/cart';
  import { toast } from '$lib/stores/toast';
  import { api } from '$lib/api/client';
  import { formatRupiah } from '$lib/utils/formatters';
  import type { Customer, Transaction } from '$lib/types';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import BarcodeScannerModal from '$lib/components/pos/BarcodeScannerModal.svelte';

  interface Props {
    open: boolean;
    customers: Customer[];
    onclose: () => void;
    onsuccess: (transaction: Transaction) => void;
  }

  let { open, customers, onclose, onsuccess }: Props = $props();

  let paymentMethod = $state<'cash' | 'qris' | 'transfer'>('cash');
  let cashPaid = $state<number>(0);
  let selectedCustomerId = $state<number | ''>('');
  let checkoutLoading = $state(false);

  // Member search & scan
  let memberSearchQuery = $state('');
  let memberScannerOpen = $state(false);
  let searchingMember = $state(false);

  const selectedCustomer = $derived(
    customers.find((c) => c.id === selectedCustomerId) || null
  );

  async function handleSearchMember(queryOverride?: string) {
    const raw = (queryOverride !== undefined ? queryOverride : memberSearchQuery).trim();
    if (!raw) return;
    const q = raw.toLowerCase();

    // 1. Search local list
    const foundLocal = customers.find(
      (c) =>
        (c.member_code && c.member_code.toLowerCase() === q) ||
        c.phone.toLowerCase() === q ||
        c.phone.replace(/[^0-9]/g, '') === q.replace(/[^0-9]/g, '') ||
        c.name.toLowerCase().includes(q)
    );

    if (foundLocal) {
      selectedCustomerId = foundLocal.id;
      memberSearchQuery = '';
      toast.success(`Member terhubung: ${foundLocal.name} (${foundLocal.member_code || foundLocal.phone})`);
      return;
    }

    // 2. Lookup from API
    searchingMember = true;
    try {
      const res = await api.get<Customer>(`/customers/lookup/${encodeURIComponent(raw)}`);
      if (res.data) {
        const cust = res.data;
        if (!customers.some((c) => c.id === cust.id)) {
          customers.push(cust);
        }
        selectedCustomerId = cust.id;
        memberSearchQuery = '';
        toast.success(`Member terhubung: ${cust.name} (${cust.member_code || cust.phone})`);
        return;
      }
    } catch {
      // not found in backend
    } finally {
      searchingMember = false;
    }

    toast.error(`Member dengan no kartu / HP "${raw}" tidak ditemukan.`);
  }

  function clearMember() {
    selectedCustomerId = '';
    memberSearchQuery = '';
  }

  // QRIS Simulation state
  let qrisStep = $state<'prepare' | 'waiting' | 'success'>('prepare');
  let currentTransaction = $state<Transaction | null>(null);
  let qrSimulating = $state(false);

  $effect(() => {
    if (open) {
      cashPaid = $cartTotals.total;
      qrisStep = 'prepare';
      currentTransaction = null;
    }
  });

  const changeAmount = $derived(
    paymentMethod === 'cash' ? Math.max(0, cashPaid - $cartTotals.total) : 0
  );

  const canSubmitCash = $derived(
    paymentMethod !== 'cash' || cashPaid >= $cartTotals.total
  );

  const estimatedPoints = $derived(Math.floor($cartTotals.total / 10000));

  function setQuickCash(amount: number) {
    cashPaid = amount;
  }

  async function handleCheckout() {
    if ($cart.items.length === 0) {
      toast.error('Keranjang belanja kosong');
      return;
    }

    if (paymentMethod === 'cash' && cashPaid < $cartTotals.total) {
      toast.error('Uang tunai yang diterima kurang dari total tagihan');
      return;
    }

    checkoutLoading = true;
    try {
      const payload = {
        customer_id: selectedCustomerId ? Number(selectedCustomerId) : undefined,
        items: $cart.items.map((i) => ({
          product_id: i.product.id,
          qty: i.qty,
          discount_type: i.discount_type,
          discount_value: i.discount_value
        })),
        discount_type: $cart.discountType,
        discount_value: $cart.discountValue,
        tax_enabled: $cart.taxEnabled,
        payment_method: paymentMethod,
        paid_amount: paymentMethod === 'cash' ? cashPaid : $cartTotals.total,
        notes: $cart.notes
      };

      const res = await api.post<Transaction>('/pos/checkout', payload);
      const trans = res.data;
      if (!trans) throw new Error('Data transaksi gagal diproses');

      if (paymentMethod === 'qris') {
        currentTransaction = trans;
        qrisStep = 'waiting';
        toast.info('QRIS siap, silakan scan menggunakan QRIS Simulator Midtrans');
      } else {
        toast.success('Transaksi sukses!');
        cart.clearCart();
        onsuccess(trans);
      }
    } catch (e: any) {
      toast.error(e.message || 'Gagal memproses checkout');
    } finally {
      checkoutLoading = false;
    }
  }

  async function handleSimulateQRISPayment() {
    if (!currentTransaction) return;
    qrSimulating = true;
    try {
      const res = await api.post<Transaction>(`/pos/qris/confirm/${currentTransaction.invoice_no}`);
      if (res.data) {
        toast.success('Pembayaran QRIS Berhasil Diverifikasi!');
        cart.clearCart();
        onsuccess(res.data);
      }
    } catch (e: any) {
      toast.error(e.message || 'Gagal verifikasi pembayaran QRIS');
    } finally {
      qrSimulating = false;
    }
  }
</script>

<Modal
  {open}
  title="CHECKOUT & PEMBAYARAN"
  maxWidth="max-w-xl"
  {onclose}
>
  {#if qrisStep === 'waiting' && currentTransaction}
    <!-- QRIS Waiting & Simulator Screen -->
    <div class="space-y-4 text-center py-2">
      <div class="bg-[#FFE600] p-3 border-2 border-black inline-block font-black text-xs uppercase tracking-wider">
        Scan QRIS Menggunakan Midtrans Simulator
      </div>

      <!-- QR Box Display -->
      <div class="my-4 p-5 bg-white border-3 border-black shadow-[6px_6px_0px_0px_#000000] inline-block max-w-[280px]">
        <!-- High Quality QR Code Generator using QuickChart API / Midtrans standard -->
        <img
          src="https://api.qrserver.com/v1/create-qr-code/?size=240x240&data={encodeURIComponent('00020101021226680016ID.CO.MIDTRANS01189360091800000000000215' + currentTransaction.invoice_no + '53033605408' + currentTransaction.total_amount + '5802ID5911KASIRPRO6007JAKARTA6304ABCD')}"
          alt="QRIS Midtrans"
          class="w-56 h-56 mx-auto border-2 border-black"
        />
        <div class="mt-3 font-mono font-black text-sm">
          {formatRupiah(currentTransaction.total_amount)}
        </div>
        <p class="text-[10px] text-neutral-500 font-bold uppercase mt-1">Invoice: {currentTransaction.invoice_no}</p>
      </div>

      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-300 max-w-sm mx-auto">
        Untuk testing Midtrans QRIS Simulator, Anda dapat menekan tombol konfirmasi pembayaran di bawah.
      </p>

      <div class="flex justify-center gap-3 pt-2">
        <Button variant="outline" onclick={onclose}>Batalkan</Button>
        <Button
          variant="success"
          size="lg"
          loading={qrSimulating}
          onclick={handleSimulateQRISPayment}
        >
          SIMULASI PEMBAYARAN SUKSES
        </Button>
      </div>
    </div>
  {:else}
    <!-- Normal Checkout Flow -->
    <div class="space-y-5">
      <!-- Total Summary Bar -->
      <div class="p-4 bg-[#FFE600] text-black border-3 border-black shadow-[4px_4px_0px_0px_#000000] flex justify-between items-center">
        <div>
          <span class="text-xs font-black uppercase tracking-wider text-neutral-800">Total Tagihan</span>
          <div class="text-2xl sm:text-3xl font-black font-mono">
            {formatRupiah($cartTotals.total)}
          </div>
        </div>
        <div class="text-right text-xs font-bold">
          <div>{$cartTotals.itemCount} Item Barang</div>
          <div>PPN 11%: {formatRupiah($cartTotals.taxAmount)}</div>
        </div>
      </div>

      <!-- Customer / Member Card Selection -->
      <div class="p-3 bg-neutral-50 dark:bg-[#202020] border-2 border-black dark:border-white space-y-2">
        <div class="flex items-center justify-between">
          <label for="member-query" class="text-xs font-black uppercase tracking-wider flex items-center gap-1.5 text-black dark:text-white">
            <svg class="w-4 h-4 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
            <span>Member Toko (Poin Belanja)</span>
          </label>
          {#if selectedCustomer}
            <button
              type="button"
              onclick={clearMember}
              class="text-[11px] font-black underline text-red-600 dark:text-red-400 hover:opacity-80"
            >
              Lepas Member
            </button>
          {/if}
        </div>

        {#if selectedCustomer}
          <!-- Connected Member Banner -->
          <div class="p-2.5 bg-[#00E676]/15 dark:bg-[#00E676]/10 border-2 border-[#00E676] flex items-center justify-between">
            <div>
              <div class="flex items-center gap-2">
                <span class="font-black text-sm text-black dark:text-white">{selectedCustomer.name}</span>
                <span class="bg-black text-white px-1.5 py-0.5 text-[10px] font-mono font-black border border-black">
                  {selectedCustomer.member_code || ('MBR-' + String(selectedCustomer.id).padStart(4, '0'))}
                </span>
              </div>
              <div class="text-xs text-neutral-500 font-mono mt-0.5">{selectedCustomer.phone || '-'}</div>
            </div>
            <div class="text-right">
              <div class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Saldo: {selectedCustomer.points} Poin</div>
              <div class="text-xs font-black text-green-600 dark:text-[#33eb91]">+{estimatedPoints} Poin Belanja</div>
            </div>
          </div>
        {:else}
          <!-- Search / Scan Input Row -->
          <div class="flex gap-1.5">
            <div class="flex-1">
              <input
                id="member-query"
                type="text"
                bind:value={memberSearchQuery}
                onkeydown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault();
                    handleSearchMember();
                  }
                }}
                placeholder="Scan / Ketik No. Kartu Member (MBR-...) atau No. HP..."
                class="neo-input text-xs font-bold py-2 bg-white dark:bg-[#252525] text-black dark:text-white"
              />
            </div>
            <button
              type="button"
              onclick={() => handleSearchMember()}
              disabled={searchingMember}
              class="neo-btn bg-[#FFE600] text-black px-3 text-xs font-black"
              title="Cari Member"
            >
              CARI
            </button>
            <button
              type="button"
              onclick={() => (memberScannerOpen = true)}
              class="neo-btn bg-[#00F0FF] text-black px-3 text-xs font-black flex items-center gap-1.5"
              title="Scan QR Member dari HP Pelanggan"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>SCAN</span>
            </button>
          </div>

          <!-- Dropdown Fallback -->
          <div class="pt-0.5">
            <select
              id="checkout-cust"
              bind:value={selectedCustomerId}
              class="neo-input font-bold text-xs py-1.5 bg-white dark:bg-[#222] text-black dark:text-white"
            >
              <option value="">-- Atau Pilih Manual dari Daftar ({customers.length} Member) --</option>
              {#each customers as cust}
                <option value={cust.id}>
                  {cust.member_code || ('MBR-' + String(cust.id).padStart(4, '0'))} - {cust.name} ({cust.phone || '-'}) — {cust.points} Poin
                </option>
              {/each}
            </select>
          </div>
        {/if}
      </div>

      <!-- Payment Method Tabs -->
      <div>
        <span class="text-xs font-black uppercase tracking-wider block mb-2">Metode Pembayaran</span>
        <div class="grid grid-cols-3 gap-2">
          <button
            type="button"
            onclick={() => (paymentMethod = 'cash')}
            class="neo-btn py-3 text-xs font-black flex items-center justify-center gap-1.5 {paymentMethod === 'cash' ? 'bg-[#00F0FF] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <span>TUNAI (CASH)</span>
          </button>
          <button
            type="button"
            onclick={() => (paymentMethod = 'qris')}
            class="neo-btn py-3 text-xs font-black flex items-center justify-center gap-1.5 {paymentMethod === 'qris' ? 'bg-[#FFE600] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
            <span>QRIS MIDTRANS</span>
          </button>
          <button
            type="button"
            onclick={() => (paymentMethod = 'transfer')}
            class="neo-btn py-3 text-xs font-black flex items-center justify-center gap-1.5 {paymentMethod === 'transfer' ? 'bg-[#B388FF] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z" />
            </svg>
            <span>TRANSFER BANK</span>
          </button>
        </div>
      </div>

      <!-- Cash Calculation Panel -->
      {#if paymentMethod === 'cash'}
        <div class="p-4 bg-neutral-100 dark:bg-[#222] border-2 border-black dark:border-white space-y-3">
          <Input
            label="Nominal Uang Tunai Diterima"
            type="number"
            bind:value={cashPaid}
            required
          />

          <!-- Quick Cash Presets -->
          <div class="flex flex-wrap gap-2 pt-1">
            <button
              type="button"
              onclick={() => setQuickCash($cartTotals.total)}
              class="neo-btn bg-white dark:bg-[#333] text-black dark:text-white px-2.5 py-1 text-xs font-black"
            >
              Uang Pas
            </button>
            <button
              type="button"
              onclick={() => setQuickCash(20000)}
              class="neo-btn bg-white dark:bg-[#333] text-black dark:text-white px-2.5 py-1 text-xs font-black"
            >
              Rp 20.000
            </button>
            <button
              type="button"
              onclick={() => setQuickCash(50000)}
              class="neo-btn bg-white dark:bg-[#333] text-black dark:text-white px-2.5 py-1 text-xs font-black"
            >
              Rp 50.000
            </button>
            <button
              type="button"
              onclick={() => setQuickCash(100000)}
              class="neo-btn bg-white dark:bg-[#333] text-black dark:text-white px-2.5 py-1 text-xs font-black"
            >
              Rp 100.000
            </button>
            <button
              type="button"
              onclick={() => setQuickCash(200000)}
              class="neo-btn bg-white dark:bg-[#333] text-black dark:text-white px-2.5 py-1 text-xs font-black"
            >
              Rp 200.000
            </button>
          </div>

          <!-- Kembalian Display -->
          <div class="pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-600 flex justify-between items-center">
            <span class="text-sm font-black uppercase">Kembalian:</span>
            <span class="text-xl font-black font-mono text-[#00E676]">
              {formatRupiah(changeAmount)}
            </span>
          </div>
        </div>
      {/if}

      <!-- Action Buttons -->
      <div class="flex justify-end gap-3 pt-3 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700">
        <Button variant="outline" onclick={onclose}>Batal</Button>
        <Button
          variant="primary"
          size="lg"
          disabled={!canSubmitCash}
          loading={checkoutLoading}
          onclick={handleCheckout}
        >
          {paymentMethod === 'qris' ? 'GENERATE QRIS SEKARANG' : 'SELESAIKAN PEMBAYARAN'}
        </Button>
      </div>
    </div>
  {/if}
</Modal>

<!-- Member QR Camera Scanner Modal -->
<BarcodeScannerModal
  open={memberScannerOpen}
  onclose={() => (memberScannerOpen = false)}
  onscan={(scannedCode) => {
    memberScannerOpen = false;
    handleSearchMember(scannedCode);
  }}
/>
