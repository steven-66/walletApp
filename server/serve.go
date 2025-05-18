package server

import (
	"context"
	"fmt"
	"os"
	"time"
	"walletApp/config"
	"walletApp/dto"
	"walletApp/server/handler"
)

type App struct {
	BalanceHandler     *handler.BalanceHandler
	TransactionHandler *handler.TransactionHandler
}

func NewApp() *App {
	config.InitDB()
	app := &App{
		BalanceHandler:     handler.NewBalanceHandler(),
		TransactionHandler: handler.NewTransactionHandler(),
	}

	return app
}

func (a *App) Start() {
	for {
		fmt.Println("\nWallet App CLI")
		fmt.Println("----------")
		fmt.Println("1. Deposit Money")
		fmt.Println("2. Withdraw Money")
		fmt.Println("3. Check Balance")
		fmt.Println("4. View Transaction History")
		fmt.Println("5. Transfer Money")
		fmt.Println("6. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		ctx := context.Background() // Create a context for each request

		switch choice {
		case 1:
			fmt.Print("Enter user ID: ")
			var userID uint
			fmt.Scan(&userID)
			fmt.Print("Enter amount to deposit: ")
			var amount float64
			fmt.Scan(&amount)
			resp, err := a.BalanceHandler.Deposit(ctx, &dto.DepositRequest{
				UserID: userID,
				Amount: amount,
			})
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Deposit successful!")
				fmt.Printf("New Balance: %.2f\n", resp.Balance)
			}
		case 2:
			fmt.Print("Enter user ID: ")
			var userID uint
			fmt.Scan(&userID)
			fmt.Print("Enter amount to withdraw: ")
			var amount float64
			fmt.Scan(&amount)
			newBalance, err := a.BalanceHandler.Withdraw(ctx, &dto.WithdrawRequest{
				UserID: userID,
				Amount: amount,
			})
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Withdrawal successful!")
				fmt.Printf("New Balance: %.2f\n", newBalance.Balance)
			}
		case 3:
			fmt.Print("Enter user ID: ")
			var userID uint
			fmt.Scan(&userID)
			balance, err := a.BalanceHandler.CheckBalance(ctx, userID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Balance fetched successfully!")
				fmt.Printf("Balance: %.2f\n", balance)
			}
		case 4:
			fmt.Print("Enter user ID: ")
			var userID uint
			fmt.Scan(&userID)
			resp, err := a.TransactionHandler.ViewTransactionHistory(ctx, userID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Transaction History:")
				fmt.Println("--------------------")
				for _, transaction := range resp.Transactions {
					fmt.Printf("%s: %.2f at %s\n", transaction.Type, transaction.Amount, transaction.Timestamp.Format("2006-01-02 15:04:05"))
				}
			}
		case 5:
			fmt.Print("Enter sender user ID: ")
			var fromUserID uint
			fmt.Scan(&fromUserID)
			fmt.Print("Enter recipient user ID: ")
			var toUserID uint
			fmt.Scan(&toUserID)
			fmt.Print("Enter amount to transfer: ")
			var amount float64
			fmt.Scan(&amount)
			response, err := a.BalanceHandler.Transfer(ctx, &dto.TransferRequest{
				FromUserID: fromUserID,
				ToUserID:   toUserID,
				Amount:     amount,
			})
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(response.Message)
				if response.Success {
					data := response.Data
					fmt.Printf("Sender's New Balance: %.2f\n", data["sender_balance"])
					fmt.Printf("Recipient's New Balance: %.2f\n", data["recipient_balance"])
				}
			}
		case 6:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		countdownToMainMenu()
	}
}

// Countdown function to return to the main menu
func countdownToMainMenu() {
	for i := 3; i > 0; i-- {
		fmt.Printf("\nReturning to the main menu in: %d seconds\n", i)
		time.Sleep(1 * time.Second)
	}
}
