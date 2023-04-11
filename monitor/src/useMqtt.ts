import mqtt from "mqtt";
import { useEffect, useState } from "react";

type Record = {
    Angle: number;
    lng: number;
    odometer: number;
    raw_message: string;
    Imei: string;
    battery: number;
    ignition: number;
    lat: number;
    speed: number;
    direction: number;
    timestamp: string
};


export const useMQTT = (brokerUrl: string, topic: string) => {
    const [mqttMessages, setMqttMessages] = useState<Record[]>([]);

    useEffect(() => {
        const client = mqtt.connect(brokerUrl, { port: 443 });
        console.log(`Connecting to listen to ${topic}`);

        client.on("connect", () => {
            console.log("Connected to MQTT broker");
            client.subscribe(topic);
        });

        client.on("message", (top: string, message) => {
            if (topic === top) {
                const parse: Record = JSON.parse(message.toString());
                setMqttMessages((prevState) => {
                    return [parse, ...prevState];
                });
            }
        });

        return () => {
            console.log("Disconnecting from MQTT broker");
            client.end();
        };
    }, [topic]);

    return mqttMessages;
};
