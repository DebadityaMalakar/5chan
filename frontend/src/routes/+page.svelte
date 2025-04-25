<svelte:head>
    <title>5chan - successor to the true chaos</title>
</svelte:head>

<script>
    import { onMount } from "svelte";
    import NavBar from "../components/NavBar.svelte";

    /**
     * @type {boolean} Controls visibility of the credentials modal
     */
    let showModal = false;
    
    /**
     * @type {Object|null} Stores the randomly generated user credentials
     * @property {string} username - The generated username
     * @property {string} password - The generated password
     * @property {string} expires_at - ISO string of expiration date
     */
    let userData = null;
    
    /**
     * @type {string|null} Stores error messages for display
     */
    let error = null;
    
    /**
     * @type {Object} Tracks copy status for credential fields
     * @property {boolean} username - Whether username was copied
     * @property {boolean} password - Whether password was copied
     */
    let copied = {
        username: false,
        password: false
    };
    
    
    import {env} from '$env/dynamic/public';
    /**
     * @constant {string} Base URL for backend API from environment variables
     */
    const backendAPI = env.PUBLIC_BACKEND_API;
    /**
     * Initializes the page on mount
     * @returns {void}
     */
    onMount(() => {
        document.body.className += "bg-gray-900"
    });

    /**
     * Redirects user to the login page
     * @returns {void}
     */
    function redirectToLink() {
        window.location.href = "/login";
    }

    /**
     * Creates a random anonymous account via backend API
     * @async
     * @returns {Promise<void>}
     * @throws {Error} When account creation fails
     */
    async function createRandomAccount() {
        try {
            const response = await fetch(`${backendAPI}/random`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (!response.ok) {
                throw new Error('Failed to create anonymous account');
            }
            
            userData = await response.json();
            showModal = true;
        } catch (err) {
            // @ts-ignore
            error = err.message;
            setTimeout(() => error = null, 3000);
        }
    }

    /**
     * Copies text to clipboard and shows feedback
     * @param {string} text - Text to copy
     * @param {string} field - Field name ('username' or 'password')
     * @returns {void}
     */
    function copyToClipboard(text, field) {
        navigator.clipboard.writeText(text).then(() => {
            // @ts-ignore
            copied[field] = true;
            // @ts-ignore
            setTimeout(() => copied[field] = false, 2000);
        });
    }
</script>

<NavBar/>

<div class="w-full sm:w-[400px] md:w-[510px] mx-auto flex flex-col items-center justify-center mt-4 sm:mt-6 md:mt-8 bg-gray-400 text-red-900 px-2 sm:px-4">
    <div class="flex flex-col items-center justify-center py-2">
        <span class="text-lg sm:text-xl md:text-2xl">🔔</span>
        <p class="text-base sm:text-lg md:text-xl font-medium">Disclaimer</p>
    </div>
    <p class="gap-2 text-sm sm:text-base md:text-md p-4 sm:p-4 md:p-6 text-center sm:text-left">
        All content presented on this site is the sole responsibility of the respective authors and contributors. As the creator of this platform, I do not hold liability for any content published herein. Viewers are advised to exercise their own discretion while browsing.
        
        If you come across any material you believe should be reviewed or removed, please contact the moderators, the original author, or reach out to me directly.
    </p>
    
    <div class="flex flex-col gap-2 sm:gap-3 md:gap-4 justify-center mt-4 sm:mt-5 md:mt-6 pb-4 sm:pb-5 md:pb-6">
        <!-- Use Anonymously Button -->
        <button 
            on:click={createRandomAccount}
            class="p-3 sm:p-3 md:p-4 m-1 sm:m-1 md:m-2 flex flex-col bg-black text-white shadow-md border-2 border-gray-600 hover:bg-gray-800 transition duration-300 ease-in-out items-center gap-1 sm:gap-2 w-[120px] sm:w-[135px] md:w-[150px] transform hover:scale-105 active:scale-95">
            <span class="text-base sm:text-lg md:text-xl">🔒</span> 
            <span class="text-xs sm:text-sm md:text-base">Use Anonymously</span>
        </button>
    
        <!-- Login Button -->
        <button on:click={redirectToLink} class="p-3 sm:p-3 md:p-4 m-1 sm:m-1 md:m-2 bg-red-500 flex flex-col text-white shadow-md border-2 border-red-700 hover:bg-red-600 transition duration-300 ease-in-out items-center gap-1 sm:gap-2 w-[120px] sm:w-[135px] md:w-[150px] transform hover:scale-105 active:scale-95">
            <span class="text-base sm:text-lg md:text-xl">🔑</span>
            <span class="text-xs sm:text-sm md:text-base">Login</span>
        </button>
    </div>
</div>

{#if error}
    <div class="fixed top-4 right-4 bg-red-500 text-white px-4 py-2 rounded shadow-lg animate-fade-in">
        {error}
    </div>
{/if}

{#if showModal && userData}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-gray-800 text-white p-6 rounded-lg max-w-md w-full mx-4 border-2 border-red-500">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold">Anonymous Account Created</h2>
                <button 
                    on:click={() => showModal = false}
                    class="text-gray-400 hover:text-white">
                    ✕
                </button>
            </div>
            
            <div class="space-y-4">
                <div>
                    <div class="flex justify-between items-center mb-1">
                        <p class="text-gray-300">Username:</p>
                        <button 
                            on:click={() => copyToClipboard(userData.username, 'username')}
                            class="text-xs bg-gray-700 hover:bg-gray-600 px-2 py-1 rounded">
                            {copied.username ? 'Copied!' : 'Copy'}
                        </button>
                    </div>
                    <div class="flex items-center bg-gray-900 p-2 rounded">
                        <p class="font-mono flex-grow">{userData.username}</p>
                    </div>
                </div>
                
                <div>
                    <div class="flex justify-between items-center mb-1">
                        <p class="text-gray-300">Password:</p>
                        <button 
                            on:click={() => copyToClipboard(userData.password, 'password')}
                            class="text-xs bg-gray-700 hover:bg-gray-600 px-2 py-1 rounded">
                            {copied.password ? 'Copied!' : 'Copy'}
                        </button>
                    </div>
                    <div class="flex items-center bg-gray-900 p-2 rounded">
                        <p class="font-mono flex-grow">{userData.password}</p>
                    </div>
                </div>
                
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <p class="text-gray-300">Created:</p>
                        <p>{new Date().toLocaleString()}</p>
                    </div>
                    <div>
                        <p class="text-gray-300">Expires:</p>
                        <p>{new Date(userData.expires_at).toLocaleString()}</p>
                    </div>
                </div>
                
                <div class="pt-4">
                    <p class="text-sm text-gray-400">
                        These credentials have been copied to your clipboard.
                        Use them to login on the next page.
                    </p>
                </div>
            </div>
            
            <div class="mt-6 flex justify-end">
                <button 
                    on:click={() => {
                        copyToClipboard(`${userData.username}\n${userData.password}`, 'username');
                        redirectToLink();
                    }}
                    class="bg-red-600 hover:bg-red-700 px-4 py-2 rounded transition">
                    Proceed to Login
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .animate-fade-in {
        animation: fadeIn 0.3s ease-in-out;
    }
    
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(-10px); }
        to { opacity: 1; transform: translateY(0); }
    }
</style>