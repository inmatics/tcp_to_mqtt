// Import required dependencies
import React from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import L from 'leaflet';

// Define the custom icon for the map markers
const customIcon = new L.Icon({
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.3.1/images/marker-icon.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
});

// MapDisplay component props
type MapDisplayProps = {
    data: {
        lat: number;
        lng: number;
        Imei: string;
        speed: number;
        timestamp: string;
    }[];
};

// MapDisplay component
const MapDisplay: React.FC<MapDisplayProps> = ({ data }) => {
    return (
        <MapContainer center={[data[0]?.lat || 0, data[0]?.lng || 0]} zoom={13} style={{ height: '100%', width: '100%' }}>
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {data.map((record, index) => (
                <Marker key={index} position={[record.lat, record.lng]} icon={customIcon}>
                    <Popup>
                        IMEI: {record.Imei} <br />
                        Speed: {record.speed} <br />
                        Timestamp: {record.timestamp}
                    </Popup>
                </Marker>
            ))}
        </MapContainer>
    );
};

export default MapDisplay;
