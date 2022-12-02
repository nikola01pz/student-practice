import React, { useState, useEffect } from "react";
import "../styles/Offer.css";

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

  // useEffect(() => {}, []);

  const [showOffer, setshowOffer] = useState(false);
  const [currentlySelectedLeague, setCurrentlySelectedLeague] = useState({});

  const handleOffersDialog = (liga) => {
    console.log("liga: ", liga.razrade[0].tipovi);
    setshowOffer(true);
    setCurrentlySelectedLeague(liga);
  };

  const [tipsDialogOpen, setTipsDialogOpen] = useState(false);
  const handleTipsDialog = (liga) => {
    console.log(liga);
    setTipsDialogOpen(true);
  };

  return (
    <>
      <div className="flex-container">
        <div className="league-sidebar">
          {offerHardCoded.lige.map((league, index) => {
            return (
              <div
                className="league-name"
                key={index}
                onClick={() => handleOffersDialog(league)}
                open={showOffer}
                onClose={handleOffersDialog}
              >
                <> {league.naziv} </>
              </div>
            );
          })}
        </div>

        {showOffer && (
          <div
            className="league-tips"
            onClick={() => handleTipsDialog(currentlySelectedLeague)}
            open={tipsDialogOpen}
            onClose={handleTipsDialog}
          >
            {currentlySelectedLeague.razrade[0].tipovi.map((tips, index) => {
              return (
                <div className="single-tip" key={index}>
                  <>{tips.naziv}</>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </>
  );
}
