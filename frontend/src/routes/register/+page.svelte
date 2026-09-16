<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { toast } from '$lib/stores/toast';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let registerType = $state<'owner' | 'staff'>('owner');

  // Owner Form
  let storeName = $state('');
  let ownerName = $state('');
  let ownerEmail = $state('');
  let ownerPhone = $state('');
  let ownerPassword = $state('');

  // Staff Form
  let outletCode = $state('');
  let staffName = $state('');
  let staffEmail = $state('');
  let staffPassword = $state('');
  let staffRole = $state<'kasir' | 'admin'>('kasir');

  let loading = $state(false);
  let errorMessage = $state('');

  async function handleRegisterOwner(e: SubmitEvent) {
    e.preventDefault();
    if (ownerPassword.length < 6) {
      errorMessage = 'Password minimal 6 karakter';
      return;
    }
    loading = true;
    errorMessage = '';
    try {
      await auth.register({
        store_name: storeName,
        name: ownerName,
        email: ownerEmail,
        phone: ownerPhone,
        password: ownerPassword
      });
      toast.success('Pendaftaran outlet baru berhasil! Selamat datang di KasirPro.');
      goto('/pos');
    } catch (err: any) {
      errorMessage = err.message || 'Pendaftaran gagal. Silakan periksa kembali formulir.';
      toast.error(errorMessage);
    } finally {
      loading = false;
    }
  }

  async function handleRegisterStaff(e: SubmitEvent) {
    e.preventDefault();
    if (!outletCode.trim()) {
      errorMessage = 'Kode outlet/toko wajib diisi';
      return;
    }
    if (staffPassword.length < 6) {
      errorMessage = 'Password minimal 6 karakter';
      return;
    }
    loading = true;
    errorMessage = '';
    try {
      await auth.registerStaff({
        outlet_code: outletCode.trim(),
        name: staffName,
        email: staffEmail,
        password: staffPassword,
        role: staffRole
      });
      toast.success('Berhasil bergabung ke toko! Selamat bekerja.');
      goto('/pos');
    } catch (err: any) {
      errorMessage = err.message || 'Gagal bergabung. Pastikan kode toko benar.';
      toast.error(errorMessage);
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center p-4 bg-[#FFFDF7] dark:bg-[#121212]">
  <div class="w-full max-w-lg bg-white dark:bg-[#1e1e1e] border-4 border-black dark:border-white shadow-[8px_8px_0px_0px_#000000] dark:shadow-[8px_8px_0px_0px_#ffffff] p-6 sm:p-8">
    
    <!-- Branding Header -->
    <div class="text-center mb-6">
      <div class="inline-flex items-center gap-1.5 bg-[#00E676] px-3 py-1 border-3 border-black shadow-[3px_3px_0px_0px_#000000] mb-3">
        <span class="text-lg font-black tracking-tight text-black">REGISTRASI KASIRPRO</span>
      </div>
      <h2 class="text-2xl font-black uppercase tracking-tight text-black dark:text-white">Pendaftaran Akun</h2>
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400 mt-1">
        Buka outlet baru atau bergabung ke toko yang sudah terdaftar
      </p>
    </div>

    <!-- Registration Type Toggle Tabs -->
    <div class="grid grid-cols-2 gap-2 mb-6 border-b-3 border-black dark:border-white pb-3">
      <button
        type="button"
        onclick={() => (registerType = 'owner')}
        class="neo-btn py-2.5 text-xs font-black uppercase tracking-wider {registerType === 'owner' ? 'bg-[#FFE600] text-black' : 'bg-white dark:bg-[#2a2a2a] text-black dark:text-white'}"
      >
        🏪 Buka Toko (Owner)
      </button>
      <button
        type="button"
        onclick={() => (registerType = 'staff')}
        class="neo-btn py-2.5 text-xs font-black uppercase tracking-wider {registerType === 'staff' ? 'bg-[#00F0FF] text-black' : 'bg-white dark:bg-[#2a2a2a] text-black dark:text-white'}"
      >
        💼 Gabung Toko (Kasir / Staff)
      </button>
    </div>

    {#if registerType === 'owner'}
      <!-- OWNER FORM: CREATE NEW OUTLET -->
      <form onsubmit={handleRegisterOwner} class="space-y-3.5">
        <Input
          label="Nama Toko / Minimarket / Outlet"
          placeholder="Contoh: Toko Berkah Jaya"
          bind:value={storeName}
          required
        />

        <Input
          label="Nama Pemilik (Owner)"
          placeholder="Contoh: Budi Santoso"
          bind:value={ownerName}
          required
        />

        <Input
          label="Email Owner"
          type="email"
          placeholder="owner@tokoberkah.com"
          bind:value={ownerEmail}
          required
        />

        <Input
          label="Nomor WhatsApp Toko"
          placeholder="081234567890"
          bind:value={ownerPhone}
        />

        <Input
          label="Kata Sandi (Password)"
          type="password"
          placeholder="Minimal 6 karakter"
          bind:value={ownerPassword}
          required
        />

        {#if errorMessage}
          <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 text-xs font-bold">
            {errorMessage}
          </div>
        {/if}

        <Button
          type="submit"
          variant="primary"
          size="lg"
          class="w-full mt-2"
          loading={loading}
        >
          DAFTARKAN OUTLET BARU SEKARANG
        </Button>
      </form>
    {:else}
      <!-- STAFF FORM: JOIN EXISTING OUTLET VIA CODE -->
      <form onsubmit={handleRegisterStaff} class="space-y-3.5">
        <!-- Quick Demo Code Box -->
        <div class="p-3 bg-neutral-100 dark:bg-[#252525] border-2 border-black dark:border-white text-xs space-y-1.5">
          <div class="flex items-center justify-between">
            <span class="font-bold text-neutral-600 dark:text-neutral-400">Punya kode toko dari Owner?</span>
            <button
              type="button"
              onclick={() => (outletCode = 'KASIR-DEMO-001')}
              class="text-[11px] font-black underline text-blue-600 dark:text-blue-400"
            >
              Gunakan Kode Demo
            </button>
          </div>
          <p class="text-[11px] text-neutral-500">
            Mintalah kode toko ke Owner toko Anda (misal: <span class="font-mono font-bold text-black dark:text-white">KASIR-DEMO-001</span>).
          </p>
        </div>

        <Input
          label="Kode Toko / Outlet"
          placeholder="Contoh: KASIR-DEMO-001"
          bind:value={outletCode}
          required
        />

        <Input
          label="Nama Lengkap Anda"
          placeholder="Contoh: Siti Rahayu"
          bind:value={staffName}
          required
        />

        <Input
          label="Email Karyawan"
          type="email"
          placeholder="siti@gmail.com"
          bind:value={staffEmail}
          required
        />

        <Input
          label="Kata Sandi (Password)"
          type="password"
          placeholder="Minimal 6 karakter"
          bind:value={staffPassword}
          required
        />

        <div>
          <label for="staff-role-select" class="text-xs font-black uppercase tracking-wider block mb-1">
            Posisi / Peran di Toko
          </label>
          <select
            id="staff-role-select"
            bind:value={staffRole}
            class="neo-input font-bold text-sm bg-white dark:bg-[#222]"
          >
            <option value="kasir">Kasir (Operasional Kasir POS)</option>
            <option value="admin">Admin Toko (Kelola Produk & Stok)</option>
          </select>
        </div>

        {#if errorMessage}
          <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 text-xs font-bold">
            {errorMessage}
          </div>
        {/if}

        <Button
          type="submit"
          variant="accent"
          size="lg"
          class="w-full mt-2"
          loading={loading}
        >
          GABUNG KE TOKO SEBAGAI KASIR
        </Button>
      </form>
    {/if}

    <!-- Link back to login -->
    <div class="mt-6 pt-4 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700 text-center">
      <p class="text-xs font-bold text-neutral-600 dark:text-neutral-400">
        Sudah memiliki akun? 
        <a href="/login" class="text-black dark:text-[#FFE600] font-black underline hover:text-[#FFE600]">
          Masuk di Sini
        </a>
      </p>
    </div>
  </div>
</div>
