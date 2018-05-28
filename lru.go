package kv

// Item 存储的数据项
type Item struct {
	value      string
	expireTime int64
}

// DNode 双向链表节点
type DNode struct {
	prev *DNode
	next *DNode
	key  string
	data *Item
}

// LruList LinkedHashMap
type LruList struct {
	head  *DNode
	tail  *DNode
	items map[string]*Item
}

// NewLruList 初始化LinkedHashMap
func NewLruList() *LruList {
	items := make(map[string]*Item)
	return &LruList{
		items: items,
	}
}

// Head 获取头节点
func (ll *LruList) Head() *DNode {
	return ll.head
}

// Tail 获取尾节点
func (ll *LruList) Tail() *DNode {
	return ll.tail
}

// IsHead 判断是不是头节点
func (ll *LruList) IsHead(node *DNode) bool {
	return ll.head == node
}

// IsTail 判断是不是尾节点
func (ll *LruList) IsTail(node *DNode) bool {
	return ll.tail == node
}

// Size 获取大小
func (ll *LruList) Size() int64 {
	return int64(len(ll.items))
}

// addHeadNode 添加到头节点
func (ll *LruList) addHeadNode(node *DNode) {

	if ll.Size() != 0 {
		ll.head.prev = node
		ll.tail.next = node
		node.prev = ll.tail
		node.next = ll.head

		ll.head = node
	} else {
		ll.head = node
		ll.tail = node
	}
}

// Set 添加到hashmap和链表表头
func (ll *LruList) Set(key string, item *Item) {

	exists := ll.Exists(key)

	var newNode *DNode

	if exists {
		node := ll.Del(key)

		if node == nil {
			return
		}

		newNode = node
	} else {
		newNode = &DNode{
			key:  key,
			data: item,
		}
	}

	ll.addHeadNode(newNode)

	ll.items[key] = item
}

// Del 删除数据节点
func (ll *LruList) Del(key string) *DNode {

	exists := ll.Exists(key)

	if !exists {
		return nil
	}

	node := ll.head
	for !ll.IsTail(node) || (node != nil && ll.head == ll.tail) {
		if node.key == key {

			if ll.Size() > 1 {
				if ll.IsHead(node) {
					ll.head = node.next
				}

				if ll.IsTail(node) {
					ll.tail = node.prev
				}

				node.prev.next = node.next
				node.next.prev = node.prev
			} else {
				ll.head = nil
				ll.tail = nil
			}

			delete(ll.items, key)

			return node
		}

		node = node.next
	}

	return nil
}

// Exists key
func (ll *LruList) Exists(key string) bool {
	_, exists := ll.items[key]
	return exists
}

// RemoveTailNode 删除尾部节点
func (ll *LruList) RemoveTailNode() *DNode {

	node := ll.Tail()

	if node != nil {

		if ll.Size() > 1 {
			ll.tail = ll.tail.prev

			ll.tail.next = ll.head
			ll.head.prev = ll.tail
		} else {
			ll.head = nil
			ll.tail = nil
		}

		delete(ll.items, node.key)

		return node
	}

	return nil
}
