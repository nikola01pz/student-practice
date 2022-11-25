import React, { useState, useEffect } from "react";

export default function OfferSidebar() {
  const [offer, setOffer] = useState({});

  // mala pomoc u razumijevanju responsea:
  //   {"lige": [{"ime-lige" : "moja-liga"}, {"ime-lige": "druga-liga"}]}
  // istraziti map funkciju u JS
  // { "key": "value"} - json objekt
  // "lige" : [] - key value par, ali i: "lige" : "nekitekst" - key value par, ali i:  "lige" : 12  - key value par

  // u ovom slucaju key ce biti array koji sadrzi objekt
  // iz tog objekta mapirati po odredenom keyu (recimo ime-lige)

  useEffect(() => {
    //   async function fetchOffer() {
    //     const response = await fetch(
    //       url
    //     );
    //     console.log("response ", response);
    //     const json = response.json();
    //     console.log("json ", json);
    //     setOffer(offer);
    //     console.log("offer ", offer);
    //   }
  }, []);

  return <>Ponuda</>;
}
