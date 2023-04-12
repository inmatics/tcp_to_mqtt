import { initializeMap, createMarkerManager } from './map';
import { initializeMQTTClient, subscribeToTopic, handleMessage } from './mqtt';


const map = initializeMap();
const clientId = 'mqttjs_' + Math.random().toString(16).substr(2, 8);
const brokerUrl = 'wss://mqtt.inmatics.io:443'; // Replace with your MQTT broker WebSocket URL and port
const topic = 'devices/new'; // Replace with the topic you want to subscribe to
const markerManager = createMarkerManager();

const client = initializeMQTTClient(brokerUrl, clientId);

client.on('connect', () => {
    subscribeToTopic(client, topic);
});

client.on('message', (topic: string, message: Buffer) => {
    handleMessage(markerManager, topic, message, map);
});
