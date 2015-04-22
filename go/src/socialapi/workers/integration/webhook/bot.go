package webhook

import (
	"errors"
	"socialapi/models"
	"strconv"

	"github.com/koding/bongo"
)

const botNick = "bot"

var (
	ErrAccountIsNotParticipant = errors.New("account is not participant of the channel")
	ErrAccountNotFound         = errors.New("account not found")
	ErrGroupNotFound           = errors.New("group not found")
)

type Bot struct {
	account *models.Account
}

type Message struct {
	Body                 string // TODO check for XSS
	ChannelId            int64
	ChannelIntegrationId int64
}

func NewBot() (*Bot, error) {
	acc := models.NewAccount()
	if err := acc.ByNick(botNick); err != nil {
		return nil, err
	}

	return &Bot{account: acc}, nil
}

func (b *Bot) SendMessage(m *Message) error {
	cm, err := b.createMessage(m)
	if err != nil {
		return err
	}

	return b.createMessageList(cm, m.ChannelId)
}

func (b *Bot) createMessage(m *Message) (*models.ChannelMessage, error) {
	cm := models.NewChannelMessage()
	cm.AccountId = b.account.Id
	cm.InitialChannelId = m.ChannelId
	cm.Body = m.Body
	cm.TypeConstant = models.ChannelMessage_TYPE_BOT
	tid := strconv.FormatInt(m.ChannelIntegrationId, 10)
	cm.SetPayload("channelIntegrationId", tid)

	return cm, cm.Create()
}

func (b *Bot) createMessageList(cm *models.ChannelMessage, channelId int64) error {
	cml := models.NewChannelMessageList()
	cml.ChannelId = channelId
	cml.MessageId = cm.Id

	return cml.Create()
}

func (b *Bot) FetchBotChannel(username, groupName string) (*models.Channel, error) {

	// fetch account id
	acc := models.NewAccount()
	err := acc.ByNick(username)
	if err == bongo.RecordNotFound {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	// prevent sending bot messages when the user is not participant
	// of the given group
	canOpen, err := b.checkParticipation(acc, groupName)
	if err != nil {
		return nil, err
	}

	if !canOpen {
		return nil, ErrAccountIsNotParticipant
	}

	c, err := b.fetchOrCreateChannel(acc, groupName)
	if err != nil {
		return nil, err
	}

	// add user as participant
	_, err = c.AddParticipant(acc.Id)

	return c, err
}

func (b *Bot) fetchOrCreateChannel(a *models.Account, groupName string) (*models.Channel, error) {

	// fetch or create channel
	c, err := b.fetchBotChannel(a, groupName)
	if err == bongo.RecordNotFound {
		return b.createBotChannel(a, groupName)
	}

	if err != nil {
		return nil, err
	}

	return c, err
}

func (b *Bot) fetchBotChannel(a *models.Account, groupName string) (*models.Channel, error) {

	c := models.NewChannel()
	selector := map[string]interface{}{
		"creator_id":    a.Id,
		"type_constant": models.Channel_TYPE_BOT,
		"group_name":    groupName,
	}

	// if err is nil
	// it means we already have that channel
	err := c.One(bongo.NewQS(selector))

	return c, err
}

func (b *Bot) createBotChannel(a *models.Account, groupName string) (*models.Channel, error) {
	c := models.NewChannel()

	c.CreatorId = a.Id
	c.GroupName = groupName
	c.Name = models.RandomName()
	c.TypeConstant = models.Channel_TYPE_BOT

	err := c.Create()

	return c, err
}

func (b *Bot) checkParticipation(a *models.Account, groupName string) (bool, error) {
	c := models.NewChannel()

	selector := map[string]interface{}{
		"type_constant": models.Channel_TYPE_GROUP,
		"group_name":    groupName,
	}

	err := c.One(bongo.NewQS(selector))
	if err == bongo.RecordNotFound {
		return false, ErrGroupNotFound
	}

	if err != nil {
		return false, err
	}

	return c.CanOpen(a.Id)
}
