<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { cart, cartTotals } from '$lib/stores/cart';
  import { formatRupiah } from '$lib/utils/formatters';
  import type { Product, Category, Customer, Transaction } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import CheckoutModal from '$lib/components/pos/CheckoutModal.svelte';
  import ThermalReceipt from '$lib/components/pos/ThermalReceipt.svelte';
  import MemberStandeeModal from '$lib/components/pos/MemberStandeeModal.svelte';
  import BarcodeScannerModal from '$lib/components/pos/BarcodeScannerModal.svelte';

  let products = $state<Product[]>([]);
  let categories = $state<Category[]>([]);
  let customers = $state<Customer[]>([]);
  let loading = $state(true);

  let standeeOpen = $state(false);
  let itemScannerOpen = $state(false);

  // Filters
  let search = $state('');
  let selectedCategory = $state<number | null>(null);

  // Mobile Active View ('catalog' | 'cart')
  let mobileTab = $state<'catalog' | 'cart'>('catalog');

  // Modals
  let checkoutModalOpen = $state(false);
  let receiptModalOpen = $state(false);
  let completedTransaction = $state<Transaction | null>(null);

  // Item discount modal / inline
  let itemDiscountModalOpen = $state(false);
  let discountItemProductId = $state<number | null>(null);
  let discountItemType = $state<'none' | 'percentage' | 'fixed'>('none');
  let discountItemValue = $state<number>(0);

  async function loadData() {
    loading = true;
    try {
      const [prodRes, catRes, custRes] = await Promise.all([
        api.get<Product[]>('/products?limit=100'),
        api.get<Category[]>('/categories'),
        api.get<Customer[]>('/customers')
      ]);
      products = prodRes.data || [];
      categories = catRes.data || [];
      customers = custRes.data || [];
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat katalog kasir');
    } finally {
      loading = false;
    }
  }

  // Filter products by search and category
  const filteredProducts = $derived(
    products.filter((p) => {
      const matchCat = !selectedCategory || p.category_id === selectedCategory;
      const term = search.toLowerCase().trim();
      const matchSearch =
        !term ||
        p.name.toLowerCase().includes(term) ||
        p.barcode.toLowerCase().includes(term) ||
        p.sku.toLowerCase().includes(term);
      return matchCat && matchSearch;
    })
  );

  function handleScannedProduct(barcode: string) {
    const term = barcode.trim();
    if (!term) return;

    // Search by exact barcode, numeric barcode, or SKU
    const found = products.find(
      (p) =>
        p.barcode === term ||
        p.barcode.replace(/[^0-9A-Za-z]/g, '') === term.replace(/[^0-9A-Za-z]/g, '') ||
        p.sku.toLowerCase() === term.toLowerCase()
    );

    if (found) {
      if (found.stock > 0) {
        cart.addItem(found);
        toast.success(`+ ${found.name} (${formatRupiah(found.sell_price)})`);
      } else {
        toast.error(`Stok "${found.name}" habis!`);
      }
    } else {
      toast.error(`Barang barcode "${term}" tidak ditemukan di katalog toko.`);
    }
  }

  function handleBarcodeSearch(e: KeyboardEvent) {
    if (e.key === 'Enter' && search.trim()) {
      handleScannedProduct(search);
      search = '';
    }
  }

  function handleAddToCart(product: Product) {
    if (product.stock <= 0) {
      toast.error('Stok produk habis!');
      return;
    }
    cart.addItem(product);
  }

  function openItemDiscount(productId: number, curType: string, curVal: number) {
    discountItemProductId = productId;
    discountItemType = (curType as any) || 'none';
    discountItemValue = curVal || 0;
    itemDiscountModalOpen = true;
  }

  function saveItemDiscount() {
    if (discountItemProductId) {
      cart.setItemDiscount(discountItemProductId, discountItemType, discountItemValue);
      itemDiscountModalOpen = false;
      toast.success('Diskon item disimpan');
    }
  }

  function handleTransactionSuccess(transaction: Transaction) {
    completedTransaction = transaction;
    checkoutModalOpen = false;
    receiptModalOpen = true;
    // Reload products to update stock numbers
    loadData();
  }

  onMount(loadData);
</script>

<div class="h-[calc(100vh-6.5rem)] flex flex-col">
  <!-- Mobile Tab Switcher -->
  <div class="lg:hidden flex border-b-3 border-black mb-3">
    <button
      onclick={() => (mobileTab = 'catalog')}
      class="flex-1 py-2.5 font-black text-xs uppercase text-center border-r-2 border-black {mobileTab === 'catalog' ? 'bg-[#FFE600] text-black' : 'bg-white dark:bg-[#222]'}"
    >
      Katalog Produk
    </button>
    <button
      onclick={() => (mobileTab = 'cart')}
      class="flex-1 py-2.5 font-black text-xs uppercase text-center {mobileTab === 'cart' ? 'bg-[#00F0FF] text-black' : 'bg-white dark:bg-[#222]'}"
    >
      Keranjang ({$cartTotals.itemCount}) — {formatRupiah($cartTotals.total)}
    </button>
  </div>

  <!-- Main 2-Column POS Layout -->
  <div class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-4 min-h-0">
    
    <!-- LEFT COLUMN: Product Catalog (7 cols) -->
    <div class="lg:col-span-7 flex flex-col min-h-0 {mobileTab === 'cart' ? 'hidden lg:flex' : 'flex'}">
      <!-- Search and Scan Bar -->
      <div class="mb-3 flex flex-wrap sm:flex-nowrap gap-2">
        <div class="relative flex-1 min-w-[200px]">
          <Input
            placeholder="Scan barcode / ketik nama produk..."
            bind:value={search}
            onkeydown={handleBarcodeSearch}
            class="text-sm font-bold pl-9"
          />
          <svg class="w-5 h-5 absolute left-2.5 top-3 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        {#if search}
          <Button variant="outline" size="sm" onclick={() => (search = '')}>Reset</Button>
        {/if}
        <button
          type="button"
          onclick={() => (itemScannerOpen = true)}
          class="neo-btn bg-[#00F0FF] text-black px-3.5 py-2 text-xs font-black flex items-center gap-1.5 whitespace-nowrap shadow-[3px_3px_0px_0px_#000000]"
          title="Buka Kamera Barcode Scanner untuk memasukkan barang belanjaan"
        >
          <span>📷</span> <span>SCAN BARANG</span>
        </button>
        <Button variant="accent" size="sm" onclick={() => (standeeOpen = true)}>
          📱 QR Member
        </Button>
      </div>

      <!-- Category Filter Chips -->
      <div class="flex gap-2 overflow-x-auto pb-2 mb-3 no-scrollbar flex-shrink-0">
        <button
          onclick={() => (selectedCategory = null)}
          class="neo-btn text-xs px-3 py-1.5 whitespace-nowrap {selectedCategory === null ? 'bg-black text-white dark:bg-white dark:text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
        >
          Semua ({products.length})
        </button>
        {#each categories as cat}
          <button
            onclick={() => (selectedCategory = cat.id)}
            class="neo-btn text-xs px-3 py-1.5 whitespace-nowrap {selectedCategory === cat.id ? 'bg-[#FFE600] text-black' : 'bg-white dark:bg-[#222] text-black dark:text-white'}"
          >
            {cat.name}
          </button>
        {/each}
      </div>

      <!-- Products Grid -->
      <div class="flex-1 overflow-y-auto pr-1">
        {#if loading}
          <div class="p-12 text-center font-bold text-neutral-500">Memuat produk kasir...</div>
        {:else if filteredProducts.length === 0}
          <div class="p-12 text-center font-bold text-neutral-500 bg-white dark:bg-[#1a1a1a] neo-box">
            Barang tidak ditemukan. Periksa kata kunci atau barcode.
          </div>
        {:else}
          <div class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3 pb-4">
            {#each filteredProducts as p (p.id)}
              <button
                type="button"
                onclick={() => handleAddToCart(p)}
                disabled={p.stock <= 0}
                class="neo-box p-3 flex flex-col justify-between text-left transition-all hover:-translate-y-1 active:translate-y-0.5 cursor-pointer bg-white dark:bg-[#1f1f1f] group {p.stock <= 0 ? 'opacity-50 cursor-not-allowed' : ''}"
              >
                <div>
                  <!-- Stock & Category Tag -->
                  <div class="flex items-center justify-between gap-1 mb-2">
                    <span class="text-[10px] font-black uppercase text-neutral-500 truncate max-w-[80px]">
                      {p.category?.name || 'Umum'}
                    </span>
                    <Badge variant={p.stock <= 5 ? 'red' : p.stock <= 15 ? 'yellow' : 'green'} class="text-[9px] px-1 py-0">
                      {p.stock} {p.unit}
                    </Badge>
                  </div>

                  <!-- Name -->
                  <h4 class="font-black text-xs sm:text-sm line-clamp-2 leading-snug group-hover:text-amber-600 dark:group-hover:text-[#FFE600] mb-1">
                    {p.name}
                  </h4>
                  <div class="text-[10px] font-mono text-neutral-400 truncate">{p.barcode || p.sku}</div>
                </div>

                <!-- Price -->
                <div class="mt-3 pt-2 border-t-2 border-dashed border-neutral-200 dark:border-neutral-700 flex justify-between items-center">
                  <span class="font-mono font-black text-xs sm:text-sm text-black dark:text-white">
                    {formatRupiah(p.sell_price)}
                  </span>
                  <span class="text-[10px] font-black bg-[#FFE600] text-black px-1.5 py-0.5 border border-black shadow-[1px_1px_0px_0px_#000000]">
                    + PILIH
                  </span>
                </div>
              </button>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <!-- RIGHT COLUMN: Cart Panel (5 cols) -->
    <div class="lg:col-span-5 neo-box flex flex-col min-h-0 bg-white dark:bg-[#1a1a1a] {mobileTab === 'catalog' ? 'hidden lg:flex' : 'flex'}">
      <!-- Cart Header -->
      <div class="p-3.5 border-b-3 border-black dark:border-white bg-[#FFE600] text-black flex items-center justify-between">
        <div class="flex items-center gap-2">
          <h3 class="font-black text-sm uppercase tracking-tight">Keranjang Belanja</h3>
          <span class="bg-black text-white text-xs font-black px-2 py-0.5 border border-black">
            {$cartTotals.itemCount} Item
          </span>
        </div>
        {#if $cart.items.length > 0}
          <button
            onclick={() => cart.clearCart()}
            class="text-xs font-black underline text-red-700 hover:text-black"
          >
            Kosongkan
          </button>
        {/if}
      </div>

      <!-- Cart Items List -->
      <div class="flex-1 overflow-y-auto p-3 divide-y-2 divide-neutral-200 dark:divide-neutral-700">
        {#if $cart.items.length === 0}
          <div class="h-full flex flex-col items-center justify-center text-center p-6 text-neutral-400">
            <svg class="w-12 h-12 mb-2 text-neutral-300 dark:text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <p class="font-black text-sm text-neutral-500">Keranjang masih kosong</p>
            <p class="text-xs">Klik produk di katalog untuk menambahkan ke transaksi</p>
          </div>
        {:else}
          {#each $cart.items as item (item.product.id)}
            {@const grossPrice = item.product.sell_price * item.qty}
            {@const itemDisc = item.discount_type === 'percentage'
              ? grossPrice * (item.discount_value / 100)
              : item.discount_type === 'fixed'
              ? item.discount_value
              : 0}
            {@const itemNet = Math.max(0, grossPrice - itemDisc)}

            <div class="py-2.5 flex flex-col gap-1.5">
              <div class="flex items-start justify-between gap-2">
                <div class="flex-1">
                  <h5 class="text-xs font-black leading-tight text-black dark:text-white">
                    {item.product.name}
                  </h5>
                  <div class="text-[11px] font-mono text-neutral-500">
                    {formatRupiah(item.product.sell_price)} / {item.product.unit}
                  </div>
                </div>
                <!-- Subtotal item -->
                <div class="text-right">
                  <span class="font-mono font-black text-sm text-black dark:text-white">
                    {formatRupiah(itemNet)}
                  </span>
                  {#if itemDisc > 0}
                    <div class="text-[10px] text-red-500 font-mono">
                      -{formatRupiah(itemDisc)}
                    </div>
                  {/if}
                </div>
              </div>

              <!-- Item Controls: Qty + Discount + Delete -->
              <div class="flex items-center justify-between pt-1">
                <div class="flex items-center gap-1 border-2 border-black dark:border-white bg-white dark:bg-[#222]">
                  <button
                    onclick={() => cart.updateQty(item.product.id, item.qty - 1)}
                    class="w-6 h-6 flex items-center justify-center font-black hover:bg-neutral-200 dark:hover:bg-neutral-700"
                  >
                    -
                  </button>
                  <span class="w-8 text-center text-xs font-mono font-black">{item.qty}</span>
                  <button
                    onclick={() => cart.updateQty(item.product.id, item.qty + 1)}
                    class="w-6 h-6 flex items-center justify-center font-black hover:bg-neutral-200 dark:hover:bg-neutral-700"
                  >
                    +
                  </button>
                </div>

                <div class="flex items-center gap-2">
                  <button
                    onclick={() => openItemDiscount(item.product.id, item.discount_type, item.discount_value)}
                    class="text-[11px] font-black underline {item.discount_value > 0 ? 'text-[#00E676] dark:text-[#33eb91]' : 'text-neutral-500 hover:text-black dark:hover:text-white'}"
                  >
                    {item.discount_value > 0 ? `Diskon (${item.discount_value}${item.discount_type === 'percentage' ? '%' : 'rb'})` : 'Beri Diskon'}
                  </button>
                  <button
                    onclick={() => cart.removeItem(item.product.id)}
                    class="text-[11px] font-black text-red-500 hover:text-red-700"
                  >
                    Hapus
                  </button>
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>

      <!-- Cart Bottom Calculations & Checkout Button -->
      {#if $cart.items.length > 0}
        <div class="p-4 border-t-3 border-black dark:border-white bg-neutral-50 dark:bg-[#151515] space-y-2.5">
          <!-- Discount & Tax switches -->
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-neutral-600 dark:text-neutral-400">PPN (11%):</span>
            <button
              type="button"
              onclick={() => cart.toggleTax()}
              class="font-black px-2 py-0.5 border border-black text-[11px] {$cart.taxEnabled ? 'bg-[#00E676] text-black' : 'bg-neutral-200 text-neutral-500'}"
            >
              {$cart.taxEnabled ? 'AKTIF' : 'NON-AKTIF'}
            </button>
          </div>

          <!-- Calculation Lines -->
          <div class="space-y-1 text-xs font-bold border-t border-dashed border-neutral-300 dark:border-neutral-700 pt-2">
            <div class="flex justify-between text-neutral-600 dark:text-neutral-400">
              <span>Subtotal:</span>
              <span class="font-mono">{formatRupiah($cartTotals.subtotal)}</span>
            </div>
            {#if $cartTotals.taxAmount > 0}
              <div class="flex justify-between text-neutral-600 dark:text-neutral-400">
                <span>PPN 11%:</span>
                <span class="font-mono">{formatRupiah($cartTotals.taxAmount)}</span>
              </div>
            {/if}
          </div>

          <!-- Grand Total Bar -->
          <div class="p-3 bg-[#FFE600] text-black border-2 border-black flex justify-between items-center shadow-[3px_3px_0px_0px_#000000]">
            <span class="text-xs font-black uppercase">Total Bayar:</span>
            <span class="text-xl font-black font-mono">
              {formatRupiah($cartTotals.total)}
            </span>
          </div>

          <!-- Checkout Giant Button -->
          <Button
            variant="primary"
            size="lg"
            class="w-full text-base font-black py-4 bg-[#00E676] hover:bg-[#2ae48c] text-black"
            onclick={() => (checkoutModalOpen = true)}
          >
            CHECKOUT & BAYAR (F9)
          </Button>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Item Discount Modal -->
<Modal
  open={itemDiscountModalOpen}
  title="DISKON PRODUK"
  maxWidth="max-w-sm"
  onclose={() => (itemDiscountModalOpen = false)}
>
  <div class="space-y-3">
    <div>
      <label for="disc-type" class="text-xs font-black uppercase tracking-wider block mb-1">Tipe Diskon</label>
      <select id="disc-type" bind:value={discountItemType} class="neo-input font-bold text-sm bg-white dark:bg-[#222]">
        <option value="none">Tanpa Diskon</option>
        <option value="percentage">Persentase (%)</option>
        <option value="fixed">Nominal (Rp)</option>
      </select>
    </div>

    {#if discountItemType !== 'none'}
      <Input
        label={discountItemType === 'percentage' ? 'Persen Diskon (%)' : 'Nominal Diskon (Rp)'}
        type="number"
        bind:value={discountItemValue}
      />
    {/if}

    <div class="flex justify-end gap-2 pt-2">
      <Button variant="outline" onclick={() => (itemDiscountModalOpen = false)}>Batal</Button>
      <Button variant="primary" onclick={saveItemDiscount}>Terapkan</Button>
    </div>
  </div>
</Modal>

<!-- Checkout Modal with QRIS Midtrans & Cash Support -->
<CheckoutModal
  open={checkoutModalOpen}
  {customers}
  onclose={() => (checkoutModalOpen = false)}
  onsuccess={handleTransactionSuccess}
/>

<!-- Thermal Receipt Modal -->
<Modal
  open={receiptModalOpen}
  title="STRUK PEMBAYARAN KASIR"
  maxWidth="max-w-md"
  onclose={() => (receiptModalOpen = false)}
>
  {#if completedTransaction}
    <ThermalReceipt
      transaction={completedTransaction}
      onclose={() => (receiptModalOpen = false)}
    />
  {/if}
</Modal>

<!-- Member Standee Modal -->
<MemberStandeeModal
  open={standeeOpen}
  onclose={() => (standeeOpen = false)}
/>

<!-- Continuous Item Barcode Scanner Modal for Cart -->
<BarcodeScannerModal
  open={itemScannerOpen}
  title="SCAN BARANG KE KERANJANG (KAMERA)"
  continuous={true}
  onclose={() => (itemScannerOpen = false)}
  onscan={handleScannedProduct}
/>
