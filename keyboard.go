package skeleton

import (
	"errors"
	"sort"
	"strconv"
	"sync"

	"github.com/Syfaro/telegram-bot-api"
)

// Keyboarder
type Keyboarder interface {
	Generate() *KeyboardResult
	Read(m *tgbotapi.Message, k *Keyboards) error
	Save(m *tgbotapi.Message, k *Keyboards) error
}

type Keyboards struct {
	sync.Mutex
	state map[int]interface{}
}

// set()
func (k *Keyboards) set(id int, kb interface{}) {
	k.Lock()
	defer k.Unlock()

	k.state[id] = kb
}

// get()
func (k *Keyboards) get(id int) (interface{}, bool) {
	k.Lock()
	defer k.Unlock()

	kb, ok := k.state[id]
	return kb, ok
}

// del()
func (k *Keyboards) del(id int) {
	k.Lock()
	defer k.Unlock()

	delete(k.state, id)
}

// ReadState
func (k *Keyboards) ReadState(m *tgbotapi.Message, kb Keyboarder) {
	kb.Read(m, k)
}

// SaveState
func (k *Keyboards) SaveState(m *tgbotapi.Message, kb Keyboarder) {
	kb.Save(m, k)
}

// newKeyboards()
func newKeyboards() *Keyboards {
	return &Keyboards{state: map[int]interface{}{}}
}

// KeyboardResult
type KeyboardResult struct {
	replyKeyboard  *tgbotapi.ReplyKeyboardMarkup
	inlineKeyboard *tgbotapi.InlineKeyboardMarkup
}

// ReplyKeyboardMarkup
func (kr *KeyboardResult) ReplyKeyboardMarkup() *tgbotapi.ReplyKeyboardMarkup {
	return kr.replyKeyboard
}

// InlineKeyboardMarkup
func (kr *KeyboardResult) InlineKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	return kr.inlineKeyboard
}

// -- Reply keyboard markup

// ReplyKeyboard
type ReplyKeyboard struct {
	ColumnsCount int
	Buttons      ReplyButton
}

// NewReplyKeyboard
func NewReplyKeyboard(columnsCount int) *ReplyKeyboard {
	return &ReplyKeyboard{ColumnsCount: columnsCount}
}

// Generate markup keyboard
func (k *ReplyKeyboard) Generate() *KeyboardResult {

	if len(k.Buttons) == 0 {
		return nil
	}

	var keyboard [][]tgbotapi.KeyboardButton

	for i := 0; i < len(k.Buttons); i = i + k.ColumnsCount {

		d := i + k.ColumnsCount
		if d > len(k.Buttons) {
			d = len(k.Buttons)
		}

		keyboard = append(keyboard, k.Buttons[i:d])
	}

	return &KeyboardResult{
		replyKeyboard: &tgbotapi.ReplyKeyboardMarkup{
			Keyboard:        keyboard,
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		},
	}
}

// save() save keyboard in memory
func (k *ReplyKeyboard) Save(m *tgbotapi.Message) error {

	return nil
}

// ReplyButton
type ReplyButton []tgbotapi.KeyboardButton

// Add() add new markup button
func (b *ReplyButton) Add(text string) {
	*b = append(*b, tgbotapi.NewKeyboardButton(text))
}

// -- Reply keyboard markup

// -- Inline keyboard markup

// InlineKeyboard
type InlineKeyboard struct {
	// using MessageID
	Id int
	// Title menu
	Title string
	// using user ChatID
	ChatID int64
	// columns count in keyboard
	ColumnsCount int
	// columns count in one page keyboard
	StringCount int
	// map all Pages
	Pages map[int][][]tgbotapi.InlineKeyboardButton
	// all buttons
	Buttons InlineButton
	// control buttons
	ControlButtons InlineButton
	// current open page
	CurrentPage int
	// if you need dont save keyboard in cacheMemory
	// actual on show method
	DontSave bool
}

// NewInlineKeyboard
func NewInlineKeyboard(columnsCount, stringCount int) *InlineKeyboard {
	return &InlineKeyboard{
		ColumnsCount: columnsCount,
		StringCount:  stringCount,
	}
}

// NewAbortPipelineKeyboard
func NewAbortPipelineKeyboard(title string) *tgbotapi.InlineKeyboardMarkup {
	kb := NewInlineKeyboard(1, 1)
	kb.Buttons.Add(title, "abort-pipeline")

	return kb.Generate().InlineKeyboardMarkup()
}

// Generate inline keyboard
func (k *InlineKeyboard) Generate() *KeyboardResult {

	var keyboard [][]tgbotapi.InlineKeyboardButton

	// ++ Create simple keyboard

	// create columns
	for i := 0; i < len(k.Buttons); i = i + k.ColumnsCount {

		d := i + k.ColumnsCount
		if d > len(k.Buttons) {
			d = len(k.Buttons)
		}

		keyboard = append(keyboard, k.Buttons[i:d])

	}

	// back simple keyboard
	if k.StringCount == 0 {
		// if have control button app end
		if len(k.ControlButtons) > 0 {
			keyboard = append(keyboard, k.ControlButtons)
		}
		return &KeyboardResult{inlineKeyboard: &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
	}

	// -- Create simple keyboard

	// ++ Create page list keyboard

	// number first page
	num := 1
	// make map Pages
	k.Pages = make(map[int][][]tgbotapi.InlineKeyboardButton)

	// create Pages
	for i := 0; i < len(keyboard); i = i + k.StringCount {

		d := i + k.StringCount
		if d > len(keyboard) {
			d = len(keyboard)
		}

		k.Pages[num] = make([][]tgbotapi.InlineKeyboardButton, len(keyboard[i:d]))

		copy(k.Pages[num], keyboard[i:d])

		if len(k.ControlButtons) > 0 {
			if num == 1 {
				k.Pages[num] = append(k.Pages[num], k.ControlButtons[1:])
			} else {
				k.Pages[num] = append(k.Pages[num], k.ControlButtons)
			}
		}

		num++
	}

	// set current page
	k.CurrentPage = 1
	keyboard = k.Pages[1]

	// -- Create page list

	return &KeyboardResult{inlineKeyboard: &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}}

}

// Read
func (k *InlineKeyboard) Read(message *tgbotapi.Message, ks *Keyboards) error {

	kb, ok := ks.get(message.MessageID)
	if !ok {
		return errors.New("Keyboard not found in cache!\n")
	}

	*k = kb.(InlineKeyboard)

	return nil
}

// Save
func (k *InlineKeyboard) Save(message *tgbotapi.Message, ks *Keyboards) error {

	k.Id = message.MessageID
	k.ChatID = message.Chat.ID
	k.Title = message.Text

	ks.set(k.Id, *k)

	return nil
}

// ++ default control func`s

// Next() next page
func (k *InlineKeyboard) Next() tgbotapi.EditMessageReplyMarkupConfig {
	if k.CurrentPage <= len(k.Pages)-1 {
		k.CurrentPage++
	}

	return k.Show()
}

func next(c *Context) bool {

	k := NewInlineKeyboard(0, 0)

	c.Keyboards().ReadState(c.Update.CallbackQuery.Message, k)

	m, _ := c.BotAPI.Send(k.Next())
	c.Keyboards().SaveState(&m, k)

	return true
}

// Prev() previous page
func (k *InlineKeyboard) Prev() tgbotapi.EditMessageReplyMarkupConfig {
	k.CurrentPage--
	if k.CurrentPage <= 1 {
		k.CurrentPage = 1
	}

	return k.Show()
}

func prev(c *Context) bool {

	k := NewInlineKeyboard(0, 0)

	c.Keyboards().ReadState(c.Update.CallbackQuery.Message, k)

	m, _ := c.BotAPI.Send(k.Prev())
	c.Keyboards().SaveState(&m, k)

	return true
}

// List() all Pages list
func (k *InlineKeyboard) List() tgbotapi.EditMessageReplyMarkupConfig {
	list := NewInlineKeyboard(3, 0)
	list.Id = k.Id
	list.ChatID = k.ChatID

	var pages []int
	for v := range k.Pages {
		pages = append(pages, v)
	}

	sort.Ints(pages)

	for _, page := range pages {
		list.Buttons.Add(
			strconv.Itoa(page),
			"show-"+strconv.Itoa(page))
	}

	list.ControlButtons.Add("⬅️back", "back")

	return tgbotapi.NewEditMessageReplyMarkup(k.ChatID, k.Id, *list.Generate().InlineKeyboardMarkup())
}

func list(c *Context) bool {

	k := NewInlineKeyboard(0, 0)

	c.Keyboards().ReadState(c.Update.CallbackQuery.Message, k)

	m, _ := c.BotAPI.Send(k.List())
	c.Keyboards().SaveState(&m, k)

	return true
}

// Show() show current list
func (k *InlineKeyboard) Show() tgbotapi.EditMessageReplyMarkupConfig {

	return tgbotapi.NewEditMessageReplyMarkup(k.ChatID, k.Id, tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: k.Pages[k.CurrentPage],
	})
}

func show(c *Context) bool {

	k := NewInlineKeyboard(0, 0)

	c.Keyboards().ReadState(c.Update.CallbackQuery.Message, k)
	k.CurrentPage, _ = strconv.Atoi(c.RegexpResult[1])

	m, _ := c.BotAPI.Send(k.Show())
	c.Keyboards().SaveState(&m, k)

	return true
}

// Back()
func (k *InlineKeyboard) Back() tgbotapi.EditMessageTextConfig {

	msg := tgbotapi.NewEditMessageText(k.ChatID, k.Id, k.Title)
	msg.ParseMode = parseMode
	msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: k.Pages[k.CurrentPage],
	}

	return msg
}

func back(c *Context) bool {

	k := NewInlineKeyboard(0, 0)

	c.Keyboards().ReadState(c.Update.CallbackQuery.Message, k)

	c.BotAPI.Send(k.Back())

	return true
}

// -- default control func`s

// InlineButton
type InlineButton []tgbotapi.InlineKeyboardButton

// Add() add Data type button
func (b *InlineButton) Add(text, data string) {
	*b = append(*b, tgbotapi.NewInlineKeyboardButtonData(text, data))
}

// AddURL() add URI type button
func (b *InlineButton) AddURL(text, url string) {
	*b = append(*b, tgbotapi.NewInlineKeyboardButtonURL(text, url))
}

// AddSwitch() add Switch type button
func (b *InlineButton) AddSwitch(text, swh string) {
	*b = append(*b, tgbotapi.NewInlineKeyboardButtonSwitch(text, swh))
}

// AddPagesControl
func (b *InlineButton) AddPagesControl() {
	b.Add("⬅️", "prev")
	b.Add("...", "list")
	b.Add("➡️", "next")
}
