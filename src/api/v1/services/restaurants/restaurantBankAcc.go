package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	postgres "dine-server/src/config/database"
	"dine-server/src/config/env"
	"dine-server/src/config/payments"
	models_restaurant "dine-server/src/models/restaurants"

	"github.com/gin-gonic/gin"
)

// ConnectRestaurantBankAccount connects a restaurant's bank account
// @Summary Connect a restaurant's bank account
// @Description Connect a restaurant's bank account
// @Tags Restaurant Bank Account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models_restaurant.AddBankAccount true "Bank account data"
// @Router /api/v1/restaurants/bank-account [post]
func ConnectRestaurantBankAccount(c *gin.Context) {
	var BankAccountData models_restaurant.AddBankAccount

	// Validate input data
	if err := c.ShouldBindJSON(&BankAccountData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Create contact on Razorpay
	contactID, err := createRazorpayContact(BankAccountData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create contact: " + err.Error()})
		return
	}

	// Create fund account for the contact
	fundAccountID, err := createFundAccount(contactID, BankAccountData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create fund account: " + err.Error()})
		return
	}

	// Save bank account details in the database
	BankAccount := models_restaurant.RestaurantBankAccount{
		RestaurantID:  BankAccountData.RestaurantID,
		ContactID:     contactID,
		Email:         BankAccountData.Email,
		Phone:         BankAccountData.Phone,
		FundAccID:     fundAccountID,
		BankName:      BankAccountData.BankName,
		AccountName:   BankAccountData.AccountName,
		AccountNumber: BankAccountData.AccountNumber,
		IFSCCode:      BankAccountData.IFSCCode,
		Branch:        BankAccountData.Branch,
	}

	if err := postgres.DB.Create(&BankAccount).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save bank account details: " + err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":     "Bank account connected successfully",
		"bankAccount": BankAccount,
	})
}

// createRazorpayContact creates a contact on Razorpay and returns the contact ID
func createRazorpayContact(data models_restaurant.AddBankAccount) (string, error) {
	contactPayload := map[string]interface{}{
		"name":         data.AccountName,
		"email":        data.Email,
		"contact":      data.Phone,
		"type":         "customer", // Adjust type as necessary
		"reference_id": data.RestaurantID.String(),
	}

	payloadBytes, _ := json.Marshal(contactPayload)
	req, _ := http.NewRequest("POST", "https://api.razorpay.com/v1/contacts", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(env.PaymentsVar["RAZORPAY_KEY_ID"], env.PaymentsVar["RAZORPAY_SECRET_KEY"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(bodyBytes))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", errors.New("failed to decode Razorpay response")
	}

	contactID, ok := response["id"].(string)
	if !ok {
		return "", errors.New("missing contact ID in Razorpay response")
	}

	return contactID, nil
}

// createFundAccount creates a fund account for the given contact ID and returns the fund account ID
func createFundAccount(contactID string, data models_restaurant.AddBankAccount) (string, error) {
	fundAccountPayload := map[string]interface{}{
		"contact_id":   contactID,
		"account_type": "bank_account",
		"bank_account": map[string]interface{}{
			"name":           data.AccountName,
			"account_number": data.AccountNumber,
			"ifsc":           data.IFSCCode,
		},
	}

	fundAccount, err := payments.RazorpayClient.FundAccount.Create(fundAccountPayload, nil)
	if err != nil {
		return "", err
	}

	fundAccountID, ok := fundAccount["id"].(string)
	if !ok {
		return "", errors.New("missing fund account ID in Razorpay response")
	}

	return fundAccountID, nil
}
