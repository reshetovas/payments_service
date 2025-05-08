package routes

import "github.com/gorilla/mux"

func MainRouter(paymentRoutes *PaymentRoutes, bonusRoutes *BonusRoutes, userRoutes *UserRoutes, wsRoutes *WSRoutes) *mux.Router {
	mainRouter := mux.NewRouter()

	paymentRouter := paymentRoutes.PaymentRouter()
	bonusRouter := bonusRoutes.BonusRouter()
	userRouter := userRoutes.UserRouter()
	wsRouter := wsRoutes.WSRouter()

	mainRouter.PathPrefix("/payment").Handler(paymentRouter)
	mainRouter.PathPrefix("/payments").Handler(paymentRouter)
	mainRouter.PathPrefix("/bonus").Handler(bonusRouter)
	mainRouter.PathPrefix("/bonuses").Handler(bonusRouter)
	mainRouter.PathPrefix("/users").Handler(userRouter)
	mainRouter.PathPrefix("/ws").Handler(wsRouter)

	return mainRouter
}
