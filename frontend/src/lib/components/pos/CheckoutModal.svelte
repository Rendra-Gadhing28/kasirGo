<script lang="ts">
  import { cart, cartTotals } from '$lib/stores/cart';
  import { toast } from '$lib/stores/toast';
  import { api } from '$lib/api/client';
  import { formatRupiah } from '$lib/utils/formatters';
  import type { Customer, Transaction } from '$lib/types';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

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

      <!-- Customer Selection -->
      <div>
        <label for="checkout-cust" class="text-xs font-black uppercase tracking-wider block mb-1.5">
          Pilih Pelanggan / Member (Opsional)
        </label>
        <select
          id="checkout-cust"
          bind:value={selectedCustomerId}
          class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
        >
          <option value="">Pelanggan Umum (Walk-in / Tanpa Member)</option>
          {#each customers as cust}
            <option value={cust.id}>
              {cust.name} ({cust.phone || '-'}) — Poin: {cust.points}
            </option>
          {/each}
        </select>
        {#if selectedCustomerId}
          <div class="mt-1 text-xs text-[#00E676] font-black">
            + Pelanggan akan mendapatkan {estimatedPoints} poin loyalitas dari transaksi ini.
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
            class="neo-btn py-3 text-xs font-black {paymentMethod === 'cash' ? 'bg-[#00F0FF] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            💵 TUNAI (CASH)
          </button>
          <button
            type="button"
            onclick={() => (paymentMethod = 'qris')}
            class="neo-btn py-3 text-xs font-black {paymentMethod === 'qris' ? 'bg-[#FFE600] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            📱 QRIS MIDTRANS
          </button>
          <button
            type="button"
            onclick={() => (paymentMethod = 'transfer')}
            class="neo-btn py-3 text-xs font-black {paymentMethod === 'transfer' ? 'bg-[#B388FF] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            🏦 TRANSFER BANK
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
