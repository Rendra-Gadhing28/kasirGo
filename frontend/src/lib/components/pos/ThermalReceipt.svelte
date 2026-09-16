<script lang="ts">
  import { formatRupiah, formatDateTime } from '$lib/utils/formatters';
  import type { Transaction } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    transaction: Transaction;
    onclose?: () => void;
  }

  let { transaction, onclose }: Props = $props();

  function triggerPrint() {
    window.print();
  }
</script>

<div class="space-y-4">
  <!-- Print and Close Buttons (Hidden on print) -->
  <div class="flex justify-between items-center print:hidden border-b-2 border-black pb-3">
    <span class="text-xs font-black uppercase text-neutral-600 dark:text-neutral-300">Format Struk Thermal 80mm</span>
    <div class="flex gap-2">
      {#if onclose}
        <Button variant="outline" size="sm" onclick={onclose}>Tutup</Button>
      {/if}
      <Button variant="primary" size="sm" onclick={triggerPrint}>
        Cetak Struk
      </Button>
    </div>
  </div>

  <!-- Thermal Paper Container -->
  <div
    id="thermal-receipt"
    class="bg-white text-black font-mono text-xs p-4 mx-auto max-w-[320px] border-2 border-dashed border-neutral-400 shadow-md leading-relaxed"
  >
    <!-- Store Header -->
    <div class="text-center pb-3 border-b border-dashed border-black">
      <h2 class="text-base font-black uppercase tracking-wider">{transaction.outlet?.name || 'KASIRPRO STORE'}</h2>
      <p class="text-[11px] text-neutral-700">{transaction.outlet?.address || 'Pasar Modern Indonesia'}</p>
      <p class="text-[11px] text-neutral-700">Telp: {transaction.outlet?.phone || '081234567890'}</p>
    </div>

    <!-- Meta Info -->
    <div class="py-2 border-b border-dashed border-black text-[11px] space-y-0.5">
      <div class="flex justify-between">
        <span>No:</span>
        <span class="font-bold">{transaction.invoice_no}</span>
      </div>
      <div class="flex justify-between">
        <span>Waktu:</span>
        <span>{formatDateTime(transaction.created_at)}</span>
      </div>
      <div class="flex justify-between">
        <span>Kasir:</span>
        <span>{transaction.user?.name || 'Kasir'}</span>
      </div>
      {#if transaction.customer}
        <div class="flex justify-between">
          <span>Pelanggan:</span>
          <span class="font-bold">{transaction.customer.name}</span>
        </div>
      {/if}
    </div>

    <!-- Items List -->
    <div class="py-2 border-b border-dashed border-black space-y-2">
      {#each transaction.items as item}
        <div>
          <div class="font-bold text-[11px]">{item.product_name}</div>
          <div class="flex justify-between text-[11px] pl-2 text-neutral-800">
            <span>{item.qty} x {formatRupiah(item.price)}</span>
            <span class="font-bold">{formatRupiah(item.subtotal)}</span>
          </div>
          {#if item.discount_amount > 0}
            <div class="text-[10px] text-red-600 pl-2">
              (Diskon: -{formatRupiah(item.discount_amount)})
            </div>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Calculations -->
    <div class="py-2 border-b border-dashed border-black text-[11px] space-y-1">
      <div class="flex justify-between">
        <span>Subtotal:</span>
        <span>{formatRupiah(transaction.subtotal)}</span>
      </div>

      {#if transaction.discount_amount > 0}
        <div class="flex justify-between text-red-600">
          <span>Diskon Total:</span>
          <span>-{formatRupiah(transaction.discount_amount)}</span>
        </div>
      {/if}

      {#if transaction.tax_amount > 0}
        <div class="flex justify-between">
          <span>PPN ({transaction.tax_rate}%):</span>
          <span>{formatRupiah(transaction.tax_amount)}</span>
        </div>
      {/if}

      <div class="flex justify-between text-sm font-black pt-1 border-t border-dotted border-black">
        <span>TOTAL:</span>
        <span>{formatRupiah(transaction.total_amount)}</span>
      </div>
    </div>

    <!-- Payment Detail -->
    <div class="py-2 border-b border-dashed border-black text-[11px] space-y-1">
      <div class="flex justify-between">
        <span>Metode:</span>
        <span class="uppercase font-bold">{transaction.payment_method}</span>
      </div>
      {#if transaction.payment_method === 'cash'}
        <div class="flex justify-between">
          <span>Bayar (Tunai):</span>
          <span>{formatRupiah(transaction.paid_amount)}</span>
        </div>
        <div class="flex justify-between font-bold">
          <span>Kembalian:</span>
          <span>{formatRupiah(transaction.change_amount)}</span>
        </div>
      {:else}
        <div class="flex justify-between">
          <span>Status Bayar:</span>
          <span class="font-bold uppercase text-green-700">{transaction.payment_status}</span>
        </div>
      {/if}
    </div>

    <!-- Footer Greeting -->
    <div class="text-center pt-3 text-[10px] text-neutral-600 space-y-1">
      <p class="font-bold uppercase">Terima Kasih Atas Kunjungan Anda!</p>
      <p>Barang yang sudah dibeli tidak dapat ditukar kecuali ada perjanjian.</p>
      <p class="font-bold pt-1 text-[9px] text-neutral-400">Powered by KasirPro</p>
    </div>
  </div>
</div>
