import React, { useState, useEffect } from "react";

export default function OfferSidebar() {
  const [offer, setOffer] = useState({});

  const offerHardCoded = {
    lige: [
      {
        naziv: "SUPER TIP- NOGOMET",
        razrade: [
          {
            tipovi: [
              {
                naziv: "1",
              },
              {
                naziv: "X",
              },
              {
                naziv: "2",
              },
              {
                naziv: "1X",
              },
              {
                naziv: "X2",
              },
              {
                naziv: "12",
              },
              {
                naziv: "f+2",
              },
            ],
            ponude: [8991909, 8991908, 8991913, 8991933],
          },
        ],
      },
      {
        naziv: "SUPER TIP-TENIS",
        razrade: [
          {
            tipovi: [
              {
                naziv: "1",
              },
              {
                naziv: "2",
              },
            ],
            ponude: [8991284],
          },
        ],
      },
      {
        naziv: "Najtipovaniji događaji",
        razrade: [
          {
            tipovi: [
              {
                naziv: "1",
              },
              {
                naziv: "X",
              },
              {
                naziv: "2",
              },
              {
                naziv: "1X",
              },
              {
                naziv: "X2",
              },
              {
                naziv: "12",
              },
              {
                naziv: "f+2",
              },
            ],
            ponude: [
              8679662, 8679707, 8679702, 8679700, 8679699, 8679704, 8679705,
              8974498, 8679698, 8977288, 8679718, 8679715, 8679717, 8679713,
              8679709, 8679712,
            ],
          },
        ],
      },
    ],
  };

  const mzObj = { lige: [{naziv: "moja liga" }]};

  console.log(offerHardCoded);

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

  return <>Ponuda: {mzObj.lige[0].naziv}</>;
}
