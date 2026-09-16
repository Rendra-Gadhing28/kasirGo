<script lang="ts">
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { auth } from '$lib/stores/auth';

  interface Props {
    open: boolean;
    onclose: () => void;
  }

  let { open, onclose }: Props = $props();

  const outletCode = $derived($auth.user?.outlet_code || 'KASIR-DEMO-001');
  const outletName = $derived($auth.user?.outlet_name || 'KasirPro Modern Market');
  
  let registerUrl = $state('');

  $effect(() => {
    if (typeof window !== 'undefined') {
      const origin = (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
        ? 'https://fish-warming-logos-lots.trycloudflare.com'
        : window.location.origin;
      registerUrl = `${origin}/join-member/${outletCode}`;
    }
  });

  function handlePrintStandee() {
    window.print();
  }
</script>

<Modal
  {open}
  title="QR STANDEE PENDAFTARAN MEMBER MANDIRI"
  maxWidth="max-w-md"
  {onclose}
>
  <div class="space-y-4 text-center">
    <p class="text-xs text-neutral-600 dark:text-neutral-300 font-bold print:hidden">
      Cetak kertas ini dan letakkan di meja kasir. Pelanggan cukup scan dengan kamera HP untuk mendaftar mandiri.
    </p>

    <!-- Standee Card to Print -->
    <div
      id="standee-print-area"
      class="bg-white text-black p-6 border-4 border-black shadow-[6px_6px_0px_0px_#000000] max-w-[340px] mx-auto text-center font-sans"
    >
      <div class="inline-block bg-[#FFE600] px-3 py-1 border-2 border-black font-black text-xs uppercase mb-2">
        MEMBER RESMI
      </div>
      <h3 class="text-xl font-black uppercase tracking-tight text-black">{outletName}</h3>
      <p class="text-xs font-bold text-neutral-600 mt-0.5">Scan untuk Daftar & Dapatkan Poin</p>

      <!-- QR Code Display -->
      <div class="my-4 p-3 bg-white border-3 border-black inline-block">
        {#if registerUrl}
          <img
            src="https://api.qrserver.com/v1/create-qr-code/?size=220x220&data={encodeURIComponent(registerUrl)}"
            alt="QR Pendaftaran Member"
            class="w-48 h-48 mx-auto"
          />
        {/if}
      </div>

      <div class="p-2 bg-[#00E676] border-2 border-black font-black text-xs uppercase mb-2">
        Bonus 10 Poin Selamat Datang!
      </div>

      <p class="text-[10px] font-bold text-neutral-600">
        Kode Toko: <span class="font-mono font-black text-black text-xs">{outletCode}</span>
      </p>
      <p class="text-[9px] text-neutral-400 mt-1 font-bold">kasirpro.id</p>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-2 pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700 print:hidden">
      <Button variant="outline" onclick={onclose}>Tutup</Button>
      <Button variant="primary" onclick={handlePrintStandee}>
        Cetak Standee Meja (Print)
      </Button>
    </div>
  </div>
</Modal>

<style>
  @media print {
    :global(body *) {
      visibility: hidden;
    }
    #standee-print-area, #standee-print-area * {
      visibility: visible;
    }
    #standee-print-area {
      position: fixed;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      width: 100mm;
      max-width: none;
      margin: 0;
      border: 3px solid black;
      box-shadow: none;
    }
  }
</style>
