import L from 'leaflet';

export interface MarkerManager {
    addOrUpdateMarker: (imei: string, lat: number, lng: number, map: L.Map) => void;
}

export function initializeMap(): L.Map {
    const map = L.map('map').setView([51.505, -0.09], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
        tileSize: 512,
        zoomOffset: -1
    }).addTo(map);

    return map;
}

export function createMarkerManager(): MarkerManager {
    const markersByIMEI: { [imei: string]: L.Marker } = {};

    function addOrUpdateMarker(imei: string, lat: number, lng: number, map: L.Map): void {
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

    return {
        addOrUpdateMarker
    };
}
