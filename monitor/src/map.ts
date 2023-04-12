import L from 'leaflet';
import {Record} from "./types";


export interface MarkerManager {
    addOrUpdateMarker: (record: Record, map: L.Map) => void;
}

export function initializeMap(): L.Map {
    const map = L.map('map').setView([-34.61, -58.4], 13);

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

    function addOrUpdateMarker(record: Record, map: L.Map): void {
        const {Imei, lat, lng} = record

        const content = `<div>
                                    IMEI: ${Imei}<br>
                                    Speed: ${record.speed} km/h<br>
                                    Battery: ${record.battery}<br>
                                    Ignition: ${record.ignition}<br>
                                    </div>`;
        if (markersByIMEI.hasOwnProperty(Imei)) {
            markersByIMEI[Imei].setLatLng([lat, lng]);
            markersByIMEI[Imei].bindPopup(content);

        } else {
            markersByIMEI[Imei] = L.marker([lat, lng]).addTo(map)
                .bindPopup(content);
        }
    }

    return {
        addOrUpdateMarker
    };
}
