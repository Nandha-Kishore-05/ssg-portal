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

// In-memory store to hold OTPs temporarily (key: email, value: OTP)
var otpStore = make(map[string]string)

// SMTP settings (replace these with your actual values)
const (
    smtpHost     = "smtp.gmail.com"
    smtpPort     = "587"
    smtpUser     = "nandhakishore.ct23@bitsathy.ac.in"
    smtpPassword = "awho onwp dxkb ydbo"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
    rand.Seed(time.Now().UnixNano())
    otp := fmt.Sprintf("%06d", rand.Intn(1000000))
    return otp
}

// Check if the email exists in the master_user table
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

// sendEmail sends an email with the provided subject and body to the given recipient
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

// SendOTPHandler handles sending OTP to the user (via email)
func SendOTPHandler(c *gin.Context) {
    var request struct {
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Check if the email exists in the database
    if !emailExists(request.Email) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
        return
    }

    // Generate a new OTP and store it temporarily in memory
    otp := GenerateOTP()
    otpStore[request.Email] = otp

    // Send the OTP via email
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

// VerifyOTPHandler handles OTP verification
func VerifyOTPHandler(c *gin.Context) {
    var request struct {
        Email string `json:"email"`
        Otp   string `json:"otp"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Check if the OTP exists for the email
    storedOtp, exists := otpStore[request.Email]
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not found for this email"})
        return
    }

    // Verify if the OTP matches
    if storedOtp != request.Otp {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
        return
    }

    // If OTP is correct, remove it from the store (to ensure it's only used once)
    delete(otpStore, request.Email)

    // Retrieve user info from the database (assuming you fetch the user after OTP verification)
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

    // Retrieve the auth token from the user_login table
    var authToken string
    authQuery := "SELECT auth_token FROM user_login WHERE user_id = ?"
    err = config.Database.QueryRow(authQuery, user.ID).Scan(&authToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve auth token"})
        return
    }

    // Return the auth token and user details
    c.JSON(http.StatusOK, gin.H{

        "auth_token": authToken,
        "user":       user,
    })
}