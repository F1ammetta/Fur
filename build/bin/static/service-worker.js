     // Cache names
    const CACHE_NAME = 'xfs-cache-v1';
    const urlsToCache = [
      '/',
      '/static/htmx.js',
    ];

    // Install event
    self.addEventListener('install', event => {
      event.waitUntil(
        caches.open(CACHE_NAME)
          .then(cache => {
            return cache.addAll(urlsToCache);
          })
      );
    });

    // Activate event
    self.addEventListener('activate', event => {
      event.waitUntil(
        caches.keys()
          .then(cacheNames => {
            return Promise.all(
              cacheNames.filter(cacheName => {
                return cacheName !== CACHE_NAME;
              }).map(cacheName => {
                return caches.delete(cacheName);
              })
            );
          })
      );
    });

    // Fetch event
    self.addEventListener('fetch', event => {
      event.respondWith(
        caches.match(event.request)
          .then(response => {
            // Cache hit - return the response
            if (response) {
              return response;
            }

            // Clone the request
            const fetchRequest = event.request.clone();

            return fetch(fetchRequest)
              .then(response => {
                // Check if we received a valid response
                if (!response || response.status !== 200 || response.type !== 'basic') {
                  return response;
                }

                // Clone the response
                const responseToCache = response.clone();

                caches.open(CACHE_NAME)
                  .then(cache => {
                    cache.put(event.request, responseToCache);
                  });

                return response;
              })
              .catch(() => {
                // Offline fallback page
                return caches.match('/offline.html');
              });
          })
      );
    });
