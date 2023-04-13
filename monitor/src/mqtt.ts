import {MqttClient} from 'mqtt';
import {UIManager} from './ui';
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

export function handleMessage(ui: UIManager, topic: string, message: Buffer, map: L.Map): void {
    const data : Record = JSON.parse(message.toString());

    if (!isNaN(data.lat) && !isNaN(data.lng)) {
        ui.update(data, map);
    } else {
        console.error('Invalid latitude and longitude in received message');
    }
}
