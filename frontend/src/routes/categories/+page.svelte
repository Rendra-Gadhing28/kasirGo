<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { auth } from '$lib/stores/auth';
  import type { Category } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';

  let categories = $state<Category[]>([]);
  let loading = $state(true);
  let modalOpen = $state(false);
  let editId = $state<number | null>(null);

  let formName = $state('');
  let formDesc = $state('');
  let formLoading = $state(false);

  let deleteModalOpen = $state(false);
  let deleteId = $state<number | null>(null);

  async function loadCategories() {
    loading = true;
    try {
      const res = await api.get<Category[]>('/categories');
      categories = res.data || [];
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat kategori');
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    editId = null;
    formName = '';
    formDesc = '';
    modalOpen = true;
  }

  function openEditModal(cat: Category) {
    editId = cat.id;
    formName = cat.name;
    formDesc = cat.description;
    modalOpen = true;
  }

  async function handleSubmit() {
    if (!formName.trim()) {
      toast.error('Nama kategori wajib diisi');
      return;
    }
    formLoading = true;
    try {
      if (editId) {
        await api.put(`/categories/${editId}`, {
          name: formName,
          description: formDesc
        });
        toast.success('Kategori berhasil diperbarui');
      } else {
        await api.post('/categories', {
          name: formName,
          description: formDesc
        });
        toast.success('Kategori baru berhasil ditambahkan');
      }
      modalOpen = false;
      await loadCategories();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menyimpan kategori');
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
      await api.delete(`/categories/${deleteId}`);
      toast.success('Kategori berhasil dihapus');
      deleteModalOpen = false;
      await loadCategories();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menghapus kategori');
    }
  }

  onMount(loadCategories);
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Kategori Produk</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Kelola kelompok jenis barang untuk kemudahan kasir</p>
    </div>
    {#if $auth.user?.role !== 'kasir'}
      <Button variant="primary" onclick={openCreateModal}>
        + TAMBAH KATEGORI
      </Button>
    {/if}
  </div>

  <!-- Table Card -->
  <div class="neo-box overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-[#FFE600] text-black border-b-3 border-black dark:border-white font-black text-xs uppercase tracking-wider">
            <th class="p-3.5 w-16">No</th>
            <th class="p-3.5">Nama Kategori</th>
            <th class="p-3.5">Deskripsi</th>
            {#if $auth.user?.role !== 'kasir'}
              <th class="p-3.5 w-36 text-center">Aksi</th>
            {/if}
          </tr>
        </thead>
        <tbody class="divide-y-2 divide-black dark:divide-neutral-700 text-sm font-bold text-black dark:text-white">
          {#if loading}
            <tr>
              <td colspan="4" class="p-8 text-center text-neutral-500 font-bold">
                Memuat data kategori...
              </td>
            </tr>
          {:else if categories.length === 0}
            <tr>
              <td colspan="4" class="p-8 text-center text-neutral-500 font-bold">
                Belum ada data kategori. Klik tombol tambah kategori di atas.
              </td>
            </tr>
          {:else}
            {#each categories as cat, i (cat.id)}
              <tr class="hover:bg-yellow-50 dark:hover:bg-[#252525] transition-colors">
                <td class="p-3.5 font-mono">{i + 1}</td>
                <td class="p-3.5 font-black text-base">{cat.name}</td>
                <td class="p-3.5 text-neutral-600 dark:text-neutral-300">{cat.description || '-'}</td>
                {#if $auth.user?.role !== 'kasir'}
                  <td class="p-3.5 text-center">
                    <div class="inline-flex gap-2">
                      <button
                        onclick={() => openEditModal(cat)}
                        class="px-2.5 py-1 text-xs font-black bg-[#00F0FF] text-black border-2 border-black shadow-[2px_2px_0px_0px_#000000] hover:translate-x-0.5 hover:translate-y-0.5 active:shadow-none"
                      >
                        Edit
                      </button>
                      <button
                        onclick={() => confirmDelete(cat.id)}
                        class="px-2.5 py-1 text-xs font-black bg-[#FF5252] text-white border-2 border-black shadow-[2px_2px_0px_0px_#000000] hover:translate-x-0.5 hover:translate-y-0.5 active:shadow-none"
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
  </div>
</div>

<!-- Modal Create / Edit -->
<Modal
  open={modalOpen}
  title={editId ? 'UBAH KATEGORI' : 'TAMBAH KATEGORI BARU'}
  onclose={() => (modalOpen = false)}
>
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
    <Input
      label="Nama Kategori"
      placeholder="Contoh: Minuman Dingin"
      bind:value={formName}
      required
    />
    <div class="flex flex-col gap-1.5">
      <label for="cat-desc" class="text-xs font-black uppercase tracking-wider">Deskripsi Kategori</label>
      <textarea
        id="cat-desc"
        rows="3"
        bind:value={formDesc}
        placeholder="Keterangan kategori..."
        class="neo-input"
      ></textarea>
    </div>
    <div class="flex justify-end gap-2 pt-2">
      <Button variant="outline" onclick={() => (modalOpen = false)}>Batal</Button>
      <Button variant="primary" type="submit" loading={formLoading}>Simpan</Button>
    </div>
  </form>
</Modal>

<!-- Modal Confirm Delete -->
<Modal
  open={deleteModalOpen}
  title="KONFIRMASI HAPUS"
  maxWidth="max-w-sm"
  onclose={() => (deleteModalOpen = false)}
>
  <p class="font-bold text-sm mb-4">Apakah Anda yakin ingin menghapus kategori ini? Data produk yang terkait mungkin terpengaruh.</p>
  <div class="flex justify-end gap-2">
    <Button variant="outline" onclick={() => (deleteModalOpen = false)}>Batal</Button>
    <Button variant="danger" onclick={handleDelete}>Hapus Permanen</Button>
  </div>
</Modal>
