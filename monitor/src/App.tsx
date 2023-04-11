import React from 'react';
import {Table} from "antd";
import {useMQTT} from "./useMqtt";

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

function App() {
    const topic = "devices/new";
    const brokerUrl = "wss://mqtt.example.com";

    return <Table dataSource={useMQTT(brokerUrl, topic)} columns={columns}/>;
}

export default App;
