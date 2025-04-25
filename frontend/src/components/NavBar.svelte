<script>
    import Logo from "../assets/5clover-blackclover-removebg-preview.png";
    import { onMount } from "svelte";
    
    let mobileMenuOpen = false;
    let isScrolled = false;
    
    function toggleMobileMenu() {
      mobileMenuOpen = !mobileMenuOpen;
    }
    
    onMount(() => {
      const handleResize = () => {
        if (window.innerWidth >= 768) {
          mobileMenuOpen = false;
        }
      };
      
      const handleScroll = () => {
        isScrolled = window.scrollY > 10;
      };
      
      window.addEventListener('resize', handleResize);
      window.addEventListener('scroll', handleScroll);
      
      return () => {
        window.removeEventListener('resize', handleResize);
        window.removeEventListener('scroll', handleScroll);
      };
    });
  </script>
  
  <header class="sticky top-0 z-50 flex flex-col md:flex-row items-center justify-between p-4 lg:p-3 bg-gray-900 text-white w-full transition-all duration-200 shadow-md {isScrolled ? 'shadow-lg' : ''}">
    <!-- Profile Button (left side corner) -->
    <div class="w-full md:w-auto mb-4 md:mb-0 flex justify-between items-center">
      <button class="bg-transparent border border-white text-white px-4 py-2 hover:bg-gray-400 hover:text-red-900">
        Profile
      </button>
      
      <!-- Mobile menu button - only visible on small screens -->
      <button class="md:hidden text-white" on:click={toggleMobileMenu} aria-label="hamburger">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={mobileMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"}></path>
        </svg>
      </button>
    </div>
    
    <!-- Navigation Links (center) -->
    <div class={`${mobileMenuOpen ? 'flex' : 'hidden'} md:flex flex-col md:flex-row space-y-3 md:space-y-0 md:space-x-6 w-full md:w-auto items-center mb-4 md:mb-0`}>
      <a href="#all-boards" class="hover:text-red-500">All Boards</a>
      <a href="#random" class="hover:text-red-500">Random (/b/)</a>
      <a href="#catalog" class="hover:text-red-500">Catalog</a>
      <a href="#support" class="hover:text-red-500">Support</a>
      <a href="#search" class="hover:text-red-500">Search</a>
    </div>
  
    <!-- Site Logo and Text (properly right-aligned) -->
    <div class="flex items-center justify-center md:justify-end space-x-2">
      <img src={Logo} alt="Site Logo" class="w-16 sm:w-20 md:w-24 h-auto">
      <p class="text-lg sm:text-xl md:text-2xl font-bold text-white">5chan</p>
    </div>
  </header>
  
  <style>
    /* Add this style to ensure the body content has appropriate margin/padding */
    :global(main), :global(.main-content) {
      padding-top: 8px;
    }
  </style>