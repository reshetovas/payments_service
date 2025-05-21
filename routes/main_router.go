package routes

import "github.com/gorilla/mux"

func MainRouter(paymentRoutes *PaymentRoutes, bonusRoutes *BonusRoutes, userRoutes *UserRoutes) *mux.Router {
	mainRouter := mux.NewRouter()

	paymentRouter := paymentRoutes.PaymentRouter()
	bonusRouter := bonusRoutes.BonusRouter()
	userRouter := userRoutes.UserRouter()

	mainRouter.PathPrefix("/payment").Handler(paymentRouter)
	mainRouter.PathPrefix("/payments").Handler(paymentRouter)
	mainRouter.PathPrefix("/bonus").Handler(bonusRouter)
	mainRouter.PathPrefix("/bonuses").Handler(bonusRouter)
	mainRouter.PathPrefix("/users").Handler(userRouter)

	return mainRouter
}

func MainWSRouter(wsRoutes *WSRoutes) *mux.Router {
	mainRouter := mux.NewRouter()
	wsRouter := wsRoutes.WSRouter()

	mainRouter.PathPrefix("/").Handler(wsRouter)

	return mainRouter
}
