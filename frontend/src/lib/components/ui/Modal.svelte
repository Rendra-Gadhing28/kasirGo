<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    open?: boolean;
    title?: string;
    maxWidth?: string;
    onclose: () => void;
    children?: Snippet;
    actions?: Snippet;
  }

  let {
    open = false,
    title = '',
    maxWidth = 'max-w-lg',
    onclose,
    children,
    actions
  }: Props = $props();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && open) {
      onclose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-[2px]">
    <!-- Backdrop dismiss -->
    <button
      type="button"
      class="fixed inset-0 w-full h-full cursor-default"
      onclick={onclose}
      aria-label="Tutup modal"
    ></button>

    <!-- Modal Box -->
    <div
      class="relative z-10 w-full {maxWidth} bg-white dark:bg-[#1a1a1a] border-3 border-black dark:border-white shadow-[8px_8px_0px_0px_#000000] dark:shadow-[8px_8px_0px_0px_#ffffff] max-h-[90vh] flex flex-col overflow-hidden"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-5 py-4 border-b-3 border-black dark:border-white bg-[#FFE600] text-black">
        <h3 class="text-lg font-black tracking-tight">{title}</h3>
        <button
          onclick={onclose}
          class="w-8 h-8 flex items-center justify-center bg-black text-white hover:bg-red-500 font-black text-lg border-2 border-black transition-colors cursor-pointer"
        >
          ×
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto flex-1 text-black dark:text-white">
        {#if children}
          {@render children()}
        {/if}
      </div>

      <!-- Footer Actions -->
      {#if actions}
        <div class="px-6 py-4 border-t-3 border-black dark:border-white bg-neutral-50 dark:bg-[#222] flex justify-end gap-3">
          {@render actions()}
        </div>
      {/if}
    </div>
  </div>
{/if}
