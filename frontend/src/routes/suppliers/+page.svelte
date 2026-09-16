<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { Supplier } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';

  let suppliers = $state<Supplier[]>([]);
  let loading = $state(true);
  let search = $state('');

  let modalOpen = $state(false);
  let editId = $state<number | null>(null);
  let formName = $state('');
  let formCP = $state('');
  let formPhone = $state('');
  let formEmail = $state('');
  let formAddress = $state('');
  let formLoading = $state(false);

  let deleteModalOpen = $state(false);
  let deleteId = $state<number | null>(null);

  async function loadSuppliers() {
    loading = true;
    try {
      let url = '/suppliers';
      if (search) url += `?search=${encodeURIComponent(search)}`;
      const res = await api.get<Supplier[]>(url);
      suppliers = res.data || [];
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat supplier');
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    editId = null;
    formName = '';
    formCP = '';
    formPhone = '';
    formEmail = '';
    formAddress = '';
    modalOpen = true;
  }

  function openEditModal(s: Supplier) {
    editId = s.id;
    formName = s.name;
    formCP = s.contact_person;
    formPhone = s.phone;
    formEmail = s.email;
    formAddress = s.address;
    modalOpen = true;
  }

  async function handleSubmit() {
    if (!formName.trim()) {
      toast.error('Nama distributor/supplier wajib diisi');
      return;
    }
    formLoading = true;
    try {
      const payload = {
        name: formName,
        contact_person: formCP,
        phone: formPhone,
        email: formEmail,
        address: formAddress
      };
      if (editId) {
        await api.put(`/suppliers/${editId}`, payload);
        toast.success('Supplier berhasil diperbarui');
      } else {
        await api.post('/suppliers', payload);
        toast.success('Supplier baru berhasil didaftarkan');
      }
      modalOpen = false;
      await loadSuppliers();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menyimpan supplier');
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
      await api.delete(`/suppliers/${deleteId}`);
      toast.success('Supplier berhasil dihapus');
      deleteModalOpen = false;
      await loadSuppliers();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menghapus supplier');
    }
  }

  onMount(loadSuppliers);
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Data Supplier & Distributor</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Pemasok bahan baku dan barang dagangan toko</p>
    </div>
    <Button variant="primary" onclick={openCreateModal}>
      + TAMBAH SUPPLIER BARU
    </Button>
  </div>

  <div class="flex gap-2 max-w-md">
    <Input
      placeholder="Cari nama supplier atau nama kontak..."
      bind:value={search}
      onkeydown={(e) => { if (e.key === 'Enter') loadSuppliers(); }}
    />
    <Button variant="accent" onclick={loadSuppliers}>Cari</Button>
  </div>

  <div class="neo-box overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-[#FFE600] text-black border-b-3 border-black font-black text-xs uppercase tracking-wider">
            <th class="p-3.5">Nama Perusahaan / Supplier</th>
            <th class="p-3.5">PIC / Sales</th>
            <th class="p-3.5">Telepon</th>
            <th class="p-3.5">Email</th>
            <th class="p-3.5">Alamat Gudang</th>
            <th class="p-3.5 text-center w-28">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y-2 divide-black dark:divide-neutral-700 text-sm font-bold">
          {#if loading}
            <tr><td colspan="6" class="p-8 text-center text-neutral-500 font-bold">Memuat supplier...</td></tr>
          {:else if suppliers.length === 0}
            <tr><td colspan="6" class="p-8 text-center text-neutral-500 font-bold">Belum ada supplier terdaftar.</td></tr>
          {:else}
            {#each suppliers as sup (sup.id)}
              <tr class="hover:bg-yellow-50 dark:hover:bg-[#252525] transition-colors">
                <td class="p-3.5 font-black">{sup.name}</td>
                <td class="p-3.5 font-mono text-xs">{sup.contact_person || '-'}</td>
                <td class="p-3.5 font-mono text-xs">{sup.phone || '-'}</td>
                <td class="p-3.5 text-xs text-neutral-500">{sup.email || '-'}</td>
                <td class="p-3.5 text-xs text-neutral-600 dark:text-neutral-400">{sup.address || '-'}</td>
                <td class="p-3.5 text-center">
                  <div class="inline-flex gap-1.5">
                    <button
                      onclick={() => openEditModal(sup)}
                      class="px-2 py-1 text-xs font-black bg-[#00F0FF] text-black border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                    >
                      Edit
                    </button>
                    <button
                      onclick={() => confirmDelete(sup.id)}
                      class="px-2 py-1 text-xs font-black bg-[#FF5252] text-white border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                    >
                      Hapus
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<Modal open={modalOpen} title={editId ? 'UBAH DATA SUPPLIER' : 'SUPPLIER BARU'} onclose={() => (modalOpen = false)}>
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-3.5">
    <Input label="Nama Perusahaan / Supplier" placeholder="PT Sukses Pangan Makmur" bind:value={formName} required />
    <Input label="Kontak Person (PIC)" placeholder="Bpk. Rahmat" bind:value={formCP} />
    <Input label="Nomor Telepon Kantor / WA" placeholder="021-12345678" bind:value={formPhone} />
    <Input label="Email" type="email" placeholder="order@supplier.co.id" bind:value={formEmail} />
    <div class="flex flex-col gap-1.5">
      <label for="sup-address" class="text-xs font-black uppercase tracking-wider">Alamat Gudang</label>
      <textarea id="sup-address" rows="2" bind:value={formAddress} placeholder="Alamat lengkap distributor" class="neo-input"></textarea>
    </div>
    <div class="flex justify-end gap-2 pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700">
      <Button variant="outline" onclick={() => (modalOpen = false)}>Batal</Button>
      <Button variant="primary" type="submit" loading={formLoading}>Simpan</Button>
    </div>
  </form>
</Modal>

<Modal open={deleteModalOpen} title="HAPUS SUPPLIER" maxWidth="max-w-sm" onclose={() => (deleteModalOpen = false)}>
  <p class="font-bold text-sm mb-4">Hapus supplier ini dari daftar toko?</p>
  <div class="flex justify-end gap-2">
    <Button variant="outline" onclick={() => (deleteModalOpen = false)}>Batal</Button>
    <Button variant="danger" onclick={handleDelete}>Hapus</Button>
  </div>
</Modal>
