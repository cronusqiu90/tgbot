package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

const (
	AppID    = 4
	AppHash  = "014b35b6184100b085b0d0572f9b5103"
	Password = ""
)

var (
	logger *zap.Logger
)

func handleNewMessage(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	msg, ok := update.Message.(*tg.Message)
	if !ok {
		logger.Info("OnNewMessage", zap.String("update", update.String()))
		return nil
	}
	msg.GetDate()
	msg.GetEditDate()
	msg.GetEditHide()
	msg.GetEffect()
	msg.GetEntities()
	msg.GetFactcheck()
	msg.GetForwards()
	msg.GetFromBoostsApplied()
	msg.GetFromID()
	msg.GetFromScheduled()
	msg.GetFwdFrom()
	msg.GetGroupedID()
	msg.GetID()
	msg.GetInvertMedia()
	msg.GetLegacy()
	if media, ok := msg.GetMedia(); ok {
		switch v := media.(type) {
		case *tg.MessageMediaEmpty: // messageMediaEmpty#3ded6320
		case *tg.MessageMediaPhoto: // messageMediaPhoto#695150d7
		case *tg.MessageMediaGeo: // messageMediaGeo#56e0d474
		case *tg.MessageMediaContact: // messageMediaContact#70322949
		case *tg.MessageMediaUnsupported: // messageMediaUnsupported#9f84f49e
		case *tg.MessageMediaDocument: // messageMediaDocument#dd570bd5
		case *tg.MessageMediaWebPage: // messageMediaWebPage#ddf10c3b
		case *tg.MessageMediaVenue: // messageMediaVenue#2ec0533f
		case *tg.MessageMediaGame: // messageMediaGame#fdb19008
		case *tg.MessageMediaInvoice: // messageMediaInvoice#f6a548d3
		case *tg.MessageMediaGeoLive: // messageMediaGeoLive#b940c666
		case *tg.MessageMediaPoll: // messageMediaPoll#4bd6e798
		case *tg.MessageMediaDice: // messageMediaDice#3f7ee58b
		case *tg.MessageMediaStory: // messageMediaStory#68cb6283
		case *tg.MessageMediaGiveaway: // messageMediaGiveaway#aa073beb
		case *tg.MessageMediaGiveawayResults: // messageMediaGiveawayResults#ceaa3ea1
		case *tg.MessageMediaPaidMedia: // messageMediaPaidMedia#a8852491
		default:
			panic(v)
		}
	}
	msg.GetMediaUnread()
	msg.GetMentioned()
	msg.GetMessage()
	msg.GetNoforwards()
	msg.GetOffline()
	msg.GetOut()
	msg.GetPeerID()
	msg.GetPinned()
	msg.GetPost()
	msg.GetPostAuthor()
	msg.GetQuickReplyShortcutID()
	msg.GetReactions()
	msg.GetReplies()
	msg.GetReplyMarkup()
	msg.GetReplyTo()
	msg.GetRestrictionReason()
	msg.GetSavedPeerID()
	msg.GetSilent()
	msg.GetTTLPeriod()
	msg.GetViaBotID()
	msg.GetViaBusinessBotID()
	msg.GetViews()

	logger.Info("OnNewMessage", zap.String("message", msg.String()))
	return nil
}

func signIn(ctx context.Context, phoneNumber string, client *telegram.Client) error {
	authCli := client.Auth()
	status, err := authCli.Status(ctx)
	if err != nil {
		return err
	}
	logger.Info("AuthStatus", zap.Bool("authorized", status.Authorized))

	if !status.Authorized {
		code, err := authCli.SendCode(ctx, phoneNumber, auth.SendCodeOptions{
			AllowAppHash: true,
		})
		if err != nil {
			return err
		}
		sentCode := code.(*tg.AuthSentCode)

		var verifyCode string
		fmt.Print("Enter verify code you received:")
		_, err = fmt.Scanln(&verifyCode)
		if err != nil {
			return err
		}

		_, err = authCli.SignIn(ctx, phoneNumber, verifyCode, sentCode.PhoneCodeHash)
		if err != nil {
			if err != auth.ErrPasswordAuthNeeded {
				return err
			}
			if _, err := authCli.Password(ctx, Password); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "time",
			// NameKey:      "logger",
			CallerKey: "caller",
			// FunctionKey:  "func",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	phoneNumber := ""
	sessionPath := ""
	flag.StringVar(&sessionPath, "s", "", "")
	flag.StringVar(&phoneNumber, "n", "", "")
	flag.Parse()

	sessionStorage := &telegram.FileSessionStorage{
		Path: sessionPath,
	}

	msgHandler := tg.NewUpdateDispatcher()

	client := telegram.NewClient(
		AppID,
		AppHash,
		telegram.Options{
			NoUpdates:      false,
			SessionStorage: sessionStorage,
			// Huawei BMH-AN20
			// X Plus 0.26.11.1733
			Device: telegram.DeviceConfig{
				SystemLangCode: "en",
				LangPack:       "",
				AppVersion:     "0.26.11.1733",
				SystemVersion:  "10.8.2",
				DeviceModel:    "Huawei BMH-AN20",
				LangCode:       "en",
			},
			UpdateHandler: msgHandler,
			Logger:        logger,
		},
	)
	msgHandler.OnNewMessage(handleNewMessage)

	if err := client.Run(ctx, func(ctx context.Context) error {
		if err := signIn(ctx, phoneNumber, client); err != nil {
			log.Fatal(err)
		}

		me, err := client.Self(ctx)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Login",
			zap.String("first_name", me.FirstName),
			zap.String("last_name", me.LastName),
			zap.String("user_name", me.Username),
			zap.Int64("uid", me.ID),
		)

		<-ctx.Done()

		return ctx.Err()
	}); err != nil {
		log.Fatal(err)
	}

}
