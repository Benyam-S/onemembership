package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/feedback"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/serviceprovider"
	"github.com/Benyam-S/onemembership/user"
)

// Service is a type that defines a feedback service
type Service struct {
	feedbackRepo        feedback.IFeedbackRepository
	serviceProviderRepo serviceprovider.IServiceProviderRepository
	userRepo            user.IUserRepository
	logger              *log.Logger
}

// NewFeedbackService is a function that returns a new feedback service
func NewFeedbackService(feedbackRepository feedback.IFeedbackRepository,
	serviceProviderRepository serviceprovider.IServiceProviderRepository,
	userRepository user.IUserRepository, feedbackLogger *log.Logger) feedback.IService {
	return &Service{feedbackRepo: feedbackRepository, serviceProviderRepo: serviceProviderRepository,
		userRepo: userRepository, logger: feedbackLogger}
}

// AddFeedback is a method that adds a new feedback to the system
func (service *Service) AddFeedback(newFeedback *entity.Feedback) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started feedback adding process, Feedback => %s", newFeedback.ToString()),
		service.logger.Logs.ServerLogFile)

	err := service.feedbackRepo.Create(newFeedback)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Feedback => %s, %s",
			newFeedback.ToString(), err.Error()))

		return errors.New("unable to add new feedback")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished feedback adding process, Feedback => %s", newFeedback.ToString()),
		service.logger.Logs.ServerLogFile)

	return nil
}

// ValidateFeedback is a method that validates a feedback entries.
// It checks if the feedback has a valid entries or not and return map of errors if any.
func (service *Service) ValidateFeedback(feedback *entity.Feedback) entity.ErrMap {

	errMap := make(map[string]error)

	emptyComment, _ := regexp.MatchString(`^\s*$`, feedback.Comment)
	if emptyComment {
		errMap["comment"] = errors.New("comment can not be empty")
	} else if len(feedback.Comment) > 1000 {
		errMap["comment"] = errors.New("comment can not exceed 1000 characters")
	}

	if feedback.ClientID != "" {
		_, err1 := service.userRepo.Find(feedback.ClientID)
		_, err2 := service.serviceProviderRepo.Find(feedback.ClientID)
		if err1 != nil && err2 != nil {
			errMap["client_id"] = errors.New("no client found for the provided client id")
		}
	} else {
		errMap["client_id"] = errors.New("no client found for the provided client id")
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindFeedback is a method that find and return a feedback that matches the id value
func (service *Service) FindFeedback(id string) (*entity.Feedback, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Feedback finding process { ID : %s }", id),
		service.logger.Logs.ServerLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no feedback found")
	}

	feedback, err := service.feedbackRepo.Find(id)
	if err != nil {
		return nil, errors.New("no feedback found")
	}
	return feedback, nil
}

// FindMultipleFeedbacks is a method that find and return multiple feedbacks that matchs the clientID value
func (service *Service) FindMultipleFeedbacks(clientID string) []*entity.Feedback {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple feedbacks finding process { Client ID : %s }", clientID),
		service.logger.Logs.ServerLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, clientID)
	if empty {
		return []*entity.Feedback{}
	}

	return service.feedbackRepo.FindMultiple(clientID)
}

// AllFeedbacks is a method that returns all the feedbacks with pagination
func (service *Service) AllFeedbacks(status string, pageNum int64) ([]*entity.Feedback, int64) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Getting all feedbacks process { Status : %s, Page Number : %d }", status, pageNum),
		service.logger.Logs.ServerLogFile)

	var seenStatus int64
	if status == entity.FeedbackUnseen {
		seenStatus = 0
	} else if status == entity.FeedbackSeen {
		seenStatus = 1
	} else {
		seenStatus = 2
	}

	return service.feedbackRepo.FindAll(seenStatus, pageNum)
}

// SearchFeedbacks is a method that searchs and returns a set of feedbacks related to the key identifier
func (service *Service) SearchFeedbacks(key, status string, pageNum int64, extra ...string) ([]*entity.Feedback, int64) {
	/* ---------------------------- Logging ---------------------------- */
	extraLog := ""
	for index, extraValue := range extra {
		extraLog += fmt.Sprintf(", Extra%d : %s", index, extraValue)
	}
	service.logger.Log(fmt.Sprintf("Searching feedbacks process { Key : %s, Status : %s, Page Number : %d%s }",
		key, status, pageNum, extraLog), service.logger.Logs.ServerLogFile)

	var seenStatus int64
	if status == entity.FeedbackUnseen {
		seenStatus = 0
	} else if status == entity.FeedbackSeen {
		seenStatus = 1
	} else {
		seenStatus = 2
	}

	defaultSearchColumnsRegx := []string{"comment"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "client_id"}

	result1 := make([]*entity.Feedback, 0)
	result2 := make([]*entity.Feedback, 0)
	results := make([]*entity.Feedback, 0)
	resultsMap := make(map[string]*entity.Feedback)
	var pageCount1 int64 = 0
	var pageCount2 int64 = 0
	var pageCount int64 = 0

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results, 0
	}

	result1, pageCount1 = service.feedbackRepo.Search(key, seenStatus, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2, pageCount2 = service.feedbackRepo.SearchWRegx(key, seenStatus, pageNum, defaultSearchColumnsRegx...)
	}

	for _, feedback := range result1 {
		resultsMap[feedback.ID] = feedback
	}

	for _, feedback := range result2 {
		resultsMap[feedback.ID] = feedback
	}

	for _, uniqueFeedback := range resultsMap {
		results = append(results, uniqueFeedback)
	}

	pageCount = pageCount1
	if pageCount < pageCount2 {
		pageCount = pageCount2
	}

	return results, pageCount
}

// MarkAsSeen is a method that mark a feedback as seen
func (service *Service) MarkAsSeen(feedbackID string) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started marking feedback as seen process { ID : %s }",
		feedbackID), service.logger.Logs.ServerLogFile)

	feedback, err := service.feedbackRepo.Find(feedbackID)
	if err != nil {
		return errors.New("feedback not found")
	}

	if feedback.Seen {
		return errors.New("unable to perform operation")
	}

	feedback.Seen = true
	err = service.feedbackRepo.Update(feedback)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For marking feedback as seen { ID : %s }, %s",
			feedbackID, err.Error()))

		return errors.New("unable to update feedback")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished marking feedback as seen process, Feedback => %s",
		feedback.ToString()), service.logger.Logs.ServerLogFile)

	return nil
}

// SetFeedbackClientIDNull is a method that set the client id to null for all feedbacks that matches the given client id
func (service *Service) SetFeedbackClientIDNull(clientID string) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started feedbacks' client id erasing process { Client ID : %s }", clientID),
		service.logger.Logs.ServerLogFile)

	err := service.feedbackRepo.SetToNull(clientID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For feedbacks' client id erasing { Client ID : %s }, %s",
			clientID, err.Error()))

		return errors.New("unable to erase feedbacks client id")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished feedbacks' client id erasing process, { Client ID : %s }",
		clientID), service.logger.Logs.ServerLogFile)

	return nil
}

// DeleteFeedback is a method that deletes a feedback from the system using an id
func (service *Service) DeleteFeedback(id string) (*entity.Feedback, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started feedback deleting process { ID : %s }", id),
		service.logger.Logs.ServerLogFile)

	feedback, err := service.feedbackRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting feedback { ID : %s }, %s", id, err.Error()))

		return nil, errors.New("unable to delete feedback")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished feedback deleting process, Deleted Feedback => %s",
		feedback.ToString()), service.logger.Logs.ServerLogFile)

	return feedback, nil
}

// DeleteMultipleFeedbacks is a method that deletes multiple feedbacks from the system that match the given clientID
func (service *Service) DeleteMultipleFeedbacks(clientID string) []*entity.Feedback {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple feedbacks deleting process { Client ID : %s }",
		clientID), service.logger.Logs.ServerLogFile)

	return service.feedbackRepo.DeleteMultiple(clientID)
}
