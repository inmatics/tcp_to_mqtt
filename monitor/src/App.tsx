import React from 'react';
import {Table} from "antd";
import {useMQTT} from "./useMqtt";

function App() {
    const topic = "devices/new";
    const brokerUrl = "wss://mqtt.inmatics.io";

    const locations = useMQTT(brokerUrl,topic);

    const columns = [
        {
            title: 'IMEI',
            dataIndex: 'Imei',
            key: 'Imei',
        },
        {
            title: 'Speed',
            dataIndex: 'speed',
            key: 'speed',
        },
        {
            title: 'Timestamp',
            dataIndex: 'timestamp',
            key: 'timestamp',
        },
        {
            title: 'Lat',
            dataIndex: 'lat',
            key: 'lat',
        },
        {
            title: 'Lng',
            dataIndex: 'lng',
            key: 'lng',
        },
    ];

    return <Table dataSource={locations} columns={columns}/>;
}

export default App;
