// Define the Record type for incoming messages
export type Record = {
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