package http

import (
	"bettingAPI/internal/mysql"
)

func isStakeSufficient(stake float32) bool {
	return stake >= 2
}

func hasSufficientFunds(balance, stake float32) bool {
	return balance >= stake
}

func (h *handler) calculateTotalCoefficient(betSlip mysql.BetSlipRequest) (float32, error) {
	offerTips, err := h.db.GetOfferTipCoefficients(betSlip.Bets)
	if err != nil {
		return 0, err
	}
	var totalCoefficient float32 = 1.0
	for i := range offerTips {
		totalCoefficient = totalCoefficient * float32(offerTips[i].Coefficient)
	}
	return totalCoefficient, nil
}

func hasOnlyOneTipPerOffer(bets []mysql.Bet) bool {
	for i := range bets {
		for j := range bets {
			if i != j {
				if bets[i].OfferID == bets[j].OfferID {
					return false
				}
			}
		}
	}
	return true
}
