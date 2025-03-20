package otp

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"ssg-portal/config"
	"time"

	"github.com/gin-gonic/gin"
)

var otpStore = make(map[string]string)

const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUser     = "nandhakishore.ct23@bitsathy.ac.in"
	smtpPassword = "auoa atyq ymid tmgn"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}

func emailExists(email string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM master_user WHERE email = ?)"
	err := config.Database.QueryRow(query, email).Scan(&exists)
	if err != nil {
		fmt.Println("Database error:", err)
		return false
	}
	return exists
}

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}

func SendOTPHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !emailExists(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	otp := GenerateOTP()
	otpStore[request.Email] = otp

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	err := sendEmail(request.Email, subject, body)
	if err != nil {
		fmt.Println("Error sending email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email"})
}

func VerifyOTPHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	storedOtp, exists := otpStore[request.Email]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not found for this email"})
		return
	}

	if storedOtp != request.Otp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	delete(otpStore, request.Email)

	var user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	query := "SELECT id, name FROM master_user WHERE email = ?"
	err := config.Database.QueryRow(query, request.Email).Scan(&user.ID, &user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	var authToken string
	authQuery := "SELECT auth_token FROM user_login WHERE user_id = ?"
	err = config.Database.QueryRow(authQuery, user.ID).Scan(&authToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve auth token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"auth_token": authToken,
		"user":       user,
	})
}
