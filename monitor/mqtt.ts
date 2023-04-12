import * as mqtt from 'mqtt';
import { MqttClient } from 'mqtt';
import { MarkerManager } from './map';


export function initializeMQTTClient(brokerUrl: string, clientId: string): MqttClient {
    const client = mqtt.connect(brokerUrl, { clientId });
    return client;
}

export function subscribeToTopic(client: MqttClient, topic: string): void {
    client.subscribe(topic, (err) => {
        if (err) {
            console.error('Failed to subscribe to topic:', err);
        } else {
            console.log(`Subscribed to topic: ${topic}`);
        }
    });
}

export function handleMessage(markerManager: MarkerManager, topic: string, message: Buffer, map: L.Map): void {
    const data = JSON.parse(message.toString());
    const { lat, lng, Imei } = data;

    if (!isNaN(lat) && !isNaN(lng)) {
        markerManager.addOrUpdateMarker(Imei, lat, lng, map);
    } else {
        console.error('Invalid latitude and longitude in received message');
    }
}
