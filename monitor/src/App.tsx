import React from 'react';
import { useMQTT } from "./useMqtt";
import MapDisplay from './MapDisplay';
import RecordsTable from "./Table";

function App() {
    // Define the MQTT topic and broker URL
    const topic = "devices/new";
    const brokerUrl = "wss://mqtt.inmatics.io";

    // Fetch MQTT data using the custom hook
    const mqttData = useMQTT(brokerUrl, topic);

    return (
        <div>
            <MapDisplay data={mqttData} />
            <RecordsTable records={mqttData} />
        </div>
    );
}

// Export the App component as the default export
export default App;
