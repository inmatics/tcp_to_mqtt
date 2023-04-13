import {MqttClient} from 'mqtt';
import {createMarkerManager, MarkerManager} from './map';
import {Record} from "./types";

export function subscribeToTopic(client: MqttClient, topic: string): void {
    client.subscribe(topic, (err) => {
        if (err) {
            console.error('Failed to subscribe to topic:', err);
        } else {
            console.log(`Subscribed to topic: ${topic}`);
        }
    });
}

export function handleMessage( topic: string, message: Buffer, map: L.Map): void {
    const markerManager = createMarkerManager()
    const data : Record = JSON.parse(message.toString());

    if (!isNaN(data.lat) && !isNaN(data.lng)) {
        markerManager.addOrUpdateMarker(data, map);
    } else {
        console.error('Invalid latitude and longitude in received message');
    }
}