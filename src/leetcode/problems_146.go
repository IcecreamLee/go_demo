package leetcode

import "fmt"

//146. LRU缓存机制
//运用你所掌握的数据结构，设计和实现一个  LRU (最近最少使用) 缓存机制。它应该支持以下操作： 获取数据 get 和 写入数据 put 。
//
//获取数据 get(key) - 如果密钥 (key) 存在于缓存中，则获取密钥的值（总是正数），否则返回 -1。
//写入数据 put(key, value) - 如果密钥已经存在，则变更其数据值；如果密钥不存在，则插入该组「密钥/数据值」。当缓存容量达到上限时，它应该在写入新数据之前删除最久未使用的数据值，从而为新的数据值留出空间。
//
//进阶:
//
//你是否可以在 O(1) 时间复杂度内完成这两种操作？
//
//示例:
//
//LRUCache cache = new LRUCache( 2 /* 缓存容量 */ );
//
//cache.put(1, 1);
//cache.put(2, 2);
//cache.get(1);       // 返回  1
//cache.put(3, 3);    // 该操作会使得密钥 2 作废
//cache.get(2);       // 返回 -1 (未找到)
//cache.put(4, 4);    // 该操作会使得密钥 1 作废
//cache.get(1);       // 返回 -1 (未找到)
//cache.get(3);       // 返回  3
//cache.get(4);       // 返回  4

// 定义双向链表结构体
type DLinkedNode struct {
	key, value int
	prev, next *DLinkedNode
}

// 定义LRUCache结构体
type LRUCache struct {
	size, capacity int
	head           *DLinkedNode
	cache          map[int]*DLinkedNode
}

func Constructor(capacity int) LRUCache {
	lruCache := LRUCache{
		size:     0,
		capacity: capacity,
		head:     &DLinkedNode{},
		//tail:     &DLinkedNode{},
		cache: map[int]*DLinkedNode{},
	}
	// 循环双向链表
	lruCache.head.prev = lruCache.head
	lruCache.head.next = lruCache.head
	//lruCache.tail.prev = lruCache.head
	return lruCache
}

func (c *LRUCache) Get(key int) int {
	if _, ok := c.cache[key]; ok {
		c.move2head(c.cache[key])
		fmt.Println(c.head.next, c.head.prev)
		return c.cache[key].value
	}
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	if _, ok := c.cache[key]; ok {
		c.move2head(c.cache[key])
		c.cache[key].value = value
		fmt.Println(c.head.next, c.head.prev)
		return
	}

	c.cache[key] = &DLinkedNode{
		key:   key,
		value: value,
	}
	c.add2head(c.cache[key])
	c.size++
	if c.size > c.capacity {
		tailNode := c.removeTail()
		delete(c.cache, tailNode.key)
		c.size--
	}
	fmt.Println(c.head.next, c.head.prev)
}

func (c *LRUCache) add2head(node *DLinkedNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

func (c *LRUCache) move2head(node *DLinkedNode) {
	c.removeNode(node)
	c.add2head(node)
}

func (c *LRUCache) removeTail() *DLinkedNode {
	var node = c.head.prev
	c.removeNode(node)
	return node
}

func (c *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}
