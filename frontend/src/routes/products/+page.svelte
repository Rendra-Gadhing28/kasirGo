<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah } from '$lib/utils/formatters';
  import type { Product, Category } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import BarcodeScannerModal from '$lib/components/pos/BarcodeScannerModal.svelte';
  import { lookupBarcodeInfo } from '$lib/utils/barcodeLookup';

  let products = $state<Product[]>([]);
  let categories = $state<Category[]>([]);
  let loading = $state(true);

  // Search & Filter
  let search = $state('');
  let selectedCategory = $state<number | ''>('');
  let page = $state(1);
  let limit = 15;
  let totalPages = $state(1);
  let totalCount = $state(0);

  // Modal Create/Edit
  let modalOpen = $state(false);
  let scannerOpen = $state(false);
  let lookupLoading = $state(false);
  let editId = $state<number | null>(null);
  let formCategoryId = $state<number | ''>('');
  let formName = $state('');
  let formSKU = $state('');
  let formBarcode = $state('');
  let formBuyPrice = $state<number>(0);
  let formSellPrice = $state<number>(0);
  let formStock = $state<number>(0);
  let formUnit = $state('pcs');
  let formImageUrl = $state('');
  let formBrand = $state('');
  let formLoading = $state(false);

  // Delete modal
  let deleteModalOpen = $state(false);
  let deleteId = $state<number | null>(null);

  const units = ['pcs', 'botol', 'kotak', 'bungkus', 'kg', 'gram', 'liter', 'karung', 'pack'];

  async function loadCategories() {
    try {
      const res = await api.get<Category[]>('/categories');
      categories = res.data || [];
    } catch (e) {
      // ignore
    }
  }

  async function loadProducts() {
    loading = true;
    try {
      let query = `/products?page=${page}&limit=${limit}`;
      if (search) query += `&search=${encodeURIComponent(search)}`;
      if (selectedCategory) query += `&category_id=${selectedCategory}`;

      const res = await api.get<Product[]>(query);
      products = res.data || [];
      if (res.meta) {
        totalPages = res.meta.total_pages || 1;
        totalCount = res.meta.total || 0;
      }
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat produk');
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    editId = null;
    formCategoryId = categories.length > 0 ? categories[0].id : '';
    formName = '';
    formSKU = '';
    formBarcode = '';
    formBuyPrice = 0;
    formSellPrice = 0;
    formStock = 10;
    formUnit = 'pcs';
    formImageUrl = '';
    formBrand = '';
    modalOpen = true;
  }

  function openEditModal(prd: Product) {
    editId = prd.id;
    formCategoryId = prd.category_id;
    formName = prd.name;
    formSKU = prd.sku;
    formBarcode = prd.barcode;
    formBuyPrice = prd.buy_price;
    formSellPrice = prd.sell_price;
    formStock = prd.stock;
    formUnit = prd.unit;
    formImageUrl = prd.image_url || '';
    formBrand = '';
    modalOpen = true;
  }

  async function handleBarcodeDetected(barcode: string) {
    const cleanCode = barcode.trim();
    if (!cleanCode) return;
    formBarcode = cleanCode;
    if (!formSKU) {
      formSKU = 'SKU-' + cleanCode.slice(-6);
    }
    lookupLoading = true;
    toast.info('Mencari database produk untuk barcode: ' + cleanCode + '...');
    try {
      const res = await lookupBarcodeInfo(cleanCode);
      if (res && res.found && res.name) {
        formName = res.name;
        if (res.sku) formSKU = res.sku;
        if (res.unit) formUnit = res.unit;
        if (res.imageUrl) formImageUrl = res.imageUrl;
        if (res.brand) formBrand = res.brand;

        // Auto select matched category
        if (res.categoryHint && categories.length > 0) {
          const matched = categories.find((c) =>
            c.name.toLowerCase().includes(res.categoryHint!.toLowerCase()) ||
            res.categoryHint!.toLowerCase().includes(c.name.toLowerCase())
          );
          if (matched) {
            formCategoryId = matched.id;
          }
        }

        toast.success(`Ditemukan: ${res.name} ${res.brand ? '(' + res.brand + ')' : ''}`);
      } else {
        toast.info(`Barcode ${cleanCode} belum terdaftar di Open Food Facts. Silakan isi nama produk secara manual.`);
      }
    } catch (e) {
      toast.warn('Pencarian katalog selesai. Silakan lengkapi data produk.');
    } finally {
      lookupLoading = false;
    }
  }

  function generateRandomBarcode() {
    formBarcode = '899' + Math.floor(1000000000 + Math.random() * 9000000000).toString();
    if (!formSKU) {
      formSKU = 'SKU-' + Math.random().toString(36).substring(2, 8).toUpperCase();
    }
  }

  async function handleSubmit() {
    if (!formName.trim()) {
      toast.error('Nama produk wajib diisi');
      return;
    }
    if (!formCategoryId) {
      toast.error('Kategori produk wajib dipilih');
      return;
    }
    if (formSellPrice <= 0) {
      toast.error('Harga jual harus lebih dari Rp 0');
      return;
    }

    formLoading = true;
    try {
      const payload = {
        category_id: Number(formCategoryId),
        name: formName,
        sku: formSKU,
        barcode: formBarcode,
        buy_price: Number(formBuyPrice),
        sell_price: Number(formSellPrice),
        stock: Number(formStock),
        unit: formUnit,
        image_url: formImageUrl
      };

      if (editId) {
        await api.put(`/products/${editId}`, payload);
        toast.success('Produk berhasil diperbarui');
      } else {
        await api.post('/products', payload);
        toast.success('Produk berhasil ditambahkan');
      }
      modalOpen = false;
      await loadProducts();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menyimpan produk');
    } finally {
      formLoading = false;
    }
  }

  function confirmDelete(id: number) {
    deleteId = id;
    deleteModalOpen = true;
  }

  async function handleDelete() {
    if (!deleteId) return;
    try {
      await api.delete(`/products/${deleteId}`);
      toast.success('Produk berhasil dihapus');
      deleteModalOpen = false;
      await loadProducts();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menghapus produk');
    }
  }

  onMount(async () => {
    await loadCategories();
    await loadProducts();
  });
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Katalog Produk</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Total {totalCount} barang terdaftar di inventaris</p>
    </div>
    {#if $auth.user?.role !== 'kasir'}
      <Button variant="primary" onclick={openCreateModal}>
        + TAMBAH PRODUK BARU
      </Button>
    {/if}
  </div>

  <!-- Filter & Search Bar -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div class="sm:col-span-2">
      <Input
        placeholder="Cari nama barang, barcode (899...), atau SKU..."
        bind:value={search}
        onkeydown={(e) => { if (e.key === 'Enter') { page = 1; loadProducts(); } }}
      />
    </div>
    <div class="flex gap-2">
      <select
        bind:value={selectedCategory}
        onchange={() => { page = 1; loadProducts(); }}
        class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
      >
        <option value="">Semua Kategori</option>
        {#each categories as cat}
          <option value={cat.id}>{cat.name}</option>
        {/each}
      </select>
      <Button variant="accent" onclick={() => { page = 1; loadProducts(); }}>
        Cari
      </Button>
    </div>
  </div>

  <!-- Products Table -->
  <div class="neo-box overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-[#FFE600] text-black border-b-3 border-black dark:border-white font-black text-xs uppercase tracking-wider">
            <th class="p-3">Produk</th>
            <th class="p-3">Kategori</th>
            <th class="p-3">Barcode / SKU</th>
            <th class="p-3 text-right">Harga Beli</th>
            <th class="p-3 text-right">Harga Jual</th>
            <th class="p-3 text-center">Stok</th>
            {#if $auth.user?.role !== 'kasir'}
              <th class="p-3 text-center w-28">Aksi</th>
            {/if}
          </tr>
        </thead>
        <tbody class="divide-y-2 divide-black dark:divide-neutral-700 text-sm font-bold text-black dark:text-white">
          {#if loading}
            <tr>
              <td colspan="7" class="p-8 text-center text-neutral-500 font-bold">
                Memuat data produk...
              </td>
            </tr>
          {:else if products.length === 0}
            <tr>
              <td colspan="7" class="p-8 text-center text-neutral-500 font-bold">
                Tidak ada produk yang cocok dengan pencarian.
              </td>
            </tr>
          {:else}
            {#each products as prd (prd.id)}
              <tr class="hover:bg-yellow-50 dark:hover:bg-[#252525] transition-colors">
                <td class="p-3">
                  <div class="font-black text-sm">{prd.name}</div>
                  <div class="text-xs text-neutral-500">{prd.unit}</div>
                </td>
                <td class="p-3">
                  <Badge variant="purple">{prd.category?.name || 'Umum'}</Badge>
                </td>
                <td class="p-3 font-mono text-xs">
                  <div>{prd.barcode || '-'}</div>
                  <div class="text-neutral-500">{prd.sku || ''}</div>
                </td>
                <td class="p-3 text-right text-neutral-500 font-mono text-xs">
                  {formatRupiah(prd.buy_price)}
                </td>
                <td class="p-3 text-right font-black font-mono text-sm">
                  {formatRupiah(prd.sell_price)}
                </td>
                <td class="p-3 text-center">
                  {#if prd.stock <= 5}
                    <Badge variant="red">{prd.stock} {prd.unit}</Badge>
                  {:else if prd.stock <= 15}
                    <Badge variant="yellow">{prd.stock} {prd.unit}</Badge>
                  {:else}
                    <Badge variant="green">{prd.stock} {prd.unit}</Badge>
                  {/if}
                </td>
                {#if $auth.user?.role !== 'kasir'}
                  <td class="p-3 text-center">
                    <div class="inline-flex gap-1.5">
                      <button
                        onclick={() => openEditModal(prd)}
                        class="px-2 py-1 text-xs font-black bg-[#00F0FF] text-black border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                      >
                        Edit
                      </button>
                      <button
                        onclick={() => confirmDelete(prd.id)}
                        class="px-2 py-1 text-xs font-black bg-[#FF5252] text-white border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                      >
                        Hapus
                      </button>
                    </div>
                  </td>
                {/if}
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
            onclick={() => { page--; loadProducts(); }}
          >
            Sebelumnya
          </Button>
          <Button
            size="sm"
            variant="outline"
            disabled={page >= totalPages}
            onclick={() => { page++; loadProducts(); }}
          >
            Berikutnya
          </Button>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- Modal Create / Edit Product -->
<Modal
  open={modalOpen}
  title={editId ? 'UBAH DATA PRODUK' : 'TAMBAH PRODUK BARU'}
  maxWidth="max-w-xl"
  onclose={() => (modalOpen = false)}
>
  <div class="space-y-4">
    <!-- Camera Scan Quick Banner -->
    <div class="p-3 bg-[#FFE600] text-black border-2 border-black flex flex-col sm:flex-row sm:items-center justify-between gap-2 shadow-[2px_2px_0px_0px_#000000]">
      <div>
        <span class="text-xs font-black uppercase tracking-wider block">Scan Barcode Otomatis</span>
        <p class="text-[11px] font-bold text-neutral-800">Foto / scan barcode kemasan di HP untuk mengisi nama barang & foto otomatis</p>
      </div>
      <button
        type="button"
        onclick={() => (scannerOpen = true)}
        class="neo-btn bg-black text-white px-3.5 py-1.5 text-xs font-black self-start sm:self-auto flex items-center gap-1.5"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        <span>Scan Kamera HP</span>
      </button>
    </div>

    {#if lookupLoading}
      <div class="p-2 bg-blue-100 border-2 border-blue-500 text-blue-800 text-xs font-bold flex items-center gap-2">
        <div class="w-3.5 h-3.5 border-2 border-blue-600 border-t-transparent animate-spin rounded-full"></div>
        Mencari data produk di database internasional (Open Food Facts)...
      </div>
    {/if}

    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {#if formImageUrl}
          <div class="sm:col-span-2 flex items-center gap-3 p-2 bg-neutral-100 dark:bg-[#252525] border-2 border-black">
            <img src={formImageUrl} alt="Preview Produk" class="w-16 h-16 object-contain bg-white border border-neutral-300" />
            <div class="text-xs">
              <span class="font-black text-[#00E676] flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
                <span>Informasi Produk Terdeteksi</span>
              </span>
              {#if formBrand}
                <span class="text-neutral-700 dark:text-neutral-300 font-bold block">Produsen / Merek: {formBrand}</span>
              {/if}
              <button
                type="button"
                onclick={() => { formImageUrl = ''; formBrand = ''; }}
                class="text-[10px] text-red-500 underline font-bold mt-1"
              >
                Hapus Foto
              </button>
            </div>
          </div>
        {/if}

        <div class="sm:col-span-2">
          <Input
            label="Nama Barang / Produk"
            placeholder="Contoh: Indomie Goreng Spesial 85g"
            bind:value={formName}
            required
          />
        </div>

        <div>
          <label for="form-category" class="text-xs font-black uppercase tracking-wider block mb-1.5">
            Kategori <span class="text-red-500">*</span>
          </label>
          <select
            id="form-category"
            bind:value={formCategoryId}
            class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
            required
          >
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label for="form-unit" class="text-xs font-black uppercase tracking-wider block mb-1.5">Satuan</label>
          <select
            id="form-unit"
            bind:value={formUnit}
            class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
          >
            {#each units as u}
              <option value={u}>{u}</option>
            {/each}
          </select>
        </div>

        <div>
          <div class="flex gap-1.5 items-end">
            <div class="flex-1">
              <Input
                label="Barcode (EAN-13 / UPC)"
                placeholder="Scan atau ketik barcode..."
                bind:value={formBarcode}
                onchange={() => { if (formBarcode) handleBarcodeDetected(formBarcode); }}
              />
            </div>
            <button
              type="button"
              onclick={() => (scannerOpen = true)}
              title="Buka Kamera Barcode"
              class="neo-btn bg-[#00F0FF] text-black px-3 py-2.5 mb-0.5 border-2 text-sm flex items-center justify-center"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
          </div>
        </div>

      <div>
        <div class="flex items-center justify-between mb-1.5">
          <label for="form-sku" class="text-xs font-black uppercase tracking-wider">Kode SKU</label>
          <button
            type="button"
            onclick={generateRandomBarcode}
            class="text-[10px] font-black underline text-blue-600 dark:text-blue-400"
          >
            Auto Barcode & SKU
          </button>
        </div>
        <input
          id="form-sku"
          placeholder="Auto-generate jika kosong"
          bind:value={formSKU}
          class="neo-input text-xs"
        />
      </div>

      <div>
        <Input
          label="Harga Beli (Modal)"
          type="number"
          placeholder="0"
          bind:value={formBuyPrice}
        />
      </div>

      <div>
        <Input
          label="Harga Jual"
          type="number"
          placeholder="0"
          bind:value={formSellPrice}
          required
        />
      </div>

      <div>
        <Input
          label="Jumlah Stok Fisik"
          type="number"
          placeholder="0"
          bind:value={formStock}
          required
        />
      </div>

      <!-- Live Margin Info -->
      <div class="flex flex-col justify-end">
        <div class="p-3 bg-neutral-100 dark:bg-[#292929] border-2 border-black dark:border-neutral-500 text-xs font-black">
          <div>Estimasi Profit per Unit:</div>
          <div class="text-[#00E676] font-mono text-sm">
            {formatRupiah(formSellPrice - formBuyPrice)}
            {#if formSellPrice > 0}
              ({Math.round(((formSellPrice - formBuyPrice) / formSellPrice) * 100)}%)
            {/if}
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-2 pt-3 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700">
      <Button variant="outline" onclick={() => (modalOpen = false)}>Batal</Button>
      <Button variant="primary" type="submit" loading={formLoading}>Simpan Produk</Button>
    </div>
  </form>
  </div>
</Modal>

<!-- Modal Delete -->
<Modal
  open={deleteModalOpen}
  title="HAPUS PRODUK"
  maxWidth="max-w-sm"
  onclose={() => (deleteModalOpen = false)}
>
  <p class="font-bold text-sm mb-4">Hapus produk ini dari daftar inventaris toko?</p>
  <div class="flex justify-end gap-2">
    <Button variant="outline" onclick={() => (deleteModalOpen = false)}>Batal</Button>
    <Button variant="danger" onclick={handleDelete}>Hapus</Button>
  </div>
</Modal>

<!-- Barcode Camera Scanner Modal -->
<BarcodeScannerModal
  open={scannerOpen}
  onclose={() => (scannerOpen = false)}
  onscan={handleBarcodeDetected}
/>
