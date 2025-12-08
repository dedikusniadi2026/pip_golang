package main

import (
	"auth-service/handler"
	"auth-service/repository"
	"auth-service/service"
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {

	userRepo := &repository.UserRepository{DB: db}
	driverRepo := &repository.DriverRepository{DB: db}
	bookingRepo := &repository.BookingRepository{DB: db}
	popularRepo := &repository.PopularDestinationRepository{DB: db}
	carRepo := repository.NewCarRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	authService := &service.AuthService{Repo: userRepo}
	driverService := &service.DriverService{Repo: driverRepo}
	bookingService := &service.BookingService{Repo: bookingRepo}
	popularService := &service.PopularDestinationService{Repo: popularRepo}
	carService := service.NewCarService(carRepo)
	paymentService := service.NewPaymentService(paymentRepo)

	authHandler := &handler.AuthHandler{AuthService: authService}
	driverHandler := &handler.DriverHandler{Service: driverService}
	bookingHandler := &handler.BookingHandler{BookingService: bookingService}
	popularHandler := &handler.PopularDestinationHandler{Service: popularService}
	carHandler := handler.NewCarHandler(carService)
	paymentHandler := handler.PaymentHandler{Service: paymentService}

	maintenanceRepo := repository.NewMaintenanceRepository(db)
	maintenanceService := service.NewMaintenanceService(maintenanceRepo)
	maintenanceHandler := handler.NewMaintenanceHandler(maintenanceService)

	assignmentRepo := repository.NewAssignmentsRepository(db)
	assignmentService := service.NewAssignmentsService(assignmentRepo)
	assignmentHandler := handler.NewAssignmentHandler(assignmentService)

	tripRepo := repository.NewTripRepo(db)
	tripService := service.NewTripService(tripRepo)
	tripHandler := handler.NewTripHandler(tripService)

	tripHistoryRepo := repository.NewTripHistoryRepository(db)
	tripHistoryService := service.NewTripHistoryService(tripHistoryRepo)

	tripHistoryHandler := handler.NewTripHistoryHandler(tripHistoryService)

	bookingTrendsRepo := repository.NewBookingTrendsRepository(db)
	bookingTrendsService := service.NewBookingTrendsService(bookingTrendsRepo)
	bookingTrendsHandler := handler.NewBookingTrendsHandler(bookingTrendsService)

	pdfRepo := repository.NewPDFRepository(db)
	pdfService := &service.PDFService{
		TripRepo:     pdfRepo,
		PDFGenerator: &service.DefaultPDFGenerator{},
	}
	pdfHandler := handler.NewPDFHandler(pdfService)

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	dashboardTripRepo := repository.NewDashboardTripRepository(db)
	dashboardTripService := service.NewDashboardTripService(dashboardTripRepo)
	dashboardTripHandler := handler.NewDashboardTripHandler(dashboardTripService)

	carTypeRepo := repository.NewCarTypeRepository(db)
	carModelRepo := repository.NewCarModelRepository(db)

	carTypeService := service.NewCarTypeService(carTypeRepo)
	carModelService := service.NewCarModelService(carModelRepo)

	carTypeHandler := handler.NewCarTypeHandler(carTypeService)
	carModelHandler := handler.NewCarModelRepositoryHandler(carModelService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// REGISTER ROUTES
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)

	r.GET("/drivers", driverHandler.GetAll)
	r.POST("/drivers", driverHandler.Create)
	r.GET("/drivers/:id", driverHandler.GetByID)
	r.PUT("/drivers/:id", driverHandler.Update)
	r.DELETE("/drivers/:id", driverHandler.Delete)

	r.POST("/booking", bookingHandler.Create)
	r.GET("/booking", bookingHandler.GetAll)
	r.DELETE("/booking/:id", bookingHandler.Delete)

	r.GET("/car", carHandler.GetAll)
	r.GET("/car/:id", carHandler.GetByID)
	r.POST("/car", carHandler.Create)
	r.PUT("/car/:id", carHandler.Update)
	r.DELETE("/car/:id", carHandler.Delete)

	r.GET("/payments", paymentHandler.GetPayments)
	r.GET("/paymentsStats", paymentHandler.GetPaymentStats)
	r.GET("/payments/:id", paymentHandler.GetPaymentByID)
	r.POST("/payments", paymentHandler.CreatePayment)
	r.PUT("/payments/:id", paymentHandler.UpdatePayment)
	r.DELETE("/payments/:id", paymentHandler.DeletePayment)

	r.GET("/assignments/:vehicle_id", assignmentHandler.FindByVehicle)
	r.POST("/assignments", assignmentHandler.Create)
	r.PUT("/assignments/:id", assignmentHandler.Update)
	r.DELETE("/assignments/:id", assignmentHandler.Delete)

	r.POST("/maintenance", maintenanceHandler.Create)
	r.GET("/maintenance/vehicle/:vehicle_id", maintenanceHandler.FindByVehicle)
	r.PUT("/maintenance/:id", maintenanceHandler.Update)
	r.DELETE("/maintenance/:id", maintenanceHandler.Delete)

	r.GET("/trips/vehicle/:vehicle_id", tripHandler.FindByVehicle)
	r.POST("/trips", tripHandler.Create)
	r.PUT("/trips/:id", tripHandler.Update)
	r.DELETE("/trips/:id", tripHandler.Delete)

	r.GET("/tripsStats", tripHandler.GetTripTotal)
	r.GET("/booking-trends", bookingTrendsHandler.GetTrends)
	r.GET("/trip-history", tripHistoryHandler.GetTripHistory)
	r.GET("/chart", popularHandler.GetAll)

	r.GET("/car-types", carTypeHandler.GetAll)
	r.GET("/car-models", carModelHandler.GetAll)
	r.GET("/car-models/:id", carModelHandler.GetByID)
	r.GET("/car-types/:id", carTypeHandler.GetByID)

	r.GET("/receipt/:trip_id", pdfHandler.HandlePDFReceipt)
	r.GET("/api/dashboard", dashboardHandler.GetDashboard)
	r.GET("/api/dashboardtrip/:id", dashboardTripHandler.GetDashboard)

	return r
}
