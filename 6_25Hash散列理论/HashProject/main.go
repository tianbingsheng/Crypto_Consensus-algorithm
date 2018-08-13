package main

import "HashProject/MyHashCode"

func main() {

	MyHashCode.CreateBuckets()
	MyHashCode.AddKeyValue("a1","hello China")
	MyHashCode.AddKeyValue("a1","hello kitty")
	MyHashCode.GetValueByKey("a1")
}
