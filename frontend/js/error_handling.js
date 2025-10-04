// PWA Error Handling System
class PWAErrorHandler {
    constructor() {
        this.cache = new PWACacheHandler();
        this.installPrompt = new PWAInstallPromptHandler();
        this.sse = new SSEErrorHandler();
        this.offlineQueue = [];
        this.isOnline = navigator.onLine;

        this.initializeEventListeners();
    }

    initializeEventListeners() {
        // Network status monitoring
        window.addEventListener('online', () => {
            this.handleOnlineStatus(true);
        });

        window.addEventListener('offline', () => {
            this.handleOnlineStatus(false);
        });

        // Service Worker error handling
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.addEventListener('message', (event) => {
                this.handleServiceWorkerMessage(event);
            });
        }

        // Unhandled promise rejections
        window.addEventListener('unhandledrejection', (event) => {
            this.handleUnhandledRejection(event);
        });

        // Global error handler
        window.addEventListener('error', (event) => {
            this.handleGlobalError(event);
        });
    }

    handleOnlineStatus(isOnline) {
        this.isOnline = isOnline;
        const event = new CustomEvent('pwa:networkstatus', {
            detail: { isOnline, timestamp: Date.now() }
        });
        window.dispatchEvent(event);

        if (isOnline) {
            this.processOfflineQueue();
        }
    }

    handleServiceWorkerMessage(event) {
        const { type, payload } = event.data;

        switch (type) {
            case 'CACHE_ERROR':
                this.handleCacheError(payload);
                break;
            case 'SYNC_ERROR':
                this.handleSyncError(payload);
                break;
            case 'PERFORMANCE_ISSUE':
                this.handlePerformanceIssue(payload);
                break;
        }
    }

    handleUnhandledRejection(event) {
        const error = {
            type: 'unhandled_promise_rejection',
            message: event.reason?.message || 'Unhandled promise rejection',
            stack: event.reason?.stack,
            timestamp: Date.now(),
            url: window.location.href,
            userAgent: navigator.userAgent
        };

        this.logError(error);
        this.showUserNotification('An unexpected error occurred. Please refresh the page.', 'error');
    }

    handleGlobalError(event) {
        const error = {
            type: 'javascript_error',
            message: event.message,
            filename: event.filename,
            lineno: event.lineno,
            colno: event.colno,
            stack: event.error?.stack,
            timestamp: Date.now(),
            url: window.location.href
        };

        this.logError(error);
    }

    async handleCacheError(payload) {
        console.error('Cache error:', payload);

        // Attempt cache recovery
        try {
            await this.cache.recoverFromError(payload);
            this.showUserNotification('Cache recovered successfully', 'success');
        } catch (error) {
            this.showUserNotification('Cache error occurred. Some features may not work properly.', 'warning');
        }
    }

    handleSyncError(payload) {
        console.error('Sync error:', payload);

        // Add to offline queue for retry
        this.offlineQueue.push({
            type: 'sync',
            payload,
            timestamp: Date.now()
        });
    }

    handlePerformanceIssue(payload) {
        console.warn('Performance issue detected:', payload);

        // Log performance metrics
        this.logPerformanceMetric(payload);
    }

    async processOfflineQueue() {
        if (this.offlineQueue.length === 0) return;

        console.log(`Processing ${this.offlineQueue.length} offline actions...`);

        for (const item of this.offlineQueue) {
            try {
                await this.retryOfflineAction(item);
            } catch (error) {
                console.error('Failed to retry offline action:', error);
            }
        }

        this.offlineQueue = [];
    }

    async retryOfflineAction(item) {
        // Implement retry logic based on action type
        switch (item.type) {
            case 'sync':
                // Retry sync operation
                break;
            default:
                console.warn('Unknown offline action type:', item.type);
        }
    }

    logError(error) {
        // In production, send to error tracking service
        console.error('PWA Error:', error);

        // Store in local storage for debugging
        this.storeErrorLocally(error);
    }

    logPerformanceMetric(metric) {
        console.log('Performance Metric:', metric);

        // Store performance data
        this.storePerformanceMetricLocally(metric);
    }

    storeErrorLocally(error) {
        try {
            const errors = JSON.parse(localStorage.getItem('pwa_errors') || '[]');
            errors.push(error);

            // Keep only last 50 errors
            if (errors.length > 50) {
                errors.splice(0, errors.length - 50);
            }

            localStorage.setItem('pwa_errors', JSON.stringify(errors));
        } catch (e) {
            console.error('Failed to store error locally:', e);
        }
    }

    storePerformanceMetricLocally(metric) {
        try {
            const metrics = JSON.parse(localStorage.getItem('pwa_performance') || '[]');
            metrics.push(metric);

            // Keep only last 100 metrics
            if (metrics.length > 100) {
                metrics.splice(0, metrics.length - 100);
            }

            localStorage.setItem('pwa_performance', JSON.stringify(metrics));
        } catch (e) {
            console.error('Failed to store performance metric locally:', e);
        }
    }

    showUserNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `pwa-notification pwa-notification--${type}`;
        notification.innerHTML = `
            <div class="pwa-notification__content">
                <span class="pwa-notification__message">${message}</span>
                <button class="pwa-notification__close">&times;</button>
            </div>
        `;

        // Add to DOM
        document.body.appendChild(notification);

        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 5000);

        // Manual close handler
        const closeBtn = notification.querySelector('.pwa-notification__close');
        closeBtn.addEventListener('click', () => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        });
    }

    getStoredErrors() {
        try {
            return JSON.parse(localStorage.getItem('pwa_errors') || '[]');
        } catch (e) {
            return [];
        }
    }

    getStoredPerformanceMetrics() {
        try {
            return JSON.parse(localStorage.getItem('pwa_performance') || '[]');
        } catch (e) {
            return [];
        }
    }

    clearStoredData() {
        localStorage.removeItem('pwa_errors');
        localStorage.removeItem('pwa_performance');
        this.offlineQueue = [];
    }
}

// PWA Cache Error Handler
class PWACacheHandler {
    constructor() {
        this.cacheName = 'pwa-cache-v1';
    }

    async recoverFromError(payload) {
        switch (payload.error) {
            case 'CACHE_FULL':
                await this.cleanOldCache();
                break;
            case 'CACHE_CORRUPTED':
                await this.rebuildCache();
                break;
            case 'NETWORK_ERROR':
                await this.handleNetworkError();
                break;
            default:
                throw new Error(`Unknown cache error: ${payload.error}`);
        }
    }

    async cleanOldCache() {
        try {
            const cacheNames = await caches.keys();
            const oldCaches = cacheNames.filter(name => name !== this.cacheName);

            await Promise.all(
                oldCaches.map(cacheName => caches.delete(cacheName))
            );
        } catch (error) {
            throw new Error(`Failed to clean old cache: ${error.message}`);
        }
    }

    async rebuildCache() {
        try {
            await caches.delete(this.cacheName);
            // Rebuild cache with fresh resources
            await this.initializeCache();
        } catch (error) {
            throw new Error(`Failed to rebuild cache: ${error.message}`);
        }
    }

    async handleNetworkError() {
        // Switch to cache-first strategy when network fails
        console.log('Switching to cache-first strategy due to network error');
    }

    async initializeCache() {
        // Initialize cache with essential resources
        const essentialResources = [
            '/',
            '/index.html',
            '/css/styles.css',
            '/js/app.js'
        ];

        try {
            const cache = await caches.open(this.cacheName);
            await cache.addAll(essentialResources);
        } catch (error) {
            console.error('Failed to initialize cache:', error);
        }
    }
}

// PWA Install Prompt Error Handler
class PWAInstallPromptHandler {
    constructor() {
        this.deferredPrompt = null;
        this.installButton = null;
        this.isInstalled = false;
    }

    initialize() {
        // Listen for install prompt
        window.addEventListener('beforeinstallprompt', (event) => {
            event.preventDefault();
            this.deferredPrompt = event;
            this.showInstallButton();
        });

        // Check if already installed
        if (window.matchMedia('(display-mode: standalone)').matches) {
            this.isInstalled = true;
            this.hideInstallButton();
        }
    }

    showInstallButton() {
        let button = document.getElementById('pwa-install-button');
        if (!button) {
            button = document.createElement('button');
            button.id = 'pwa-install-button';
            button.className = 'pwa-install-button';
            button.innerHTML = '📱 Install App';
            button.addEventListener('click', () => this.handleInstallClick());
            document.body.appendChild(button);
        }
        button.style.display = 'block';
        this.installButton = button;
    }

    hideInstallButton() {
        const button = document.getElementById('pwa-install-button');
        if (button) {
            button.style.display = 'none';
        }
    }

    async handleInstallClick() {
        if (!this.deferredPrompt) {
            this.showInstallError('Install prompt not available');
            return;
        }

        try {
            this.deferredPrompt.prompt();
            const { outcome } = await this.deferredPrompt.userChoice;

            if (outcome === 'accepted') {
                this.isInstalled = true;
                this.hideInstallButton();
                this.showUserNotification('App installed successfully!', 'success');
            }

            this.deferredPrompt = null;
        } catch (error) {
            this.showInstallError(`Installation failed: ${error.message}`);
        }
    }

    showInstallError(message) {
        this.showUserNotification(message, 'error');
    }

    showUserNotification(message, type) {
        // Use the global error handler's notification system
        if (window.pwaErrorHandler) {
            window.pwaErrorHandler.showUserNotification(message, type);
        }
    }
}

// SSE Error Handler for PWA
class SSEErrorHandler {
    constructor() {
        this.eventSource = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000; // Start with 1 second
        this.heartbeatInterval = 30000; // 30 seconds
        this.isConnected = false;
    }

    connect(url) {
        if (this.eventSource) {
            this.disconnect();
        }

        try {
            this.eventSource = new EventSource(url);
            this.setupEventListeners();
        } catch (error) {
            this.handleConnectionError(error);
        }
    }

    disconnect() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
        this.isConnected = false;
        this.reconnectAttempts = 0;
    }

    setupEventListeners() {
        this.eventSource.onopen = () => {
            this.handleConnectionOpen();
        };

        this.eventSource.onmessage = (event) => {
            this.handleMessage(event);
        };

        this.eventSource.onerror = (event) => {
            this.handleConnectionError(event);
        };

        // Setup heartbeat
        this.setupHeartbeat();
    }

    handleConnectionOpen() {
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;

        const event = new CustomEvent('pwa:sse:connected', {
            detail: { timestamp: Date.now() }
        });
        window.dispatchEvent(event);
    }

    handleMessage(event) {
        try {
            const data = JSON.parse(event.data);
            const messageEvent = new CustomEvent('pwa:sse:message', {
                detail: { data, timestamp: Date.now() }
            });
            window.dispatchEvent(messageEvent);
        } catch (error) {
            console.error('Failed to parse SSE message:', error);
        }
    }

    handleConnectionError(error) {
        this.isConnected = false;

        const event = new CustomEvent('pwa:sse:error', {
            detail: {
                error,
                attempt: this.reconnectAttempts + 1,
                timestamp: Date.now()
            }
        });
        window.dispatchEvent(event);

        // Attempt reconnection with exponential backoff
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnection();
        } else {
            this.handleMaxReconnectionAttempts();
        }
    }

    scheduleReconnection() {
        this.reconnectAttempts++;
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000); // Max 30 seconds

        setTimeout(() => {
            if (!this.isConnected) {
                console.log(`Attempting SSE reconnection ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                this.connect(this.eventSource?.url || '');
            }
        }, this.reconnectDelay);
    }

    handleMaxReconnectionAttempts() {
        const errorHandler = window.pwaErrorHandler;
        if (errorHandler) {
            errorHandler.showUserNotification(
                'Real-time updates are currently unavailable. Please refresh the page.',
                'warning'
            );
        }
    }

    setupHeartbeat() {
        setInterval(() => {
            if (this.isConnected) {
                // Send heartbeat comment to keep connection alive
                // This is handled by the browser's EventSource implementation
            }
        }, this.heartbeatInterval);
    }

    sendMessage(message) {
        if (!this.isConnected) {
            throw new Error('SSE not connected');
        }

        // For sending messages, you might need a different approach
        // This is a placeholder for message sending functionality
        console.log('Sending SSE message:', message);
    }
}

// PWA Service Worker Error Handler
class PWAServiceWorkerHandler {
    constructor() {
        this.serviceWorker = null;
        this.registration = null;
    }

    async register() {
        if (!('serviceWorker' in navigator)) {
            throw new Error('Service Worker not supported');
        }

        try {
            this.registration = await navigator.serviceWorker.register('/sw.js');
            this.serviceWorker = this.registration.active;

            this.setupServiceWorkerListeners();
        } catch (error) {
            throw new Error(`Service Worker registration failed: ${error.message}`);
        }
    }

    setupServiceWorkerListeners() {
        navigator.serviceWorker.addEventListener('controllerchange', () => {
            window.location.reload();
        });

        navigator.serviceWorker.addEventListener('message', (event) => {
            this.handleServiceWorkerMessage(event);
        });
    }

    handleServiceWorkerMessage(event) {
        const { type, payload } = event.data;

        switch (type) {
            case 'SKIP_WAITING':
                this.registration.waiting?.postMessage({ type: 'SKIP_WAITING' });
                break;
            case 'CACHE_UPDATED':
                console.log('Cache updated:', payload);
                break;
            case 'BACKGROUND_SYNC':
                console.log('Background sync completed:', payload);
                break;
        }
    }

    async updateServiceWorker() {
        if (!this.registration) {
            throw new Error('Service Worker not registered');
        }

        try {
            await this.registration.update();
        } catch (error) {
            throw new Error(`Service Worker update failed: ${error.message}`);
        }
    }
}

// Initialize PWA Error Handler
if (typeof window !== 'undefined') {
    window.pwaErrorHandler = new PWAErrorHandler();

    // Initialize PWA components
    document.addEventListener('DOMContentLoaded', () => {
        // Initialize install prompt handler
        const installPrompt = new PWAInstallPromptHandler();
        installPrompt.initialize();

        // Initialize service worker
        const serviceWorker = new PWAServiceWorkerHandler();
        serviceWorker.register().catch(error => {
            console.error('Service Worker registration failed:', error);
        });

        // Make components globally available
        window.pwaInstallPrompt = installPrompt;
        window.pwaServiceWorker = serviceWorker;
    });
}

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        PWAErrorHandler,
        PWACacheHandler,
        PWAInstallPromptHandler,
        SSEErrorHandler,
        PWAServiceWorkerHandler
    };
}