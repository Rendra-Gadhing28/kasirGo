<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { auth } from '$lib/stores/auth';
  import type { Customer } from '$lib/types';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import MemberStandeeModal from '$lib/components/pos/MemberStandeeModal.svelte';

  let customers = $state<Customer[]>([]);
  let loading = $state(true);
  let search = $state('');

  let standeeOpen = $state(false);

  let modalOpen = $state(false);
  let editId = $state<number | null>(null);
  let formName = $state('');
  let formPhone = $state('');
  let formEmail = $state('');
  let formAddress = $state('');
  let formLoading = $state(false);

  let deleteModalOpen = $state(false);
  let deleteId = $state<number | null>(null);

  async function loadCustomers() {
    loading = true;
    try {
      let url = '/customers';
      if (search) url += `?search=${encodeURIComponent(search)}`;
      const res = await api.get<Customer[]>(url);
      customers = res.data || [];
    } catch (e: any) {
      toast.error(e.message || 'Gagal memuat data pelanggan');
    } finally {
      loading = false;
    }
  }

  function openCreateModal() {
    editId = null;
    formName = '';
    formPhone = '';
    formEmail = '';
    formAddress = '';
    modalOpen = true;
  }

  function openEditModal(c: Customer) {
    editId = c.id;
    formName = c.name;
    formPhone = c.phone;
    formEmail = c.email;
    formAddress = c.address;
    modalOpen = true;
  }

  async function handleSubmit() {
    if (!formName.trim()) {
      toast.error('Nama pelanggan wajib diisi');
      return;
    }
    formLoading = true;
    try {
      const payload = {
        name: formName,
        phone: formPhone,
        email: formEmail,
        address: formAddress
      };
      if (editId) {
        await api.put(`/customers/${editId}`, payload);
        toast.success('Data pelanggan diperbarui');
      } else {
        await api.post('/customers', payload);
        toast.success('Pelanggan baru terdaftar');
      }
      modalOpen = false;
      await loadCustomers();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menyimpan pelanggan');
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
      await api.delete(`/customers/${deleteId}`);
      toast.success('Pelanggan berhasil dihapus');
      deleteModalOpen = false;
      await loadCustomers();
    } catch (e: any) {
      toast.error(e.message || 'Gagal menghapus pelanggan');
    }
  }

  onMount(loadCustomers);
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b-3 border-black dark:border-white pb-4">
    <div>
      <h1 class="text-2xl font-black uppercase tracking-tight">Data Pelanggan & Loyalitas</h1>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">Poin otomatis bertambah 1 poin tiap belanja Rp 10.000</p>
    </div>
    <div class="flex flex-wrap gap-2">
      <Button variant="accent" onclick={() => (standeeOpen = true)}>
        📱 CETAK QR MEMBER
      </Button>
      <Button variant="primary" onclick={openCreateModal}>
        + PELANGGAN BARU
      </Button>
    </div>
  </div>

  <div class="flex gap-2 max-w-md">
    <Input
      placeholder="Cari nama pelanggan atau nomor WhatsApp..."
      bind:value={search}
      onkeydown={(e) => { if (e.key === 'Enter') loadCustomers(); }}
    />
    <Button variant="accent" onclick={loadCustomers}>Cari</Button>
  </div>

  <div class="neo-box overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-[#FFE600] text-black border-b-3 border-black font-black text-xs uppercase tracking-wider">
            <th class="p-3.5">Nama Pelanggan</th>
            <th class="p-3.5">Kontak / WhatsApp</th>
            <th class="p-3.5">Email</th>
            <th class="p-3.5 text-center">Poin Loyalitas</th>
            <th class="p-3.5">Alamat</th>
            <th class="p-3.5 text-center w-28">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y-2 divide-black dark:divide-neutral-700 text-sm font-bold">
          {#if loading}
            <tr>
              <td colspan="6" class="p-8 text-center text-neutral-500 font-bold">Memuat pelanggan...</td>
            </tr>
          {:else if customers.length === 0}
            <tr>
              <td colspan="6" class="p-8 text-center text-neutral-500 font-bold">Tidak ada data pelanggan.</td>
            </tr>
          {:else}
            {#each customers as cust (cust.id)}
              <tr class="hover:bg-yellow-50 dark:hover:bg-[#252525] transition-colors">
                <td class="p-3.5 font-black">{cust.name}</td>
                <td class="p-3.5 font-mono text-xs">{cust.phone || '-'}</td>
                <td class="p-3.5 text-xs text-neutral-500">{cust.email || '-'}</td>
                <td class="p-3.5 text-center">
                  <Badge variant="cyan">{cust.points} Poin</Badge>
                </td>
                <td class="p-3.5 text-xs text-neutral-600 dark:text-neutral-400">{cust.address || '-'}</td>
                <td class="p-3.5 text-center">
                  <div class="inline-flex gap-1.5">
                    <button
                      onclick={() => openEditModal(cust)}
                      class="px-2 py-1 text-xs font-black bg-[#00F0FF] text-black border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                    >
                      Edit
                    </button>
                    {#if $auth.user?.role !== 'kasir'}
                      <button
                        onclick={() => confirmDelete(cust.id)}
                        class="px-2 py-1 text-xs font-black bg-[#FF5252] text-white border-2 border-black shadow-[2px_2px_0px_0px_#000000]"
                      >
                        Hapus
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
  </div>
</div>

<Modal
  open={modalOpen}
  title={editId ? 'UBAH PELANGGAN' : 'PELANGGAN BARU'}
  onclose={() => (modalOpen = false)}
>
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-3.5">
    <Input label="Nama Lengkap" placeholder="Budi Santoso" bind:value={formName} required />
    <Input label="Nomor WhatsApp" placeholder="081234567890" bind:value={formPhone} />
    <Input label="Email" type="email" placeholder="budi@gmail.com" bind:value={formEmail} />
    <div class="flex flex-col gap-1.5">
      <label for="cust-address" class="text-xs font-black uppercase tracking-wider">Alamat</label>
      <textarea id="cust-address" rows="2" bind:value={formAddress} placeholder="Alamat rumah / domisili" class="neo-input"></textarea>
    </div>
    <div class="flex justify-end gap-2 pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700">
      <Button variant="outline" onclick={() => (modalOpen = false)}>Batal</Button>
      <Button variant="primary" type="submit" loading={formLoading}>Simpan</Button>
    </div>
  </form>
</Modal>

<Modal open={deleteModalOpen} title="HAPUS PELANGGAN" maxWidth="max-w-sm" onclose={() => (deleteModalOpen = false)}>
  <p class="font-bold text-sm mb-4">Hapus data pelanggan ini dari sistem toko?</p>
  <div class="flex justify-end gap-2">
    <Button variant="outline" onclick={() => (deleteModalOpen = false)}>Batal</Button>
    <Button variant="danger" onclick={handleDelete}>Hapus</Button>
  </div>
</Modal>

<!-- Standee Print Modal -->
<MemberStandeeModal
  open={standeeOpen}
  onclose={() => (standeeOpen = false)}
/>
