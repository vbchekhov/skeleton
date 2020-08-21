package skeleton

// AllowList users
type AllowList map[int64]struct{}

// newAllowList
func newAllowList() *AllowList {
	return &AllowList{}
}

// Empty check len list
func (a *AllowList) Empty() bool {
	return len(*a) == 0
}

// Exist check exist user in allow list
func (a *AllowList) Exist(chatId int64) bool {
	c := *a
	_, ok := c[chatId]
	return ok
}

// Load load array users
func (a *AllowList) Load(arr ...int64) {
	c := *a
	for i := range arr {
		c[arr[i]] = struct{}{}
	}
	*a = c
}

// Append user in allow list
func (a *AllowList) Append(chatId int64) {
	c := *a
	c[chatId] = struct{}{}
	*a = c
}

// Delete user in allow list
func (a *AllowList) Delete(chatId int64) {
	c := *a
	delete(c, chatId)
	*a = c
}
