<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode';
  import Modal from '$lib/components/ui/Modal.svelte';
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    open: boolean;
    title?: string;
    continuous?: boolean;
    onclose: () => void;
    onscan: (barcode: string) => void;
  }

  let {
    open,
    title = 'SCAN BARCODE PRODUK',
    continuous = false,
    onclose,
    onscan
  }: Props = $props();

  let scannerContainerId = 'barcode-reader-box';
  let html5QrCode: Html5Qrcode | null = null;
  let isScanning = $state(false);
  let scanError = $state('');
  let fileInput: HTMLInputElement | null = null;

  let lastScanned = $state('');
  let lastScanTime = $state(0);
  let totalScanned = $state(0);

  function playBeepSound() {
    try {
      const AudioCtx = window.AudioContext || (window as any).webkitAudioContext;
      if (!AudioCtx) return;
      const ctx = new AudioCtx();
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'sine';
      osc.frequency.setValueAtTime(1800, ctx.currentTime);
      gain.gain.setValueAtTime(0.2, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.12);

      osc.connect(gain);
      gain.connect(ctx.destination);

      osc.start();
      osc.stop(ctx.currentTime + 0.12);
    } catch (_) {}
  }

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
        () => {
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
    const clean = barcode.trim();
    if (!clean) return;

    const now = Date.now();
    // Debounce duplicate scans within 1.5 seconds in continuous mode
    if (continuous && clean === lastScanned && now - lastScanTime < 1500) {
      return;
    }

    lastScanned = clean;
    lastScanTime = now;
    totalScanned++;

    playBeepSound();
    if ('vibrate' in navigator) {
      try {
        navigator.vibrate(100);
      } catch (_) {}
    }

    onscan(clean);

    if (!continuous) {
      stopScanner();
      onclose();
    }
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
      totalScanned = 0;
      lastScanned = '';
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
  {title}
  maxWidth="max-w-md"
  onclose={() => {
    stopScanner();
    onclose();
  }}
>
  <div class="space-y-3.5 text-center">
    <div class="flex items-center justify-between text-xs font-bold text-neutral-600 dark:text-neutral-300">
      <span>Arahkan kamera ke barcode produk atau QR</span>
      {#if continuous && totalScanned > 0}
        <span class="bg-[#00E676] text-black px-2 py-0.5 border border-black font-black text-[11px]">
          {totalScanned}x Terscan
        </span>
      {/if}
    </div>

    <!-- Scanner Viewport -->
    <div class="relative overflow-hidden bg-black border-3 border-black dark:border-white shadow-[4px_4px_0px_0px_#000000] dark:shadow-[4px_4px_0px_0px_#ffffff] min-h-[220px] rounded-sm flex items-center justify-center">
      <div id={scannerContainerId} class="w-full h-full min-h-[220px]"></div>

      {#if !isScanning && !scanError}
        <div class="absolute inset-0 flex items-center justify-center text-white text-xs font-black">
          Menyiapkan kamera...
        </div>
      {/if}

      <!-- Laser Guide Line -->
      {#if isScanning}
        <div class="absolute inset-x-0 top-1/2 -translate-y-1/2 h-0.5 bg-red-500 shadow-[0_0_8px_2px_#ef4444] pointer-events-none opacity-80"></div>
      {/if}
    </div>

    <!-- Last Scanned Indicator in Continuous Mode -->
    {#if continuous && lastScanned}
      <div class="p-2 bg-[#FFE600] text-black border-2 border-black font-black text-xs flex items-center justify-between shadow-[2px_2px_0px_0px_#000000]">
        <span>Terscan:</span>
        <span class="font-mono text-xs">{lastScanned}</span>
      </div>
    {/if}

    {#if scanError}
      <div class="p-3 bg-red-100 border-2 border-red-500 text-red-700 text-xs font-bold text-left">
        {scanError}
      </div>
    {/if}

    <!-- Buttons -->
    <div class="pt-2 border-t-2 border-dashed border-neutral-300 dark:border-neutral-700 flex flex-col gap-2">
      <input
        type="file"
        accept="image/*"
        capture="environment"
        class="hidden"
        bind:this={fileInput}
        onchange={handleFileScan}
      />

      <div class="grid {continuous ? 'grid-cols-3' : 'grid-cols-2'} gap-2">
        <button
          type="button"
          onclick={() => fileInput?.click()}
          class="neo-btn bg-[#00F0FF] text-black py-2 text-xs font-black flex items-center justify-center gap-1.5"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>File / Foto</span>
        </button>
        <button
          type="button"
          onclick={startScanner}
          class="neo-btn bg-[#FFE600] text-black py-2 text-xs font-black flex items-center justify-center gap-1.5"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>Refresh</span>
        </button>
        {#if continuous}
          <button
            type="button"
            onclick={() => { stopScanner(); onclose(); }}
            class="neo-btn bg-[#00E676] text-black py-2 text-xs font-black flex items-center justify-center gap-1.5"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
            </svg>
            <span>Selesai</span>
          </button>
        {/if}
      </div>
    </div>
  </div>
</Modal>
