import React, {Children, cloneElement, isValidElement, useEffect, useRef, useState,} from "react";
import {useDeepCompareEffectForMaps} from "./maps";
import {Record} from "./types";

export function getCenter(maps: typeof google.maps, locationEntries: Record[]) {
    const bound = new maps.LatLngBounds();
    for (let i = 0; i < locationEntries.length; i++) {
        bound.extend(
            new maps.LatLng(locationEntries[i].lat, locationEntries[i].lng)
        );
    }

    const boundCenter = bound.getCenter();
    return { bound, boundCenter };
}

export function MyMapComponent(props: any) {
    const {
        children,
        locationEntries,
        onClick,
        onIdle,
        style,
        ...options
    } = props;

    const ref = useRef<any>();
    const [map, setMap] = useState<google.maps.Map>();
    const [center, setCenter] = useState<google.maps.LatLngLiteral>();
    const [bounds, setBounds] = useState<google.maps.LatLngBounds>();

    const maps = window.google.maps;

    useEffect(() => {
        if (window.google) {
            const {bound, boundCenter} = getCenter(maps, locationEntries);

            setBounds(bound);
            setCenter({
                lat: boundCenter.lat && boundCenter.lat(),
                lng: boundCenter.lng && boundCenter.lng(),
            });
        }
    }, [locationEntries, maps, maps.LatLng, maps.LatLngBounds]);

    useEffect(() => {
        if (ref.current && !map && center?.lat) {
            let map1 = new maps.Map(ref.current, {maxZoom: 19});
            bounds && map1.fitBounds(bounds);
            setMap(map1);
        }
    }, [ref, map, center, maps.Map, bounds]);

    useDeepCompareEffectForMaps(() => {
        if (map) {
            map.setOptions(options);
        }
    }, [map, options]);

    useEffect(() => {
        if (map) {
            ["click", "idle"].forEach((eventName) =>
                maps.event.clearListeners(map, eventName)
            );

            if (onClick) {
                map.addListener("click", onClick);
            }

            if (onIdle) {
                map.addListener("idle", () => onIdle(map));
            }
        }
    }, [map, onClick, onIdle, maps.event]);

    return (
        <>
            <div ref={ref} style={style}/>
            {(Children || []).map(children, (child) => {
                if (isValidElement(child)) {
                    // @ts-ignore
                    return cloneElement(child, {map});
                }
            })}
        </>
    );
}
