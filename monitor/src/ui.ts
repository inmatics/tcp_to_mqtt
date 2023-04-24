// Import Leaflet library and the custom Record type
import L from 'leaflet';
import { Record } from "./types";

// Define the UIManager interface
export interface UIManager {
    update: (record: Record, map: L.Map) => void;
}

/**
 * Initialize a Leaflet map with OpenStreetMap tiles
 * @returns A new Leaflet map instance
 */
export function initializeMap(): L.Map {
    // Create a map and set the initial view to specific coordinates and zoom level
    const map = L.map('map').setView([-34.61, -58.4], 13);

    // Add OpenStreetMap tile layer to the map
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
        tileSize: 512,
        zoomOffset: -1
    }).addTo(map);

    return map;
}

/**
 * Create a new UIManager to manage map markers and IMEI list
 * @returns A new UIManager instance
 */
export function createUIManager(): UIManager {
    // Create empty objects to store markers and records by IMEI
    const markersByIMEI: { [imei: string]: L.Marker } = {};
    const recordsByImei: { [imei: string]: Record } = {};

    /**
     * Update the IMEI list in the DOM with the given record
     * @param record The record to be added or updated in the IMEI list
     */
    function updateImeiList(record: Record) {
        recordsByImei[record.imei] = record;

        const itemList = document.getElementById("item-list");
        if (itemList) {
            Object.values(recordsByImei).forEach((item) => {
                let li;
                li = document.getElementById(item.imei);
                if (!li) {
                    li = document.createElement("li");
                    li.id = item.imei;
                }

                const a = document.createElement("a");
                li.textContent = item.imei;
                a.href = "#";
                li.appendChild(a);
                itemList.appendChild(li);
            });
        }
    }

    /**
     * Add or update a map marker for the given record
     * @param record The record containing the data to update the marker
     * @param map The map instance to add or update the marker on
     */
    function addOrUpdateMarker(record: Record, map: L.Map): void {
        const { imei, lat, lng } = record;

        const content = `<div>
            IMEI: ${imei}<br>
            Speed: ${record.speed} km/h<br>
            Battery: ${record.battery}<br>
            Ignition: ${record.ignition}<br>
        </div>`;

        if (markersByIMEI.hasOwnProperty(imei)) {
            markersByIMEI[imei].setLatLng([lat, lng]);
            markersByIMEI[imei].bindPopup(content);

        } else {
            markersByIMEI[imei] = L.marker([lat, lng]).addTo(map)
                .bindPopup(content);
        }
    }

    /**
     * Update the map marker and IMEI list with the given record
     * @param record The record to update the marker and IMEI list with
     * @param map The map instance to update the marker on
     */
    function update(record: Record, map: L.Map): void {
        addOrUpdateMarker(record, map)
        updateImeiList(record);
    }

    return {
        update
    };
}
