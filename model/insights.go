// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"net/http"
	"time"
)

const (
	TimeRangeToday string = "today"
	TimeRange7Day  string = "7_day"
	TimeRange28Day string = "28_day"
)

type InsightsOpts struct {
	StartUnixMilli int64
	Page           int
	PerPage        int
}

type InsightsListData struct {
	HasNext bool `json:"has_next"`
}

// Top Reactions
type TopReactionList struct {
	InsightsListData
	Items []*TopReaction `json:"items"`
}

type TopReaction struct {
	EmojiName string `json:"emoji_name"`
	Count     int64  `json:"count"`
}

// Top Channels
type TopChannelList struct {
	InsightsListData
	Items []*TopChannel `json:"items"`
}

type TopChannel struct {
	ID           string      `json:"id"`
	Type         ChannelType `json:"type"`
	DisplayName  string      `json:"display_name"`
	Name         string      `json:"name"`
	TeamID       string      `json:"team_id"`
	MessageCount int64       `json:"message_count"`
}

type TopChannelByTimeList struct {
	InsightsListData
	Items []*TopChannelByTime `json:"items"`
}

type TopChannelByTime struct {
	Date  string        `json:"time"`
	Items []*TopChannel `json:"items"`
}

// GetStartUnixMilliForTimeRange gets the unix start time in milliseconds from the given time range.
// Time range can be one of: "1_day", "7_day", or "28_day".
func GetStartUnixMilliForTimeRange(timeRange string) (int64, *AppError) {
	now := time.Now()
	_, offset := now.Zone()
	switch timeRange {
	case TimeRangeToday:
		return GetStartOfDayMillis(now, offset), nil
	case TimeRange7Day:
		return GetStartOfDayMillis(now.Add(time.Hour*time.Duration(-168)), offset), nil
	case TimeRange28Day:
		return GetStartOfDayMillis(now.Add(time.Hour*time.Duration(-672)), offset), nil
	}

	return GetStartOfDayMillis(now, offset), NewAppError("Insights.IsValidRequest", "model.insights.time_range.app_error", nil, "", http.StatusBadRequest)
}

// GetTopReactionListWithPagination adds a rank to each item in the given list of TopReaction and checks if there is
// another page that can be fetched based on the given limit and offset. The given list of TopReaction is assumed to be
// sorted by Count. Returns a TopReactionList.
func GetTopReactionListWithPagination(reactions []*TopReaction, limit int) *TopReactionList {
	// Add pagination support
	var hasNext bool
	if (limit != 0) && (len(reactions) == limit+1) {
		hasNext = true
		reactions = reactions[:len(reactions)-1]
	}

	return &TopReactionList{InsightsListData: InsightsListData{HasNext: hasNext}, Items: reactions}
}

// GetTopChannelListWithPagination adds a rank to each item in the given list of TopChannel and checks if there is
// another page that can be fetched based on the given limit and offset. The given list of TopChannel is assumed to be
// sorted by Score. Returns a TopChannelList.
func GetTopChannelListWithPagination(channels []*TopChannel, limit int) *TopChannelList {
	// Add pagination support
	var hasNext bool
	if (limit != 0) && (len(channels) == limit+1) {
		hasNext = true
		channels = channels[:len(channels)-1]
	}

	return &TopChannelList{InsightsListData: InsightsListData{HasNext: hasNext}, Items: channels}
}
