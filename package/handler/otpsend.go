package handler

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

// ====================== OTP Generation ===================================
func GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// ====================== Sending OTP to User Mail =========================
func SendOtp(email, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "nuhmotp@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code for Signup")
	m.SetBody("text/plain", "Your OTP for signup is: "+otp)

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("APPEMAIL"), os.Getenv("APPPASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("--------------", err, "------------------")
		return err
	}
	return nil

}

