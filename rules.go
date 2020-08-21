package skeleton

import (
	"regexp"
	"sync"
	"time"
)

// constants to limit border use methods bot
const (
	// only exec bot methods
	Private = 1
	// only exec group methods
	Group = 2
	// only exec channel methods
	Channel = 3
)

// constants to limit callbacks
const (
	// text messages
	Messages = "MessageRule"
	// callback to inline keyboard
	Callbacks = "CallbackRule"
	// edit messages callbacks
	EditedMessages = "EditedMessageRule"
	// commands, like a /start
	Commands = "CommandRule"
	// reply to message
	ReplyToMessages = "ReplyToMessageRule"
	// channel post (only in channels)
	ChannelPosts = "ChannelPostRule"
	// edit channel post callback (only in channels)
	EditedChannelPosts = "EditedChannelPostRule"
	// inline results callback
	InlineQuerys = "InlineQueryRule"
	// chosen inline results callbacks
	ChosenInlineResults = "ChosenInlineResultRule"
	// pre chekcout query callbacks
	PreCheckoutQuerys = "PreCheckoutQueryRule"
	// shipping query callbacks
	ShippingQuerys = "ShippingQueryRule"
)

// rules
type rules struct {
	sync.Mutex
	rulesMap map[string][]*Rule
}

// newRules()
func newRules() *rules {
	return &rules{
		rulesMap: make(map[string][]*Rule),
	}
}

// Rule
type Rule struct {
	// application storage
	app *app
	// previous pipeline rule
	prev *Rule
	// border user (private, group, channel)
	borderUse int
	// execution command
	command func(ctx *Context) bool
	// compile regex string
	regexp *regexp.Regexp
	// next pipeline rule
	next *Rule
	// time to end pipeline
	timer *time.Timer
	// timeout pipeline
	timeout *time.Duration
	// allowList for current rule
	allowList *AllowList
}

// HandleFunc
func (a *app) HandleFunc(regex string, handler func(c *Context) bool) *Rule {
	return &Rule{
		app:       a,
		allowList: newAllowList(),
		command:   handler,
		regexp:    regexp.MustCompile(regex),
		borderUse: Private,
	}
}

// Func
func (r *Rule) Func(handler func(c *Context) bool) *Rule {
	r.command = handler
	return r
}

// Border
func (r *Rule) Border(use int) *Rule {
	r.borderUse = use
	return r
}

// Methods
func (r *Rule) Methods(methods ...string) *Rule {
	rs := *r.app.rules
	for _, method := range methods {
		rs.Lock()
		rs.rulesMap[method] = append(rs.rulesMap[method], r)
		rs.Unlock()
	}

	return r
}

// Append
func (r *Rule) Append() *Rule {

	r.next = &Rule{
		app:  r.app,
		prev: r,
	}

	return r.next
}

// Timeout
func (r *Rule) Timeout(t time.Duration) *Rule {
	r.timeout = &t
	return r
}

// findMatch
func (r *Rule) findMatch(command string) []string {
	return r.regexp.FindStringSubmatch(command)
}

// allowList for current rule
func (r *Rule) AllowList() *AllowList {
	return r.allowList
}
