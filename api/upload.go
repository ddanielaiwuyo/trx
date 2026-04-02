package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var logger = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

func UploadBankStatement(c *gin.Context) {
	logger.Println("Client uploaded a file")
	fileHeaders, err := c.FormFile("bankStatement")
	if err != nil {
		logger.Println("Could not get bankStatement: ", err)
		c.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	if fileHeaders == nil {
		logger.Println("No file provided: ", fileHeaders)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	bankName := c.PostForm("bankName")
	logger.Println("bankType: ", bankName)

	file, err := fileHeaders.Open()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No MIME Header was provided"})
		return
	}

	f := File{F: file, BankName: bankName}
	bank := getBank(f)
	if bank == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bank not supported"})
		return
	}

	res, err := bank.GetCanonData()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func getBank(f File) Bank {
	var bank Bank
	switch f.BankName {
	case "Monzo":
		bank = &MonzoBank{File: f}
	}

	return bank
}
