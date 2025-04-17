package routes

import "github.com/gorilla/mux"

func MainRouter(paymentRoutes *PaymentRoutes, bonusRoutes *BonusRoutes) *mux.Router {
	mainRouter := mux.NewRouter()

	paymentRouter := paymentRoutes.PaymentRouter()
	bonusRouter := bonusRoutes.BonusRouter()

	mainRouter.PathPrefix("/payment").Handler(paymentRouter)
	mainRouter.PathPrefix("/payments").Handler(paymentRouter)
	mainRouter.PathPrefix("/bonus").Handler(bonusRouter)
	mainRouter.PathPrefix("/bonuses").Handler(bonusRouter)

	return mainRouter
}
