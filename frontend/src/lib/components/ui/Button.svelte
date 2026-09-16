<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'primary' | 'accent' | 'success' | 'danger' | 'dark' | 'outline';
    size?: 'sm' | 'md' | 'lg';
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    loading?: boolean;
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    variant = 'primary',
    size = 'md',
    type = 'button',
    disabled = false,
    loading = false,
    class: className = '',
    onclick,
    children
  }: Props = $props();

  const variantClasses = {
    primary: 'bg-[#FFE600] text-black hover:bg-[#ffe933]',
    accent: 'bg-[#00F0FF] text-black hover:bg-[#33f3ff]',
    success: 'bg-[#00E676] text-black hover:bg-[#33eb91]',
    danger: 'bg-[#FF5252] text-white hover:bg-[#ff7575]',
    dark: 'bg-black text-white hover:bg-neutral-800 dark:bg-white dark:text-black dark:hover:bg-neutral-200',
    outline: 'bg-white text-black hover:bg-neutral-100 dark:bg-[#222] dark:text-white dark:hover:bg-[#2c2c2c]'
  };

  const sizeClasses = {
    sm: 'text-xs px-3 py-1.5 gap-1.5',
    md: 'text-sm px-4 py-2.5 gap-2',
    lg: 'text-base px-6 py-3.5 gap-2.5 font-extrabold'
  };
</script>

<button
  {type}
  {disabled}
  {onclick}
  class="neo-btn {variantClasses[variant]} {sizeClasses[size]} {disabled || loading ? 'opacity-50 cursor-not-allowed transform-none shadow-none' : ''} {className}"
>
  {#if loading}
    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  {/if}
  {#if children}
    {@render children()}
  {/if}
</button>
