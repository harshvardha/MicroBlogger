package cache

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"sync"
	"time"
)

type emailConfig struct {
	fromEmail   string
	subject     string
	body        string
	smtpHost    string
	smtpPort    string
	appPassword string
}

// struct to store the otp cache data
type otpCacheData struct {
	otp      string
	issuedAt time.Time
}

// otp cache struct
type OtpCache struct {
	cache              map[string]otpCacheData
	lock               sync.Mutex
	expiresAfter       time.Duration
	resendAllowedAfter time.Duration
	emailConfig        *emailConfig
}

func NewOTPCache(fromEmail string, subject string, body string, smtpHost string, smtpPort string, appPassword string) *OtpCache {
	return &OtpCache{
		cache:              make(map[string]otpCacheData),
		expiresAfter:       1 * time.Minute,
		resendAllowedAfter: 2 * time.Minute,
		emailConfig: &emailConfig{
			fromEmail:   fromEmail,
			subject:     "Registration OTP",
			smtpHost:    smtpHost,
			smtpPort:    smtpPort,
			appPassword: appPassword,
		},
	}
}

func (otpCache *OtpCache) set(verificationToken string, otp string) (bool, error) {
	otpCache.lock.Lock()
	defer otpCache.lock.Unlock()

	// checking if all the info is provided or not
	if verificationToken == "" || otp == "" {
		return false, errors.New("incomplete value to store in cache")
	}

	otpCache.cache[verificationToken] = otpCacheData{
		otp:      otp,
		issuedAt: time.Now(),
	}

	return true, nil
}

func (otpCache *OtpCache) get(verificationToken string) (otpCacheData, error) {
	otpCache.lock.Lock()
	defer otpCache.lock.Unlock()

	// checking if the verificaiton token is empty
	if len(verificationToken) == 0 {
		return otpCacheData{}, errors.New("invalid verification token")
	}

	data, exists := otpCache.cache[verificationToken]
	if !exists {
		return data, errors.New("invalid verification token: data not found")
	}

	return data, nil
}

func (otpCache *OtpCache) delete(verificationToken string) (bool, error) {
	otpCache.lock.Lock()
	defer otpCache.lock.Unlock()

	// checking if the verification token is valid
	if len(verificationToken) == 0 {
		return false, errors.New("invalid verification token: empty")
	}

	if _, exists := otpCache.cache[verificationToken]; exists {
		delete(otpCache.cache, verificationToken)
		return true, nil
	} else {
		return false, errors.New("invalid verification token: data not found")
	}
}

func generateOTPAndVerificationToken() (string, string, error) {
	// generating verification token
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("Error generating otp verification token: ", err)
		return "", "", errors.New("error generating otp verification token")
	}

	verificationToken := base64.StdEncoding.EncodeToString(randomBytes)

	// generating 6-digit otp
	buffer := make([]byte, 6)
	_, err = io.ReadFull(rand.Reader, buffer)
	if err != nil {
		log.Println("Error Generating OTP: ", err)
		return "", "", errors.New("error generating otp")
	}
	for i := range buffer {
		buffer[i] = buffer[i]%10 + '0'
	}

	return verificationToken, string(buffer), nil
}

func sendMail(emailConfig *emailConfig, otp string, to []string) error {
	// constructing the email message
	emailConfig.body += otp
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", emailConfig.fromEmail, to[0], emailConfig.subject, emailConfig.body)

	// authentication
	auth := smtp.PlainAuth("", emailConfig.fromEmail, emailConfig.appPassword, emailConfig.smtpHost)

	// sending the mail
	err := smtp.SendMail(emailConfig.smtpHost+":"+emailConfig.smtpPort, auth, emailConfig.fromEmail, to, []byte(message))
	if err != nil {
		log.Println("Error Sending OTP: ", err)
		return err
	}

	return nil
}

func (otpCache *OtpCache) SendOTP(to string) error {
	verificationToken, otp, err := generateOTPAndVerificationToken()
	if err != nil {
		return err
	}

	err = sendMail(otpCache.emailConfig, otp, []string{to})
	if err != nil {
		return err
	}
	otpCache.set(verificationToken, otp)

	return nil
}

func (otpCache *OtpCache) VerifyOTP(verificationToken string, otp string) error {
	// checking if the verification token is empty
	if len(verificationToken) == 0 {
		return errors.New("invalid verification token: empty")
	}

	// checking if the data exist
	data, err := otpCache.get(verificationToken)
	if err != nil {
		return err
	}

	// checking if the otp is expired or not
	if time.Now().After(data.issuedAt.Add(otpCache.expiresAfter)) {
		otpCache.delete(verificationToken)
		return errors.New("otp expired")
	}

	// checking if the otp is correct or not
	if data.otp != otp {
		return errors.New("incorrect otp")
	}
	otpCache.delete(verificationToken)

	return nil
}

func (otpCache *OtpCache) ResendOTP(verificationToken string, email string) error {
	otpCache.delete(verificationToken)
	verificationToken, otp, err := generateOTPAndVerificationToken()
	if err != nil {
		return err
	}

	err = sendMail(otpCache.emailConfig, otp, []string{email})
	if err != nil {
		return err
	}
	otpCache.set(verificationToken, otp)

	return nil
}

func (otpCache *OtpCache) IsResendAllowed(verificationToken string) (bool, error) {
	// checking if verification token is empty or not
	if len(verificationToken) == 0 {
		return false, errors.New("invalid verification token: empty")
	}

	// checking if the data exist
	data, err := otpCache.get(verificationToken)
	if err != nil {
		return false, err
	}

	// checking if resend is allowed or not
	if time.Now().Before(data.issuedAt.Add(otpCache.resendAllowedAfter)) {
		return false, errors.New(fmt.Sprintf("resend allowed after %f", otpCache.resendAllowedAfter.Seconds()-float64(time.Now().Sub(data.issuedAt))))
	}

	return true, nil
}
