import { useEffect, useState } from "react";
import {Record} from "./types";

export const Marker = (props: { position: Record; zoom: number }) => {
  const [marker, setMarker] = useState<any>();
  const [infoWindow, setInfoWindow] = useState<any>();

  const { position, zoom } = props;

  const infoWindowStyle = "font-size: 14px";

  const infoWindowData = `<div>info</div>`;

  //Create marker
  useEffect(() => {
    if (!marker) {
      setMarker(
        new window.google.maps.Marker({
          icon: {
            anchor: new window.google.maps.Point(0, -3),
            fillColor: "blue",
            fillOpacity: 0.6,
            path: window.google.maps.SymbolPath.BACKWARD_CLOSED_ARROW,
            rotation: position.direction - 180,
            scale: 2,
            strokeWeight: 1,
          },
        })
      );
    }

    // remove marker from map on unmount
    return () => {
      if (marker) {
        marker.setMap(null);
      }
    };
  }, [marker, position.direction, zoom]);

  useEffect(() => {
    let scale: number;

    if (zoom > 13) {
      scale = 4;
    } else {
      scale = 3;
    }

    if (marker) {
      marker.setOptions({
        icon: {
          anchor: new window.google.maps.Point(0, -3),
          fillColor: "blue",
          fillOpacity: 0.6,
          path: window.google.maps.SymbolPath.BACKWARD_CLOSED_ARROW,
          rotation: position.direction - 180,
          scale: scale,
          strokeWeight: 1,
        },
      });
    }
  }, [zoom, marker, position.direction]);

  //Set marker props
  useEffect(() => {
    if (marker) {
      marker.setOptions(props);
    }
  }, [marker, props]);

  // Create Info Window
  useEffect(() => {
    if (!infoWindow) {
      setInfoWindow(
        new window.google.maps.InfoWindow({
          content: infoWindowData,
        })
      );
    }
  }, [infoWindow, infoWindowData]);

  // Render Info Window onCLick
  useEffect(() => {
    if (marker && infoWindow) {
      marker.addListener("click", () => {
        infoWindow.open({
          anchor: marker,
          shouldFocus: false,
        });
      });
    }
  }, [marker, infoWindow]);

  return null;
};
