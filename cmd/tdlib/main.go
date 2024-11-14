package main



import (
	"fmt"

	"github.com/Arman92/go-tdlib"
)
const (
	AppID    = "4"
	AppHash  = "014b35b6184100b085b0d0572f9b5103"
)
func main(){
	tdlib.SetLogVerbosityLevel(1)

	client := tdlib.NewClient(tdlib.Config{
		APIID:               AppID,
		APIHash:             AppHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Huawei VTR-AL00",
		SystemVersion:       "SDK30",
		ApplicationVersion:  "10.15.4 (4945)",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})


	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	// Main loop
	rawUpdateMessages := client.GetRawUpdatesChannel(100)
	for updateMessage := range rawUpdateMessages {
		fmt.Println("===================================")
		for k, v := range updateMessage.Data {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}