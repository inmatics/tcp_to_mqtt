import mqtt from "mqtt";
import { useEffect, useState } from "react";
import {Record} from "./types";


// Custom hook to handle MQTT messages
export const useMQTT = (brokerUrl: string, topic: string) => {
    // Store MQTT messages in state
    const [mqttMessages, setMqttMessages] = useState<Record[]>([]);

    // Use useEffect to setup and cleanup the MQTT connection
    useEffect(() => {
        // Connect to the MQTT broker
        const client = mqtt.connect(brokerUrl, { port: 443 });
        console.log(`Connecting to listen to ${topic}`);

        // When connected, subscribe to the given topic
        client.on("connect", () => {
            console.log("Connected to MQTT broker");
            client.subscribe(topic);
        });

        // When a message is received, check if it's from the subscribed topic and process it
        client.on("message", (top: string, message) => {
            if (topic === top) {
                // Parse the message and update the state with the new record
                const parse: Record = JSON.parse(message.toString());
                setMqttMessages((prevState) => {
                    return [parse, ...prevState];
                });
            }
        });

        // Cleanup function to disconnect from the MQTT broker when the component is unmounted
        return () => {
            console.log("Disconnecting from MQTT broker");
            client.end();
        };
    }, [topic]); // The effect depends on the topic, so it's included in the dependency array

    // Return the list of MQTT messages
    return mqttMessages;
};
