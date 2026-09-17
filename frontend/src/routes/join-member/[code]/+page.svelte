<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { api } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';

  const outletCode = $derived(page.params.code || '');

  let outletInfo = $state<{ name: string; code: string; address: string; phone: string } | null>(null);
  let loadingInfo = $state(true);

  let name = $state('');
  let phone = $state('');
  let email = $state('');
  let address = $state('');
  let submitting = $state(false);
  let registeredMember = $state<any>(null);

  async function loadOutlet() {
    loadingInfo = true;
    try {
      const res = await api.get(`/public/outlet-info/${outletCode}`);
      outletInfo = res.data;
    } catch (e: any) {
      toast.error(e.message || 'Toko tidak ditemukan.');
    } finally {
      loadingInfo = false;
    }
  }

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim() || !phone.trim()) {
      toast.error('Nama dan Nomor WhatsApp wajib diisi');
      return;
    }
    submitting = true;
    try {
      const res = await api.post('/public/member-register', {
        outlet_code: outletCode,
        name,
        phone,
        email,
        address
      });
      registeredMember = res.data.customer;
      toast.success('Pendaftaran member berhasil!');
    } catch (e: any) {
      toast.error(e.message || 'Gagal mendaftar member');
    } finally {
      submitting = false;
    }
  }

  onMount(loadOutlet);
</script>

<div class="min-h-screen bg-[#FFFDF7] dark:bg-[#121212] py-8 px-4 flex flex-col items-center justify-center">
  <div class="w-full max-w-md bg-white dark:bg-[#1e1e1e] border-4 border-black dark:border-white shadow-[8px_8px_0px_0px_#000000] dark:shadow-[8px_8px_0px_0px_#ffffff] p-6 sm:p-8">
    
    {#if loadingInfo}
      <div class="text-center py-12">
        <div class="w-10 h-10 border-4 border-black border-t-[#FFE600] animate-spin mx-auto mb-3"></div>
        <p class="font-bold text-xs uppercase tracking-wider">Memuat Info Toko...</p>
      </div>
    {:else if !outletInfo}
      <div class="text-center py-8">
        <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 font-black text-sm mb-4">
          Toko dengan kode "{outletCode}" tidak ditemukan.
        </div>
        <p class="text-xs text-neutral-500">Pastikan Anda memindai kode QR standee toko yang benar.</p>
      </div>
    {:else if registeredMember}
      <!-- MEMBER SUCCESS CARD -->
      <div class="text-center space-y-4">
        <div class="inline-block bg-[#00E676] text-black px-3 py-1 border-2 border-black font-black text-xs uppercase shadow-[2px_2px_0px_0px_#000000]">
          Selamat! Kartu Member Resmi Aktif
        </div>

        <div>
          <h2 class="text-2xl font-black text-black dark:text-white uppercase">{registeredMember.name}</h2>
          <div class="inline-flex items-center gap-2 mt-1">
            <span class="text-xs font-black uppercase bg-black text-white px-2 py-0.5 border border-black">
              NO. KARTU:
            </span>
            <span class="font-mono font-black text-sm text-black dark:text-white bg-[#FFE600] px-2 py-0.5 border border-black">
              {registeredMember.member_code || ('MBR-' + String(registeredMember.id).padStart(4, '0'))}
            </span>
          </div>
          <p class="text-xs font-mono font-bold text-neutral-500 mt-1">{registeredMember.phone}</p>
          <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">{outletInfo.name}</p>
        </div>

        <!-- Bonus Points Badge -->
        <div class="p-3 bg-[#FFE600] text-black border-2 border-black shadow-[3px_3px_0px_0px_#000000] font-black">
          <span class="text-xs uppercase tracking-wider block">Saldo Poin Member:</span>
          <span class="text-2xl font-mono">{registeredMember.points} POIN</span>
          <span class="text-[10px] block mt-0.5 text-neutral-800">+1 Poin otomatis tiap belanja Rp 10.000</span>
        </div>

        <!-- Digital Member QR Code -->
        <div class="p-4 bg-white text-black border-3 border-black shadow-[4px_4px_0px_0px_#000000] inline-block">
          <img
            src="https://api.qrserver.com/v1/create-qr-code/?size=200x200&data={encodeURIComponent(registeredMember.member_code || registeredMember.phone)}"
            alt="Barcode Member"
            class="w-48 h-48 mx-auto border border-black"
          />
          <div class="mt-2 text-xs font-mono font-black text-black tracking-widest">
            {registeredMember.member_code || ('MBR-' + String(registeredMember.id).padStart(4, '0'))}
          </div>
          <div class="text-[10px] text-neutral-500 font-bold">{registeredMember.phone}</div>
        </div>

        <div class="bg-neutral-100 dark:bg-[#252525] p-3 border-2 border-black dark:border-white text-xs font-bold text-neutral-700 dark:text-neutral-300 max-w-xs mx-auto text-left space-y-1">
          <div class="font-black text-black dark:text-white uppercase flex items-center gap-1.5">
            <svg class="w-4 h-4 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Cara Menggunakan:</span>
          </div>
          <p>1. <strong>Tunjukkan QR di atas</strong> ke kasir untuk di-scan secara otomatis.</p>
          <p>2. Atau cukup <strong>sebutkan No. Kartu / No. HP</strong> ke kasir saat checkout.</p>
        </div>

        <button
          onclick={() => { registeredMember = null; }}
          class="text-xs font-black underline text-neutral-500 hover:text-black dark:hover:text-white pt-2"
        >
          Daftarkan nomor lainnya
        </button>
      </div>
    {:else}
      <!-- REGISTRATION FORM -->
      <div class="text-center mb-6">
        <Badge variant="yellow" class="mb-2">PROGRAM MEMBER LOYALITAS</Badge>
        <h2 class="text-2xl font-black uppercase tracking-tight text-black dark:text-white">
          {outletInfo.name}
        </h2>
        <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400 mt-1">
          Daftar gratis sekarang & dapatkan <span class="text-[#00E676] font-black">10 Bonus Poin</span> selamat datang!
        </p>
      </div>

      <form onsubmit={handleRegister} class="space-y-3.5">
        <Input
          label="Nama Lengkap"
          placeholder="Contoh: Budi Santoso"
          bind:value={name}
          required
        />

        <Input
          label="Nomor WhatsApp / HP"
          placeholder="Contoh: 081234567890"
          bind:value={phone}
          required
        />

        <Input
          label="Email (Opsional)"
          type="email"
          placeholder="budi@gmail.com"
          bind:value={email}
        />

        <div class="flex flex-col gap-1.5">
          <label for="reg-address" class="text-xs font-black uppercase tracking-wider text-black dark:text-neutral-200">
            Alamat Domisili (Opsional)
          </label>
          <textarea
            id="reg-address"
            rows="2"
            bind:value={address}
            placeholder="Alamat singkat..."
            class="neo-input"
          ></textarea>
        </div>

        <Button
          type="submit"
          variant="primary"
          size="lg"
          class="w-full mt-2 bg-[#FFE600] text-black"
          loading={submitting}
        >
          KLAIM KARTU MEMBER & 10 POIN
        </Button>
      </form>
    {/if}
  </div>
</div>
