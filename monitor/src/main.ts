import { initializeMap, createMarkerManager } from './map';
import { subscribeToTopic, handleMessage } from './mqtt';
import * as mqtt from "mqtt";

const map = initializeMap();
const clientId = 'mqttjs_' + Math.random().toString(16).substr(2, 8);
const brokerUrl = 'wss://mqtt.example.com:443'; // Replace with your MQTT broker WebSocket URL and port
const topic = 'devices/new'; // Replace with the topic you want to subscribe to
const markerManager = createMarkerManager();

const client = mqtt.connect(brokerUrl, {clientId})

client.on('connect', () => {
    subscribeToTopic(client, topic);
});

client.on('message', (topic: string, message: Buffer) => {
    handleMessage(markerManager, topic, message, map);
});
