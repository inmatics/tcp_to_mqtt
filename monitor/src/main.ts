import {createMarkerManager, initializeMap} from './map';
import {handleMessage, subscribeToTopic} from './mqtt';
import * as mqtt from "mqtt";

const map = initializeMap();
const client = mqtt.connect('wss://mqtt.inmatics.io:443', {clientId: 'mqttjs_' + Math.random().toString(16).substr(2, 8)})

client.on('connect', () => {
    subscribeToTopic(client, 'devices/new');
});

const markerManager = createMarkerManager();
client.on('message', (topic: string, message: Buffer) => {
    handleMessage(markerManager, topic, message, map);
});

document.addEventListener("DOMContentLoaded", function() {
    const elementId = "sidebar";

    const toggleSidebar = () => {
        const sidebar = document.getElementById(elementId);
        if (sidebar) {
            sidebar.classList.toggle("collapsed");
        }
    }

    const toggleButton = document.querySelector(".toggle-sidebar");
    if (toggleButton) {
        toggleButton.addEventListener("click", toggleSidebar);
    }
});
