<svelte:head>
    <title>Login - 5chan</title>
</svelte:head>

<script>
    import { onMount } from 'svelte';
    import { browser } from '$app/environment';
    
    /** @type {string} */
    let username = '';
    
    /** @type {string} */
    let password = '';
    
    /** @type {string} */
    let error = '';
    
    /** @type {boolean} */
    let isLoading = false;
    
    //@ts-ignore
    /** @type {import('$env/dynamic/public').PublicEnv} */
    import { env } from '$env/dynamic/public';
    
    /** @type {string} Backend API base URL */
    const backendAPI = env.PUBLIC_BACKEND_API || '';

    /**
     * Sets a cookie with JWT token
     * @param {string} name - Cookie name
     * @param {string} value - Cookie value
     * @param {number} days - Expiration in days
     * @returns {void}
     */
     function setCookie(name, value, days) {
        if (!browser) return;
        
        const date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        const expires = "expires=" + date.toUTCString();
        
        // Using exact requested format
        // @ts-ignore - document.cookie is always available in browsers
        document.cookie = `${name}=${value}; ${expires}; path=/`;
    }

    /**
     * Gets cookie value by name
     * @param {string} name - Cookie name to find
     * @returns {string|null} Cookie value or null if not found
     */
     function getCookie(name) {
        if (!browser) return null;
        
        const cookies = document.cookie.split(';');
        for (const cookie of cookies) {
            const [cookieName, cookieValue] = cookie.trim().split('=');
            if (cookieName === name) {
                return decodeURIComponent(cookieValue);
            }
        }
        return null;
    }

    /**
     * Handles login form submission
     * @param {Event} e - Form submit event
     * @returns {Promise<void>}
     */
    async function handleLogin(e) {
        e.preventDefault();
        isLoading = true;
        error = '';
        
        try {
            /** @type {Response} */
            const response = await fetch(`${backendAPI}/login`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (!response.ok) {
                /** @type {Record<string, unknown>} */
                const data = await response.json();
                throw new Error(
                    // @ts-ignore - We know error might be string
                    typeof data.error === 'string' ? data.error : 'Login failed'
                );
            }

            /** @type {{ token: string }} */
            const data = await response.json();
            
            // Determine cookie duration based on username
            const isAnonymous = username.startsWith('anon_');
            const cookieDuration = isAnonymous ? 30 : 90; // 30 days for anonymous, 90 for regular
            
            // Store token in cookie
            setCookie('jwt_token', data.token, cookieDuration);
            
            // Redirect to home
            window.location.href = '/';
            
        } catch (err) {
            // @ts-ignore - We know err will have message property
            error = err.message;
            console.error('Login error:', err);
        } finally {
            isLoading = false;
        }
    }

      /**
     * Checks for existing JWT and redirects if found
     * @returns {void}
     */
     function checkAuthAndRedirect() {
        if (!browser) return;
        
        const token = getCookie('jwt_token');
        if (token) {
            window.location.href = '/';
        }
    }

    /**
     * Initialize page on mount
     * @returns {void}
     */
    onMount(() => {
        if (browser) {
            // @ts-ignore - document is available in browser
            document.body.className = 'bg-gray-900';
            checkAuthAndRedirect();
        }
    });
</script>

<div class="min-h-screen flex items-center justify-center p-4">
    <div class="bg-gray-800 rounded-lg p-8 w-full max-w-md shadow-lg">
        <div class="text-center mb-8">
            <h1 class="text-2xl font-bold text-red-500 mb-2">5chan</h1>
            <p class="text-gray-400">Successor to the true chaos</p>
        </div>

        {#if error}
            <div class="mb-4 p-3 bg-red-900 text-red-200 rounded-md text-sm">
                {error}
            </div>
        {/if}

        <form on:submit={handleLogin} class="space-y-6">
            <div>
                <label for="username" class="block text-sm font-medium text-gray-300 mb-1">
                    Username
                </label>
                <input
                    type="text"
                    id="username"
                    bind:value={username}
                    required
                    class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:border-red-500 transition"
                    placeholder="Enter your username"
                />
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-300 mb-1">
                    Password
                </label>
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    required
                    class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:border-red-500 transition"
                    placeholder="••••••••"
                />
            </div>

            <button
                type="submit"
                disabled={isLoading}
                class="w-full bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md transition duration-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {isLoading ? 'Logging in...' : 'Login'}
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-400">
            Don't have an account?{' '}
            <a href="/" class="text-red-400 hover:text-red-300 transition">
                Use anonymously
            </a>
        </div>
    </div>
</div>

<style>
    .transition {
        transition-property: background-color, border-color, color, fill, stroke, opacity, box-shadow, transform;
        transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
        transition-duration: 150ms;
    }
</style>