<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';

  onMount(() => {
    if ($auth.user) {
      goto('/pos');
    } else {
      goto('/login');
    }
  });
</script>

<div class="flex items-center justify-center min-h-[50vh]">
  <p class="font-bold text-sm">Mengarahkan ke KasirPro...</p>
</div>
