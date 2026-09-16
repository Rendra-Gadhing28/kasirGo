<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    open: boolean;
    onclose: () => void;
    onscan: (barcode: string) => void;
  }

  let { open, onclose, onscan }: Props = $props();

  let scannerContainerId = 'barcode-reader-box';
  let html5QrCode: Html5Qrcode | null = null;
  let isScanning = $state(false);
  let scanError = $state('');
  let fileInput: HTMLInputElement | null = null;

  async function startScanner() {
    scanError = '';
    try {
      if (!html5QrCode) {
        html5QrCode = new Html5Qrcode(scannerContainerId, {
          formatsToSupport: [
            Html5QrcodeSupportedFormats.EAN_13,
            Html5QrcodeSupportedFormats.EAN_8,
            Html5QrcodeSupportedFormats.UPC_A,
            Html5QrcodeSupportedFormats.UPC_E,
            Html5QrcodeSupportedFormats.CODE_128,
            Html5QrcodeSupportedFormats.CODE_39,
            Html5QrcodeSupportedFormats.QR_CODE
          ],
          verbose: false
        });
      }

      await html5QrCode.start(
        { facingMode: 'environment' },
        {
          fps: 15,
          qrbox: { width: 260, height: 160 }
        },
        (decodedText) => {
          handleSuccess(decodedText);
        },
        (errorMessage) => {
          // ignore frame read errors while scanning
        }
      );
      isScanning = true;
    } catch (err: any) {
      console.warn('Camera scanner error:', err);
      scanError = 'Tidak dapat mengakses kamera. Pastikan izin kamera telah diberikan atau gunakan opsi ambil foto di bawah.';
      isScanning = false;
    }
  }

  async function stopScanner() {
    if (html5QrCode && isScanning) {
      try {
        await html5QrCode.stop();
      } catch (e) {
        // ignore
      }
      isScanning = false;
    }
  }

  function handleSuccess(barcode: string) {
    if ('vibrate' in navigator) {
      try {
        navigator.vibrate(100);
      } catch (e) {}
    }
    stopScanner();
    onscan(barcode);
    onclose();
  }

  async function handleFileScan(e: Event) {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      const imageFile = target.files[0];
      try {
        if (!html5QrCode) {
          html5QrCode = new Html5Qrcode(scannerContainerId);
        }
        const decodedText = await html5QrCode.scanFile(imageFile, false);
        handleSuccess(decodedText);
      } catch (err: any) {
        scanError = 'Barcode tidak terdeteksi pada gambar. Pastikan gambar barcode terang, jelas, dan tegak lurus.';
      }
    }
  }

  $effect(() => {
    if (open) {
      // Small timeout to allow DOM modal to render the scannerContainerId element
      setTimeout(() => {
        startScanner();
      }, 200);
    } else {
      stopScanner();
    }
  });

  onDestroy(() => {
    stopScanner();
  });
</script>

<Modal
  {open}
  title="SCAN BARCODE PRODUK (KAMERA HP)"
  maxWidth="max-w-md"
  onclose={() => { stopScanner(); onclose(); }}
>
  <div class="space-y-4 text-center">
    <p class="text-xs text-neutral-600 dark:text-neutral-300 font-bold">
      Arahkan kamera ke garis barcode produk (EAN-13 / UPC / QR).
    </p>

    <!-- Scanner Viewport -->
    <div class="relative overflow-hidden bg-black border-3 border-black shadow-[4px_4px_0px_0px_#000000] min-h-[220px] rounded-sm flex items-center justify-center">
      <div id={scannerContainerId} class="w-full h-full min-h-[220px]"></div>

      {#if !isScanning && !scanError}
        <div class="absolute inset-0 flex items-center justify-center text-white text-xs font-black">
          Menyiapkan kamera...
        </div>
      {/if}
    </div>

    {#if scanError}
      <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 text-xs font-bold text-left">
        {scanError}
      </div>
    {/if}

    <!-- Alternative: Snap photo / Gallery option -->
    <div class="pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700 flex flex-col gap-2">
      <input
        type="file"
        accept="image/*"
        capture="environment"
        class="hidden"
        bind:this={fileInput}
        onchange={handleFileScan}
      />

      <div class="grid grid-cols-2 gap-2">
        <button
          type="button"
          onclick={() => fileInput?.click()}
          class="neo-btn bg-[#00F0FF] text-black py-2.5 text-xs font-black"
        >
          📷 Ambil Foto Barcode
        </button>
        <button
          type="button"
          onclick={startScanner}
          class="neo-btn bg-[#FFE600] text-black py-2.5 text-xs font-black"
        >
          🔄 Ulangi Kamera
        </button>
      </div>
    </div>
  </div>
</Modal>
