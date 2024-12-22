package bot

import (
	"os"
	"strconv"
	"time"

	"github.com/NicoNex/echotron/v3"
)

var (
	token      = os.Getenv("BOT_TOKEN")
	adminID, _ = strconv.Atoi(os.Getenv("ADMIN_ID"))

	TgBot = echotron.NewAPI(token)
)

func SendTextFileToAdmin(filename, text string) {
	file := echotron.NewInputFileBytes(filename, []byte(text))

	TgBot.SendDocument(file, int64(adminID), &echotron.DocumentOptions{Caption: time.Now().String()})
}
