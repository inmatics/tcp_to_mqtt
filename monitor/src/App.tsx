import React from 'react';
import { useMQTT } from "./useMqtt";
import RecordsTable from "./Table";
import {Map} from "./Map";

function App() {
    // Define the MQTT topic and broker URL
    const topic = "devices/new";
    const brokerUrl = "wss://mqtt.inmatics.io";

    // Fetch MQTT data using the custom hook
    const mqttData = useMQTT(brokerUrl, topic);

    return (
        <div>
            <Map entries={mqttData} />
            <RecordsTable records={mqttData} />
        </div>
    );
}

// Export the App component as the default export
export default App;
