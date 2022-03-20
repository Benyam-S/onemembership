package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/Benyam-S/onemembership/transaction"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

// Service is a type that defines a transaction service
type Service struct {
	paymentGatewayRepo            transaction.IPaymentGatewayRepository
	subTransactionRepo            transaction.ISubscriptionTransactionRepository
	spSubscriptionTransactionRepo transaction.ISPSubscriptionTransactionRepository
	spPayrollTransactionRepo      transaction.ISPPayrollTransactionRepository
	TelebirrAPI                   *transaction.TelebirrAPIAccount
	cmService                     common.IService
	logger                        *log.Logger
}

// NewTransactionService is a function that returns a new transaction service
func NewTransactionService(paymentGatewayRepository transaction.IPaymentGatewayRepository,
	subscriptionTransactionRepository transaction.ISubscriptionTransactionRepository,
	spSubscriptionTransactionRepository transaction.ISPSubscriptionTransactionRepository,
	spPayrollTransactionRepository transaction.ISPPayrollTransactionRepository,
	telebirrAPIAccount *transaction.TelebirrAPIAccount, commonService common.IService,
	projectLogger *log.Logger) transaction.IService {
	return &Service{paymentGatewayRepo: paymentGatewayRepository, subTransactionRepo: subscriptionTransactionRepository,
		spSubscriptionTransactionRepo: spSubscriptionTransactionRepository,
		spPayrollTransactionRepo:      spPayrollTransactionRepository, TelebirrAPI: telebirrAPIAccount,
		cmService: commonService, logger: projectLogger}
}

// AddPaymentGateway is a method that adds a new payment gateway to the system
func (service *Service) AddPaymentGateway(newPaymentGateway *entity.PaymentGateway) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payment gateway adding process, Payment Gateway => %s",
		newPaymentGateway.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.paymentGatewayRepo.Create(newPaymentGateway)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Payment Gateway => %s, %s",
			newPaymentGateway.ToString(), err.Error()))

		return errors.New("unable to add new payment gateway")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payment gateway adding process, Payment Gateway => %s",
		newPaymentGateway.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// ValidatePaymentGateway is a method that validates a payment gateway entries.
// It checks if the payment gateway has a valid entries or not and return map of errors if any.
func (service *Service) ValidatePaymentGateway(paymentGateway *entity.PaymentGateway) entity.ErrMap {

	errMap := make(map[string]error)

	emptyName, _ := regexp.MatchString(`^\s*$`, paymentGateway.Name)
	if emptyName {
		errMap["name"] = errors.New(`gateway name can not be empty`)
	} else if len(paymentGateway.Name) > 1000 { // Since it may contain special character it is better to use '1000'
		errMap["name"] = errors.New(`gateway name should not be longer than 1000 characters`)
	}

	// Meaning a new payment gateway is being add
	if paymentGateway.ID == 0 {
		if errMap["name"] == nil && !service.cmService.IsUnique("name", paymentGateway.Name, "payment_gateways") {
			errMap["name"] = errors.New(`gateway name already exists`)
		}
	} else {
		// Meaning trying to update payment gateway
		prevPaymentGateway, err := service.paymentGatewayRepo.Find(paymentGateway.ID)

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && errMap["name"] == nil && prevPaymentGateway.Name != paymentGateway.Name {
			if !service.cmService.IsUnique("name", paymentGateway.Name, "payment_gateways") {
				errMap["name"] = errors.New(`gateway name already exists`)
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindPaymentGateway is a method that find and return a payment gateway that matches the id value
func (service *Service) FindPaymentGateway(id int64) (*entity.PaymentGateway, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single payment gateway finding process { Payment Gateway ID : %d }", id),
		service.logger.Logs.TransactionLogFile)

	paymentGateway, err := service.paymentGatewayRepo.Find(id)
	if err != nil {
		return nil, errors.New("no payment gateway found")
	}
	return paymentGateway, nil
}

// AllPaymentGateways is a method that returns all the payment gateway in the system
func (service *Service) AllPaymentGateways() []*entity.PaymentGateway {
	return service.paymentGatewayRepo.All()
}

// UpdatePaymentGateway is a method that updates a payment gateway in the system
func (service *Service) UpdatePaymentGateway(paymentGateway *entity.PaymentGateway) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payment gateway updating process, Payment Gateway => %s",
		paymentGateway.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.paymentGatewayRepo.Update(paymentGateway)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Payment Gateway => %s, %s",
			paymentGateway.ToString(), err.Error()))

		return errors.New("unable to update payment gateway")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payment gateway updating process, Payment Gateway => %s",
		paymentGateway.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// DeletePaymentGateway is a method that deletes a payment gateway from the system using an id
func (service *Service) DeletePaymentGateway(id int64) (*entity.PaymentGateway, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payment gateway deleting process { Payment Gateway ID : %d }",
		id), service.logger.Logs.TransactionLogFile)

	paymentGateway, err := service.paymentGatewayRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting payment gateway { Payment Gateway ID : %d }, %s",
			id, err.Error()))

		return nil, errors.New("unable to delete payment gateway")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payment gateway deleting process, Deleted Payment Gateway => %s",
		paymentGateway.ToString()), service.logger.Logs.TransactionLogFile)
	return paymentGateway, nil
}

// GetTelebirrH5WebURL is a method that generates a H5 web url
func (service *Service) GetTelebirrH5WebURL(userID, planID, receiverName, subject, currencyType, initiatedFrom string,
	receivedAmount float64) (string, error) {

	type TelebirrRequest struct {
		AppID string `json:"appid"`
		Sign  string `json:"sign"`
		USSD  string `json:"ussd"`
	}

	type TelebirrRequestParameters struct {
		AppID          string `json:"appId" url:"appId"`
		AppKey         string `json:"-" url:"appKey"` // since it isn't included in the ussd value no need to change it to json format
		Nonce          string `json:"nonce" url:"nonce"`
		NotifyURL      string `json:"notifyUrl" url:"notifyUrl,omitempty"`
		OutTradeNo     string `json:"outTradeNo" url:"outTradeNo"`
		ReceiverName   string `json:"receiveName" url:"receiveName,omitempty"`
		ReturnURL      string `json:"returnUrl" url:"returnUrl,omitempty"`
		ShortCode      string `json:"shortCode" url:"shortCode"`
		Subject        string `json:"subject" url:"subject,omitempty"`
		TimeoutExpress string `json:"timeoutExpress" url:"timeoutExpress"`
		Timestamp      string `json:"timestamp" url:"timestamp"`
		TotalAmount    string `json:"totalAmount" url:"totalAmount"`
	}

	type TelebirrResponse struct {
		Code    string            `json:"code"`
		Message string            `json:"msg"`
		Data    map[string]string `json:"data"`
	}

	var requestTimeout int64 = 60
	var timestamp int64 = time.Now().Unix()
	var transactionFee float64 = service.TelebirrAPI.TransactionFee
	uniqueID := strings.ReplaceAll(uuid.Must(uuid.NewRandom()).String(), "-", "")
	totalAmount := receivedAmount + transactionFee

	telebirrRequestParameters := &TelebirrRequestParameters{
		AppID:          service.TelebirrAPI.AppID,
		AppKey:         service.TelebirrAPI.AppKey,
		Nonce:          uniqueID,
		OutTradeNo:     "T_" + uniqueID,
		NotifyURL:      service.TelebirrAPI.NotifyURL,
		ReceiverName:   receiverName,
		ReturnURL:      service.TelebirrAPI.ReturnURL,
		ShortCode:      service.TelebirrAPI.ShortCode,
		Subject:        subject,
		TimeoutExpress: strconv.FormatInt(requestTimeout, 10),
		Timestamp:      strconv.FormatInt(timestamp, 10),
		TotalAmount:    strconv.FormatFloat(totalAmount, 'f', 2, 64),
	}

	subscriptionTransaction := &entity.SubscriptionTransaction{
		UserID:         userID,
		PlanID:         planID,
		AppID:          telebirrRequestParameters.AppID,
		ReceiverName:   telebirrRequestParameters.ReceiverName,
		Subject:        telebirrRequestParameters.Subject,
		ReceivedAmount: receivedAmount,
		TransactionFee: transactionFee,
		CurrencyType:   currencyType,
		TimeoutExpress: requestTimeout,
		Nonce:          telebirrRequestParameters.Nonce,
		OutTradeNo:     telebirrRequestParameters.OutTradeNo,
		InitiatedFrom:  initiatedFrom,
	}

	// Checking the nonce and out_trade_no uniqueness
	for errMap := service.ValidateSubscriptionTransaction(subscriptionTransaction); errMap["nonce"] != nil ||
		errMap["out_trade_no"] != nil; {

		uniqueID := strings.ReplaceAll(uuid.Must(uuid.NewRandom()).String(), "-", "")
		subscriptionTransaction.Nonce = uniqueID
		subscriptionTransaction.OutTradeNo = "T_" + uniqueID
	}

	telebirrParameterS, err := json.Marshal(telebirrRequestParameters)
	if err != nil {
		return "", err
	}

	publicKey, err := tools.BytesToPublicKey(service.TelebirrAPI.PublicKey, nil)
	if err != nil {
		return "", err
	}

	var encryptedBytes []byte
	source := string(telebirrParameterS)
	for len(source) > 0 {
		input := tools.Substr(source, 0, 117)
		source = tools.Substr(source, 117, len(source))

		encryptedChunkBytes, err := tools.EncryptWithPublicKey([]byte(input), publicKey)
		if err != nil {
			return "", err
		}

		encryptedBytes = append(encryptedBytes, encryptedChunkBytes...)
	}
	encryptedValue := base64.URLEncoding.EncodeToString(encryptedBytes)

	urlValues, err := query.Values(telebirrRequestParameters)
	if err != nil {
		return "", errors.New("unable to construct valid request")
	}

	// Incase
	parsedURLValue, err := url.PathUnescape(urlValues.Encode())
	if err != nil {
		return "", errors.New("unable to construct valid request")
	}

	hasher := sha256.New()
	hasher.Write([]byte(parsedURLValue))
	signedValue := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	telebirrRequest := &TelebirrRequest{AppID: service.TelebirrAPI.AppID, USSD: encryptedValue, Sign: signedValue}

	client := new(http.Client)
	jsonOutput, _ := json.MarshalIndent(telebirrRequest, "", "\t")
	output := bytes.NewBuffer(jsonOutput)
	url := service.TelebirrAPI.AccessPoint + "toTradeWebPay"

	request, err := http.NewRequest("POST", url, output)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	telebirrResponse := new(TelebirrResponse)
	err = json.Unmarshal(responseBody, telebirrResponse)
	if err != nil {
		return "", err
	}

	// Means error has occurred
	if telebirrResponse.Code != "0" {
		return "", errors.New("unable to generate web url")
	}

	// Adding the subscription transaction to the database
	subscriptionTransaction.Status = entity.TransactionStatusPending
	err = service.AddSubscriptionTransaction(subscriptionTransaction)
	if err != nil {
		return "", err
	}

	return telebirrResponse.Data["toPayUrl"], nil
}
