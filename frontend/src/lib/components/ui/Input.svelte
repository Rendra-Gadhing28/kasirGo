<script lang="ts">
  interface Props {
    id?: string;
    label?: string;
    type?: string;
    placeholder?: string;
    value?: string | number;
    error?: string;
    required?: boolean;
    disabled?: boolean;
    class?: string;
    oninput?: (e: Event) => void;
    onchange?: (e: Event) => void;
    onkeydown?: (e: KeyboardEvent) => void;
  }

  let {
    id,
    label,
    type = 'text',
    placeholder = '',
    value = $bindable(''),
    error,
    required = false,
    disabled = false,
    class: className = '',
    oninput,
    onchange,
    onkeydown
  }: Props = $props();

  const inputId = $derived(id || ('input-' + Math.random().toString(36).substring(2, 9)));
</script>

<div class="flex flex-col gap-1.5 w-full">
  {#if label}
    <label for={inputId} class="text-xs font-black uppercase tracking-wider text-black dark:text-neutral-200">
      {label} {#if required}<span class="text-red-500 font-black">*</span>{/if}
    </label>
  {/if}
  <input
    id={inputId}
    {type}
    {placeholder}
    {disabled}
    {required}
    bind:value={value}
    {oninput}
    {onchange}
    {onkeydown}
    class="neo-input {error ? 'border-red-500' : ''} {className}"
  />
  {#if error}
    <span class="text-xs font-bold text-red-500 mt-0.5">{error}</span>
  {/if}
</div>
