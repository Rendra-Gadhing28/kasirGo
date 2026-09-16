<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'yellow' | 'cyan' | 'green' | 'red' | 'purple' | 'gray';
    class?: string;
    children?: Snippet;
  }

  let {
    variant = 'yellow',
    class: className = '',
    children
  }: Props = $props();

  const variantStyles = {
    yellow: 'bg-[#FFE600] text-black border-black dark:border-white',
    cyan: 'bg-[#00F0FF] text-black border-black dark:border-white',
    green: 'bg-[#00E676] text-black border-black dark:border-white',
    red: 'bg-[#FF5252] text-white border-black dark:border-white',
    purple: 'bg-[#B388FF] text-black border-black dark:border-white',
    gray: 'bg-neutral-200 dark:bg-neutral-800 text-black dark:text-white border-black dark:border-white'
  };
</script>

<span
  class="inline-flex items-center px-2.5 py-0.5 border-2 text-xs font-black uppercase tracking-wider shadow-[2px_2px_0px_0px_#000000] dark:shadow-[2px_2px_0px_0px_#ffffff] {variantStyles[variant]} {className}"
>
  {#if children}
    {@render children()}
  {/if}
</span>
