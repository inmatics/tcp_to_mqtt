// Function to initialize the Leaflet map
export function initializeMap() {
    const map = L.map('map').setView([-34, -58], 8);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
        tileSize: 512,
        zoomOffset: -1
    }).addTo(map);

    return map;
}

// Function to add or update a marker on the map
export function addOrUpdateMarker(markersByIMEI, imei, lat, lng, map) {
    if (markersByIMEI.hasOwnProperty(imei)) {
        markersByIMEI[imei].setLatLng([lat, lng]);
        markersByIMEI[imei].bindPopup(`IMEI: ${imei}<br>Latitude: ${lat}<br>Longitude: ${lng}`).openPopup();
        setTimeout(() => {
            markersByIMEI[imei].closePopup();
        }, 2000);
    } else {
        const marker = L.marker([lat, lng]).addTo(map)
            .bindPopup(`IMEI: ${imei}<br>Latitude: ${lat}<br>Longitude: ${lng}`)
            .openPopup();

        setTimeout(() => {
            marker.closePopup();
        }, 2000);

        markersByIMEI[imei] = marker;
    }
}
