import { Wrapper } from "@googlemaps/react-wrapper";
import { Status } from "@googlemaps/react-wrapper/src";
import React from "react";
import {Record} from "./types";
import {Marker} from "./Marker";
import {MyMapComponent} from "./MyMapComponent";

const GOOGLE_MAPS_KEY = "AIzaSyB0jdQ-eoAtaVqxaK7K_zIcgGk_wqeSzPY";

const renderStatus = (status: Status) => {
  if (status === "LOADING") return <h3>{status} ..</h3>;
  if (status === "FAILURE") return <h3>{status} ...</h3>;
  return <h3>{status} ...</h3>;
};

export const Map = (props: { entries: Record[] }) => {
  const { entries } = props;
  const [zoom, setZoom] = React.useState(3); // initial zoom

  return (
    <Wrapper apiKey={GOOGLE_MAPS_KEY} render={renderStatus}>
      <MyMapComponent
        locationEntries={entries}
        onClick={undefined}
        onIdle={(m: google.maps.Map) => setZoom(m.getZoom()!)}
        style={{ height: "400px" }}
      >
        {entries.map((x: Record) => (
          <Marker position={x} key={x.timestamp} zoom={zoom} />
        ))}
      </MyMapComponent>
    </Wrapper>
  );
};
