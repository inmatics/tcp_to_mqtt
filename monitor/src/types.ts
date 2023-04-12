// Define the Record type for incoming messages
export type Record = {
    angle: number;
    lng: number;
    odometer: number;
    raw_message: string;
    imei: string;
    battery: number;
    ignition: number;
    lat: number;
    speed: number;
    direction: number;
    timestamp: string
};