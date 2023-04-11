import React from 'react';
import {Table} from "antd";
import {Record} from "./types";

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


function RecordsTable(props: {records: Record[];}) {
    const {records} = props

    return (
        <Table dataSource={records} columns={columns}/>
    );
}

export default RecordsTable;
