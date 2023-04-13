import L from 'leaflet';
import {Record} from "./types";

export interface UIManager {
    update: (record: Record, map: L.Map) => void;
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

export function createUIManager(): UIManager {
    const markersByIMEI: { [imei: string]: L.Marker } = {};
    const recordsByImei: { [imei: string]: Record } = {};

    function updateImeiList(record: Record) {
        recordsByImei[record.imei] = record

        const itemList = document.getElementById("item-list");
        if (itemList) {
            Object.values(recordsByImei).forEach((item) => {
                let li;
                li = document.getElementById(item.imei);
                if (!li) {
                    li = document.createElement("li");
                    li.id = item.imei
                }

                const a = document.createElement("a");
                li.textContent = item.imei;
                a.href = "#";
                li.appendChild(a);
                itemList.appendChild(li);
            });
        }
    }

    function addOrUpdateMarker(record: Record, map: L.Map): void {
        const {imei, lat, lng} = record

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

    function update(record: Record, map: L.Map): void {
        addOrUpdateMarker(record, map)
        updateImeiList(record);
    }

    return {
        update
    };
}
