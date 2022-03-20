package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	cmRepository "github.com/Benyam-S/onemembership/common/repository"
	cmService "github.com/Benyam-S/onemembership/common/service"
	"github.com/Benyam-S/onemembership/transaction"

	spRepository "github.com/Benyam-S/onemembership/serviceprovider/repository"
	spService "github.com/Benyam-S/onemembership/serviceprovider/service"

	urRepository "github.com/Benyam-S/onemembership/user/repository"
	urService "github.com/Benyam-S/onemembership/user/service"

	fdRepository "github.com/Benyam-S/onemembership/feedback/repository"
	fdService "github.com/Benyam-S/onemembership/feedback/service"

	sbRepository "github.com/Benyam-S/onemembership/subscription/repository"
	sbService "github.com/Benyam-S/onemembership/subscription/service"

	sbpRepository "github.com/Benyam-S/onemembership/subscriptionplan/repository"
	sbpService "github.com/Benyam-S/onemembership/subscriptionplan/service"

	ptRepository "github.com/Benyam-S/onemembership/project/repository"
	ptService "github.com/Benyam-S/onemembership/project/service"

	prRepository "github.com/Benyam-S/onemembership/preference/repository"
	prService "github.com/Benyam-S/onemembership/preference/service"

	trRepository "github.com/Benyam-S/onemembership/transaction/repository"
	trService "github.com/Benyam-S/onemembership/transaction/service"

	delRepository "github.com/Benyam-S/onemembership/deleted/repository"
	delService "github.com/Benyam-S/onemembership/deleted/service"

	tspRepository "github.com/Benyam-S/onemembership/client/bot/tempserviceprovider/repository"
	tspService "github.com/Benyam-S/onemembership/client/bot/tempserviceprovider/service"

	clRepository "github.com/Benyam-S/onemembership/client/bot/client/repository"
	clService "github.com/Benyam-S/onemembership/client/bot/client/service"

	tgBotH "github.com/Benyam-S/go-tg-bot/handler"
	tgBotL "github.com/Benyam-S/go-tg-bot/log"

	handler "github.com/Benyam-S/onemembership/client/bot/handler"

	"github.com/Benyam-S/onemembership/client/bot"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	configFilesDir string
	redisClient    *redis.Client
	mysqlDB        *gorm.DB
	sysConfig      *entity.SystemConfig
	botHandler     *handler.ServiceProviderBotHandler
	err            error
)

// initServer initialize the web server for takeoff
func initServer() {

	// Reading data from config.server.json file and creating the systemconfig  object
	sysConfigDir := filepath.Join(configFilesDir, "/config.server.json")
	sysConfigData, _ := ioutil.ReadFile(sysConfigDir)

	// Reading data from config.onemembership.json file
	onemembershipConfig := make(map[string]interface{})
	onemembershipConfigDir := filepath.Join(configFilesDir, "/config.onemembership.json")
	onemembershipConfigData, _ := ioutil.ReadFile(onemembershipConfigDir)

	// Reading data from account.api.telebirr.json file
	telebirrConfig := make(map[string]interface{})
	telebirrConfigDir := filepath.Join(configFilesDir, "accounts/account.api.telebirr.json")
	telebirrConfigData, _ := ioutil.ReadFile(telebirrConfigDir)

	// Reading data from key.telebirr.pem file
	telebirrPublicKeyDir := filepath.Join(configFilesDir, "keys/public.telebirr.pem")
	telebirrPublicKeyData, err := ioutil.ReadFile(telebirrPublicKeyDir)
	if err != nil {
		panic(err)
	}

	// Reading data from assets.bot.json file
	botAssets := new(bot.Asset)
	botAssetsDir := filepath.Join(configFilesDir, "../client/bot/assets/assets.bot.json")
	botAssetsData, _ := ioutil.ReadFile(botAssetsDir)

	err = json.Unmarshal(sysConfigData, &sysConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(onemembershipConfigData, &onemembershipConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(telebirrConfigData, &telebirrConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(botAssetsData), botAssets)
	if err != nil {
		panic(err)
	}

	botAPIAccessPoint, ok1 := onemembershipConfig["api_access_point"].(string)
	spBotAPIToken, ok2 := onemembershipConfig["sp_bot_api_token"].(string)
	spBotUsername, ok3 := onemembershipConfig["sp_bot_username"].(string)
	spBotIDString, ok4 := onemembershipConfig["sp_bot_id"].(string)

	telebirrAPIAccessPoint, ok5 := onemembershipConfig["api_access_point"].(string)
	telebirrAppID, ok6 := telebirrConfig["app_id"].(string)
	telebirrAppKey, ok7 := telebirrConfig["app_key"].(string)
	telebirrNotifyURL, ok8 := telebirrConfig["notify_url"].(string)
	telebirrReturnURL, ok9 := telebirrConfig["return_url"].(string)
	telebirrShortCode, ok10 := telebirrConfig["short_code"].(string)
	telebirrTransactionFee, ok11 := telebirrConfig["transaction_fee"].(float64)

	spBotID, err := strconv.ParseInt(spBotIDString, 10, 64)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 || !ok10 || !ok11 || err != nil {
		panic(errors.New("unable to parse onemembership config data"))
	}

	// Setting environmental variables so they can be used any where on the application
	os.Setenv("config_files_dir", configFilesDir)
	os.Setenv("bot_domain_address", sysConfig.BotDomainAddress)
	os.Setenv("bot_client_server_port", sysConfig.BotClientServerPort)
	os.Setenv("logs_path", sysConfig.LogsPath)
	os.Setenv("archives_path", sysConfig.ArchivesPath)

	telebirrAPIAccount := &transaction.TelebirrAPIAccount{
		AccessPoint:    telebirrAPIAccessPoint,
		AppID:          telebirrAppID,
		AppKey:         telebirrAppKey,
		NotifyURL:      telebirrNotifyURL,
		ReturnURL:      telebirrReturnURL,
		ShortCode:      telebirrShortCode,
		TransactionFee: telebirrTransactionFee,
		PublicKey:      telebirrPublicKeyData,
	}

	logContainer := &log.LogContainer{
		UserLogFile:             filepath.Join(sysConfig.LogsPath, sysConfig.Logs["user_log_file"]),
		ServiceProviderLogFile:  filepath.Join(sysConfig.LogsPath, sysConfig.Logs["service_provider_log_file"]),
		ProjectLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["project_log_file"]),
		SubscriptionPlanLogFile: filepath.Join(sysConfig.LogsPath, sysConfig.Logs["subscription_plan_log_file"]),
		SubscriptionLogFile:     filepath.Join(sysConfig.LogsPath, sysConfig.Logs["subscription_log_file"]),
		TransactionLogFile:      filepath.Join(sysConfig.LogsPath, sysConfig.Logs["transaction_log_file"]),
		DeletedLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["deleted_log_file"]),
		ServerLogFile:           filepath.Join(sysConfig.LogsPath, sysConfig.Logs["server_log_file"]),
		BotLogFile:              filepath.Join(sysConfig.LogsPath, sysConfig.Logs["bot_log_file"]),
		ErrorLogFile:            filepath.Join(sysConfig.LogsPath, sysConfig.Logs["error_log_file"]),
		ArchiveLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["archive_log_file"]),
	}

	logger := log.NewLogger(logContainer, log.Normal)
	logger.SetFlag(log.Normal) // Setting logging on

	// Initializing the database with the needed tables and values
	initDB()

	userRepo := urRepository.NewUserRepository(mysqlDB)
	urPasswordRepo := urRepository.NewUserPasswordRepository(mysqlDB)
	serviceProviderRepo := spRepository.NewServiceProviderRepository(mysqlDB)
	spPasswordRepo := spRepository.NewSPPasswordRepository(mysqlDB)
	spWalletRepo := spRepository.NewSPWalletRepository(mysqlDB)
	subscriptionRepo := sbRepository.NewSubscriptionRepository(mysqlDB)
	spSubscriptionRepo := sbRepository.NewSPSubscriptionRepository(mysqlDB)
	subscriptionPlanRepo := sbpRepository.NewSubscriptionPlanRepository(mysqlDB)
	spSubscriptionPlanRepo := sbpRepository.NewSPSubscriptionPlanRepository(mysqlDB)
	planChatLinkRepo := sbpRepository.NewPlanChatLinkRepository(mysqlDB)
	userChatLinkRepo := sbpRepository.NewUserChatLinkRepository(mysqlDB)
	projectRepo := ptRepository.NewProjectRepository(mysqlDB)
	projectChatLinkRepo := ptRepository.NewProjectChatLinkRepository(mysqlDB)
	paymentGatewayRepo := trRepository.NewPaymentGatewayRepository(mysqlDB)
	subscriptionTransactionRepo := trRepository.NewSubscriptionTransactionRepository(mysqlDB)
	spSubscriptionTransactionRepo := trRepository.NewSPSubscriptionTransactionRepository(mysqlDB)
	spPayrollTransactionRepo := trRepository.NewSPPayrollTransactionRepository(mysqlDB)
	preferenceRepo := prRepository.NewPreferenceRepository(mysqlDB)
	deletedUserRepo := delRepository.NewDeletedUserRepository(mysqlDB)
	deletedServiceProviderRepo := delRepository.NewDeletedServiceProviderRepository(mysqlDB)
	deletedSubscriptionTransactionRepo := delRepository.NewDeletedSubscriptionTransactionRepository(mysqlDB)
	deletedSPSTRepo := delRepository.NewDeletedSPSubscriptionTransactionRepository(mysqlDB)
	deletedSPPTRepo := delRepository.NewDeletedSPPayrollTransactionRepository(mysqlDB)
	commonRepo := cmRepository.NewCommonRepository(mysqlDB)
	languageRepo := cmRepository.NewLanguageRepository(mysqlDB)
	languageEntryRepo := cmRepository.NewLanguageEntryRepository(mysqlDB)
	feedbackRepo := fdRepository.NewFeedbackRepository(mysqlDB)

	commonService := cmService.NewCommonService(commonRepo, languageRepo, languageEntryRepo, logger)
	preferenceService := prService.NewClientPreferenceService(preferenceRepo, languageRepo, logger)
	deletedService := delService.NewDeletedService(deletedUserRepo, deletedServiceProviderRepo,
		deletedSubscriptionTransactionRepo, deletedSPSTRepo, deletedSPPTRepo, logger)
	feedbackService := fdService.NewFeedbackService(feedbackRepo, serviceProviderRepo, userRepo, logger)
	userService := urService.NewUserService(userRepo, urPasswordRepo, preferenceService,
		feedbackService, deletedService, commonService, logger)
	serviceProviderService := spService.NewServiceProviderService(serviceProviderRepo, spPasswordRepo, spWalletRepo,
		preferenceService, feedbackService, deletedService, commonService, logger)
	subscriptionService := sbService.NewSubscriptionService(subscriptionRepo, spSubscriptionRepo, logger)
	subscriptionPlanService := sbpService.NewSubscriptionPlanService(subscriptionPlanRepo, spSubscriptionPlanRepo,
		planChatLinkRepo, userChatLinkRepo, commonService, logger)
	projectService := ptService.NewProjectService(projectRepo, projectChatLinkRepo, commonService, logger)
	transactionService := trService.NewTransactionService(paymentGatewayRepo, subscriptionTransactionRepo,
		spSubscriptionTransactionRepo, spPayrollTransactionRepo, telebirrAPIAccount, commonService, logger)

	// ----- Bot level init -----
	tempServiceProviderRepo := tspRepository.NewTempServiceProviderRepository(mysqlDB)
	clientRepo := clRepository.NewClientRepository(mysqlDB)

	tempServiceProviderService := tspService.NewTempServiceProviderService(tempServiceProviderRepo,
		languageRepo, commonService, logger)
	clientService := clService.NewClientService(clientRepo, logger)

	// ----- Creating store -----
	store := tools.NewRedisStore(redisClient)

	botLogContainer := &tgBotL.LogContainer{
		BotLogFile:   filepath.Join(sysConfig.LogsPath, sysConfig.Logs["bot_log_file"]),
		ErrorLogFile: filepath.Join(sysConfig.LogsPath, sysConfig.Logs["error_log_file"]),
	}
	bot := tgBotH.NewTelegramBotHandler(botAPIAccessPoint, spBotAPIToken, spBotID, spBotUsername, logger, botLogContainer)

	botHandler = handler.NewServiceProviderBotHandler(tempServiceProviderService, clientService,
		userService, serviceProviderService, subscriptionService, subscriptionPlanService, projectService,
		transactionService, preferenceService, feedbackService, deletedService, commonService, bot, botAssets, store, logger)

}

// initDB initialize the database for takeoff
func initDB() {

	redisDB, err := strconv.ParseInt(sysConfig.RedisClient["database"], 0, 0)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     sysConfig.RedisClient["address"] + ":" + sysConfig.RedisClient["port"],
		Password: sysConfig.RedisClient["password"], // no password set
		DB:       int(redisDB),                      // use default DB
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

	if err != nil {
		panic(err)
	}

	mysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sysConfig.MysqlClient["user"], sysConfig.MysqlClient["password"],
		sysConfig.MysqlClient["address"], sysConfig.MysqlClient["port"], sysConfig.MysqlClient["database"]))

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database: mysql @GORM")

	// Creating and Migrating tables from the structures
	mysqlDB.AutoMigrate(&entity.Subscription{})
	mysqlDB.AutoMigrate(&entity.SPSubscription{})
	mysqlDB.AutoMigrate(&entity.SubscriptionPlan{})
	mysqlDB.AutoMigrate(&entity.SPSubscriptionPlan{})
	mysqlDB.AutoMigrate(&entity.PlanChatLink{})
	mysqlDB.AutoMigrate(&entity.UserChatLink{})
	mysqlDB.AutoMigrate(&entity.ServiceProvider{})
	mysqlDB.AutoMigrate(&entity.SPPassword{})
	mysqlDB.AutoMigrate(&entity.SPWallet{})
	mysqlDB.AutoMigrate(&entity.User{})
	mysqlDB.AutoMigrate(&entity.UserPassword{})
	mysqlDB.AutoMigrate(&entity.Language{})
	mysqlDB.AutoMigrate(&entity.LanguageEntry{})
	mysqlDB.AutoMigrate(&entity.ClientPreference{})
	mysqlDB.AutoMigrate(&entity.Project{})
	mysqlDB.AutoMigrate(&entity.ProjectChatLink{})
	mysqlDB.AutoMigrate(&entity.PaymentGateway{})
	mysqlDB.AutoMigrate(&entity.SubscriptionTransaction{})
	mysqlDB.AutoMigrate(&entity.SPSubscriptionTransaction{})
	mysqlDB.AutoMigrate(&entity.SPPayrollTransaction{})
	mysqlDB.AutoMigrate(&entity.Feedback{})

	mysqlDB.AutoMigrate(&entity.DeletedServiceProvider{})
	mysqlDB.AutoMigrate(&entity.DeletedUser{})
	mysqlDB.AutoMigrate(&entity.DeletedSubscriptionTransaction{})
	mysqlDB.AutoMigrate(&entity.DeletedSPSubscriptionTransaction{})
	mysqlDB.AutoMigrate(&entity.DeletedSPPayrollTransaction{})

	// ----- Bot level database -----
	mysqlDB.AutoMigrate(&bot.TempUser{})
	mysqlDB.AutoMigrate(&bot.TempServiceProvider{})
	mysqlDB.AutoMigrate(&bot.Client{})

	// Setting foreign key constraint
	mysqlDB.Model(&entity.UserPassword{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SPPassword{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SPWallet{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SPSubscription{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SPSubscription{}).AddForeignKey("subscription_plan_id", "subscription_plans(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Project{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.ProjectChatLink{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SubscriptionPlan{}).AddForeignKey("project_id", "projects(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.PlanChatLink{}).AddForeignKey("plan_id", "subscription_plans(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.UserChatLink{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.UserChatLink{}).AddForeignKey("plan_id", "subscription_plans(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("subscriber_id", "users(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("provider_id", "service_providers(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("project_id", "projects(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.Subscription{}).AddForeignKey("subscription_plan_id", "subscription_plans(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.SubscriptionTransaction{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SubscriptionTransaction{}).AddForeignKey("plan_id", "subscription_plans(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.SPSubscriptionTransaction{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.SPSubscriptionTransaction{}).AddForeignKey("plan_id", "subscription_plans(id)", "SET NULL", "CASCADE")
	mysqlDB.Model(&entity.SPPayrollTransaction{}).AddForeignKey("provider_id", "service_providers(id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.LanguageEntry{}).AddForeignKey("code", "languages(code)", "CASCADE", "CASCADE")
	mysqlDB.Model(&entity.ClientPreference{}).AddForeignKey("language", "languages(code)", "CASCADE", "CASCADE")

}

func main() {
	configFilesDir = "/Users/administrator/go/src/github.com/Benyam-S/onemembership_sp/config"

	// Initializing the server
	initServer()
	defer mysqlDB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", tools.MiddlewareFactory(botHandler.HandleWebHook, botHandler.ParseRequest))

	http.ListenAndServe(":"+os.Getenv("bot_client_server_port"), router)
}
