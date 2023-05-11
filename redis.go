package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
)

// Subscriber : Manage for Mastodon Traditional Style Relay Subscriber
type Subscriber struct {
	Domain     string `json:"domain,omitempty"`
	InboxURL   string `json:"inbox_url,omitempty"`
	ActivityID string `json:"activity_id,omitempty"`
	ActorID    string `json:"actor_id,omitempty"`
}

// Follower : Manage for LitePub Style Relay Follower
type Follower struct {
	Domain         string `json:"domain,omitempty"`
	InboxURL       string `json:"inbox_url,omitempty"`
	ActivityID     string `json:"activity_id,omitempty"`
	ActorID        string `json:"actor_id,omitempty"`
	MutuallyFollow bool   `json:"mutually_follow,omitempty"`
}

func GetSubcribesAndFollowers(rdb *redis.Client, ctx context.Context) []string {
	var subscribersAndFollowers []string
	subscribers := getSubscribes(rdb, ctx)
	followers := getFollows(rdb, ctx)
	for _, subscriber := range subscribers {
		subscribersAndFollowers = append(subscribersAndFollowers, subscriber.Domain)
	}
	for _, follower := range followers {
		subscribersAndFollowers = append(subscribersAndFollowers, follower.Domain)
	}
	return subscribersAndFollowers
}

func getFollows(rdb *redis.Client, ctx context.Context) []Follower {
	// 获取关注者列表
	var followers []Follower
	domains, _ := rdb.Keys(ctx, "relay:follower:*").Result()
	for _, domain := range domains {
		domainName := strings.Replace(domain, "relay:follower:", "", 1)
		inboxURL, _ := rdb.HGet(ctx, domain, "inbox_url").Result()
		activityID, err := rdb.HGet(ctx, domain, "activity_id").Result()
		if err != nil {
			activityID = ""
		}
		actorID, err := rdb.HGet(ctx, domain, "actor_id").Result()
		if err != nil {
			actorID = ""
		}
		mutuallyFollow, err := rdb.HGet(ctx, domain, "mutually_follow").Result()
		if err != nil {
			mutuallyFollow = "0"
		}
		followers = append(followers, Follower{domainName, inboxURL, activityID, actorID, mutuallyFollow == "1"})
		// subscribersAndFollowers = append(subscribersAndFollowers, Subscriber{domainName, inboxURL, activityID, actorID})
	}
	return followers
}

func getSubscribes(rdb *redis.Client, ctx context.Context) []Subscriber {
	// 获取订阅者列表
	var subscribers []Subscriber
	domains, _ := rdb.Keys(ctx, "relay:subscription:*").Result()
	for _, domain := range domains {
		domainName := strings.Replace(domain, "relay:subscription:", "", 1)
		inboxURL, _ := rdb.HGet(ctx, domain, "inbox_url").Result()
		activityID, err := rdb.HGet(ctx, domain, "activity_id").Result()
		if err != nil {
			activityID = ""
		}
		actorID, err := rdb.HGet(ctx, domain, "actor_id").Result()
		if err != nil {
			actorID = ""
		}
		subscribers = append(subscribers, Subscriber{domainName, inboxURL, activityID, actorID})
		// subscribersAndFollowers = append(subscribersAndFollowers, Subscriber{domainName, inboxURL, activityID, actorID})
	}
	return subscribers
}
