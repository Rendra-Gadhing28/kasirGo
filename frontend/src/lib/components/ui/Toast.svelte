<script lang="ts">
  import { toast } from '$lib/stores/toast';
</script>

{#if $toast.length > 0}
  <div class="fixed bottom-5 right-5 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
    {#each $toast as item (item.id)}
      <div
        class="pointer-events-auto border-3 border-black p-3.5 font-bold shadow-[4px_4px_0px_0px_#000000] flex items-center justify-between transition-all duration-200
        {item.type === 'success' ? 'bg-[#00E676] text-black' : ''}
        {item.type === 'error' ? 'bg-[#FF5252] text-white' : ''}
        {item.type === 'info' ? 'bg-[#00F0FF] text-black' : ''}"
      >
        <span class="text-sm">{item.message}</span>
        <button
          onclick={() => toast.remove(item.id)}
          class="ml-3 font-black hover:opacity-75 cursor-pointer text-lg leading-none"
        >
          ×
        </button>
      </div>
    {/each}
  </div>
{/if}
