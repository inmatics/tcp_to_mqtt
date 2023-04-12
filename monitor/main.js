// monitor/main.js
const map = L.map('map').setView([51.505, -0.09], 13);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    maxZoom: 19,
    tileSize: 512,
    zoomOffset: -1
}).addTo(map);

const clientId = 'mqttjs_' + Math.random().toString(16).substr(2, 8);
const brokerUrl = 'wss://mqtt.inmatics.io:443'; // Replace with your MQTT broker WebSocket URL and port
const topic = 'devices/new'; // Replace with the topic you want to subscribe to

const client = mqtt.connect(brokerUrl, { clientId });

client.on('connect', () => {
    console.log('Connected to MQTT broker');
    client.subscribe(topic, (err) => {
        if (err) {
            console.error('Failed to subscribe to topic:', err);
        } else {
            console.log(`Subscribed to topic: ${topic}`);
        }
    });
});


// Create an object to store markers by IMEI
const markersByIMEI = {};

client.on('message', (topic, message) => {

    const data = JSON.parse(message);
    const {lat, lng, Imei, speed} = data

    if (!isNaN(lat) && !isNaN(lng)) {

        if (markersByIMEI.hasOwnProperty(Imei)) {
            markersByIMEI[Imei].setLatLng([lat, lng]);
            markersByIMEI[Imei].bindPopup(`IMEI: ${Imei}<br>Latitude: ${lat}<br>Longitude: ${lng}`).openPopup();
            setTimeout(() => {
                markersByIMEI[Imei].closePopup();
            }, 5000);
        } else {
            const marker = L.marker([lat, lng]).addTo(map)
                .bindPopup(`IMEI: ${Imei}<br>Latitude: ${lat}<br>Longitude: ${lng}<br> Speed: ${speed}`)
                .openPopup();

            setTimeout(() => {
                marker.closePopup();
            }, 5000);

            markersByIMEI[Imei] = marker;
        }
    } else {
        console.error('Invalid latitude and longitude in received message');
    }
});

