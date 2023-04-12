import {addOrUpdateMarker} from './map.js';

// Function to initialize the MQTT client
export function initializeMQTTClient(brokerUrl, clientId) {
    return mqtt.connect(brokerUrl, {clientId});
}

// Function to subscribe to an MQTT topic
export function subscribeToTopic(client, topic) {
    client.subscribe(topic, (err) => {
        if (err) {
            console.error('Failed to subscribe to topic:', err);
        } else {
            console.log(`Subscribed to topic: ${topic}`);
        }
    });
}

// Function to handle messages received from the MQTT broker
export function handleMessage(markersByIMEI, topic, message, map) {
    const data = JSON.parse(message);
    const { lat, lng, Imei } = data;

    if (!isNaN(lat) && !isNaN(lng)) {
        addOrUpdateMarker(markersByIMEI, Imei, lat, lng, map);
    } else {
        console.error('Invalid latitude and longitude in received message');
    }
}
