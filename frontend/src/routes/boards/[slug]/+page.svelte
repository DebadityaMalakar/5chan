<script>
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
  
    let slug;
    $: slug = $page.params.slug;

    onMount(()=>{
        document.body.className += "bg-gray-900"
    })
  </script>

<div class="w-full flex items-center justify-center">
    <p class="text-xl text-white">{slug}</p>
</div>
  