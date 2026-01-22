package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"xiangshoufu/internal/models"
	"xiangshoufu/internal/repository"
)

// AlertService 告警服务
type AlertService struct {
	configRepo      repository.AlertConfigRepository
	logRepo         repository.AlertLogRepository
	failCounterRepo repository.JobFailCounterRepository
	httpClient      *http.Client
}

// NewAlertService 创建告警服务
func NewAlertService(
	configRepo repository.AlertConfigRepository,
	logRepo repository.AlertLogRepository,
	failCounterRepo repository.JobFailCounterRepository,
) *AlertService {
	return &AlertService{
		configRepo:      configRepo,
		logRepo:         logRepo,
		failCounterRepo: failCounterRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendAlert 发送告警（根据配置发送到所有启用的通道）
func (s *AlertService) SendAlert(req *models.AlertRequest) error {
	configs, err := s.configRepo.FindEnabled()
	if err != nil {
		return fmt.Errorf("find alert configs failed: %w", err)
	}

	if len(configs) == 0 {
		log.Printf("[AlertService] No enabled alert configs, skip sending")
		return nil
	}

	var lastErr error
	for _, config := range configs {
		// 创建告警日志
		alertLog := &models.AlertLog{
			JobName:     req.JobName,
			AlertType:   req.AlertType,
			ChannelType: config.ChannelType,
			ConfigID:    &config.ID,
			Title:       req.Title,
			Message:     req.Message,
			SendStatus:  models.AlertSendStatusPending,
			CreatedAt:   time.Now(),
		}

		if err := s.logRepo.Create(alertLog); err != nil {
			log.Printf("[AlertService] Create alert log failed: %v", err)
			continue
		}

		// 发送告警
		var sendErr error
		switch config.ChannelType {
		case models.AlertChannelDingTalk:
			sendErr = s.sendDingTalkAlert(config, req)
		case models.AlertChannelWeChatWork:
			sendErr = s.sendWeChatWorkAlert(config, req)
		case models.AlertChannelEmail:
			sendErr = s.sendEmailAlert(config, req)
		default:
			sendErr = fmt.Errorf("unknown channel type: %d", config.ChannelType)
		}

		// 更新发送状态
		if sendErr != nil {
			s.logRepo.UpdateSendStatus(alertLog.ID, models.AlertSendStatusFailed, sendErr.Error())
			lastErr = sendErr
			log.Printf("[AlertService] Send alert failed: channel=%s, err=%v",
				models.GetAlertChannelName(config.ChannelType), sendErr)
		} else {
			s.logRepo.UpdateSendStatus(alertLog.ID, models.AlertSendStatusSent, "")
			log.Printf("[AlertService] Alert sent successfully: channel=%s, job=%s",
				models.GetAlertChannelName(config.ChannelType), req.JobName)
		}
	}

	return lastErr
}

// SendJobFailAlert 发送任务失败告警
func (s *AlertService) SendJobFailAlert(jobName string, errMsg string, retryCount int) error {
	title := fmt.Sprintf("【任务失败】%s", jobName)
	message := s.buildJobFailMessage(jobName, errMsg, retryCount)

	return s.SendAlert(&models.AlertRequest{
		JobName:      jobName,
		AlertType:    models.AlertTypeJobFailed,
		Title:        title,
		Message:      message,
		ErrorMessage: errMsg,
	})
}

// SendConsecutiveFailAlert 发送连续失败告警
func (s *AlertService) SendConsecutiveFailAlert(jobName string, consecutiveFails int) error {
	title := fmt.Sprintf("【连续失败告警】%s 已连续失败 %d 次", jobName, consecutiveFails)
	message := s.buildConsecutiveFailMessage(jobName, consecutiveFails)

	return s.SendAlert(&models.AlertRequest{
		JobName:   jobName,
		AlertType: models.AlertTypeConsecutiveFail,
		Title:     title,
		Message:   message,
	})
}

// SendJobTimeoutAlert 发送任务超时告警
func (s *AlertService) SendJobTimeoutAlert(jobName string, duration time.Duration) error {
	title := fmt.Sprintf("【任务超时】%s", jobName)
	message := s.buildJobTimeoutMessage(jobName, duration)

	return s.SendAlert(&models.AlertRequest{
		JobName:   jobName,
		AlertType: models.AlertTypeJobTimeout,
		Title:     title,
		Message:   message,
	})
}

// CheckAndSendConsecutiveFailAlert 检查并发送连续失败告警
func (s *AlertService) CheckAndSendConsecutiveFailAlert(jobName string, threshold int) error {
	counter, err := s.failCounterRepo.FindByJobName(jobName)
	if err != nil {
		return nil // 找不到计数器，不需要告警
	}

	// 检查是否达到告警阈值
	if counter.ConsecutiveFails >= threshold {
		// 检查是否已经告警过（避免重复告警）
		if counter.LastAlertAt != nil {
			// 如果最后告警时间在最后失败时间之后，说明已经告警过
			if counter.LastFailAt != nil && counter.LastAlertAt.After(*counter.LastFailAt) {
				return nil
			}
		}

		// 发送告警
		if err := s.SendConsecutiveFailAlert(jobName, counter.ConsecutiveFails); err != nil {
			return err
		}

		// 更新最后告警时间
		return s.failCounterRepo.UpdateLastAlert(jobName)
	}

	return nil
}

// TestAlert 测试告警配置
func (s *AlertService) TestAlert(configID int64) error {
	config, err := s.configRepo.FindByID(configID)
	if err != nil {
		return fmt.Errorf("find config failed: %w", err)
	}

	req := &models.AlertRequest{
		JobName:   "TestJob",
		AlertType: models.AlertTypeJobFailed,
		Title:     "【测试告警】这是一条测试消息",
		Message:   fmt.Sprintf("告警配置测试\n配置名称: %s\n通道类型: %s\n发送时间: %s",
			config.Name,
			models.GetAlertChannelName(config.ChannelType),
			time.Now().Format("2006-01-02 15:04:05")),
	}

	switch config.ChannelType {
	case models.AlertChannelDingTalk:
		return s.sendDingTalkAlert(config, req)
	case models.AlertChannelWeChatWork:
		return s.sendWeChatWorkAlert(config, req)
	case models.AlertChannelEmail:
		return s.sendEmailAlert(config, req)
	default:
		return fmt.Errorf("unknown channel type: %d", config.ChannelType)
	}
}

// sendDingTalkAlert 发送钉钉告警
func (s *AlertService) sendDingTalkAlert(config *models.AlertConfig, req *models.AlertRequest) error {
	if config.WebhookURL == "" {
		return fmt.Errorf("dingtalk webhook url is empty")
	}

	// 构建请求URL（处理签名）
	webhookURL := config.WebhookURL
	if config.WebhookSecret != "" {
		timestamp := time.Now().UnixMilli()
		sign := s.generateDingTalkSign(timestamp, config.WebhookSecret)
		if strings.Contains(webhookURL, "?") {
			webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhookURL, timestamp, url.QueryEscape(sign))
		} else {
			webhookURL = fmt.Sprintf("%s?timestamp=%d&sign=%s", webhookURL, timestamp, url.QueryEscape(sign))
		}
	}

	// 构建消息体
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": req.Title,
			"text":  fmt.Sprintf("## %s\n\n%s", req.Title, req.Message),
		},
	}

	return s.sendWebhook(webhookURL, payload)
}

// sendWeChatWorkAlert 发送企业微信告警
func (s *AlertService) sendWeChatWorkAlert(config *models.AlertConfig, req *models.AlertRequest) error {
	if config.WebhookURL == "" {
		return fmt.Errorf("wechat work webhook url is empty")
	}

	// 构建消息体
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": fmt.Sprintf("## %s\n\n%s", req.Title, req.Message),
		},
	}

	return s.sendWebhook(config.WebhookURL, payload)
}

// sendEmailAlert 发送邮件告警
func (s *AlertService) sendEmailAlert(config *models.AlertConfig, req *models.AlertRequest) error {
	if config.EmailAddresses == "" {
		return fmt.Errorf("email addresses is empty")
	}
	if config.EmailSMTPHost == "" {
		return fmt.Errorf("smtp host is empty")
	}

	// 解析收件人列表
	to := strings.Split(config.EmailAddresses, ",")
	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}

	// 构建邮件内容
	subject := req.Title
	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"<html><body>"+
		"<h2>%s</h2>"+
		"<pre>%s</pre>"+
		"<p>发送时间: %s</p>"+
		"</body></html>",
		req.Title,
		strings.ReplaceAll(req.Message, "\n", "<br>"),
		time.Now().Format("2006-01-02 15:04:05"))

	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"%s",
		strings.Join(to, ","),
		subject,
		body)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.EmailSMTPHost, config.EmailSMTPPort)
	auth := smtp.PlainAuth("", config.EmailUsername, config.EmailPassword, config.EmailSMTPHost)

	return smtp.SendMail(addr, auth, config.EmailUsername, to, []byte(msg))
}

// sendWebhook 发送Webhook请求
func (s *AlertService) sendWebhook(webhookURL string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload failed: %w", err)
	}

	resp, err := s.httpClient.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("send webhook failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook response status: %d", resp.StatusCode)
	}

	// 解析响应检查是否成功
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil // 忽略解析错误，认为发送成功
	}

	// 钉钉返回 errcode
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg, _ := result["errmsg"].(string)
		return fmt.Errorf("dingtalk error: %s", errmsg)
	}

	// 企微返回 errcode
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg, _ := result["errmsg"].(string)
		return fmt.Errorf("wechat work error: %s", errmsg)
	}

	return nil
}

// generateDingTalkSign 生成钉钉签名
func (s *AlertService) generateDingTalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// buildJobFailMessage 构建任务失败消息
func (s *AlertService) buildJobFailMessage(jobName string, errMsg string, retryCount int) string {
	return fmt.Sprintf(`**任务名称**: %s
**错误信息**: %s
**重试次数**: %d
**发生时间**: %s
**服务器**: %s`,
		jobName,
		errMsg,
		retryCount,
		time.Now().Format("2006-01-02 15:04:05"),
		getHostname())
}

// buildConsecutiveFailMessage 构建连续失败消息
func (s *AlertService) buildConsecutiveFailMessage(jobName string, consecutiveFails int) string {
	return fmt.Sprintf(`**任务名称**: %s
**连续失败次数**: %d
**告警级别**: 严重
**发生时间**: %s
**服务器**: %s

> 请及时处理，避免业务受影响！`,
		jobName,
		consecutiveFails,
		time.Now().Format("2006-01-02 15:04:05"),
		getHostname())
}

// buildJobTimeoutMessage 构建任务超时消息
func (s *AlertService) buildJobTimeoutMessage(jobName string, duration time.Duration) string {
	return fmt.Sprintf(`**任务名称**: %s
**执行时长**: %s
**告警类型**: 超时告警
**发生时间**: %s
**服务器**: %s`,
		jobName,
		duration.String(),
		time.Now().Format("2006-01-02 15:04:05"),
		getHostname())
}

// getHostname 获取主机名
func getHostname() string {
	// 简单实现，实际可以使用 os.Hostname()
	return "xiangshoufu-server"
}
