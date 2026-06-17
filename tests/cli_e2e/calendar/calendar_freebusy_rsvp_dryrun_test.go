// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package calendar

import (
	"context"
	"testing"
	"time"

	clie2e "github.com/larksuite/cli/tests/cli_e2e"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCalendar_FreebusyDryRunIncludesUserAndRSVP(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())
	t.Setenv("LARKSUITE_CLI_APP_ID", "app")
	t.Setenv("LARKSUITE_CLI_APP_SECRET", "secret")
	t.Setenv("LARKSUITE_CLI_BRAND", "feishu")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	result, err := clie2e.RunCmd(ctx, clie2e.Request{
		Args: []string{
			"calendar", "+freebusy",
			"--start", "2026-04-25T09:00:00+08:00",
			"--end", "2026-04-25T18:00:00+08:00",
			"--user-id", "ou_dryrun_user",
			"--dry-run",
		},
		DefaultAs: "bot",
	})
	require.NoError(t, err)
	result.AssertExitCode(t, 0)

	out := result.Stdout
	require.Equal(t, "POST", gjson.Get(out, "api.0.method").String(), "stdout:\n%s", out)
	require.Equal(t, "/open-apis/calendar/v4/freebusy/list", gjson.Get(out, "api.0.url").String(), "stdout:\n%s", out)
	require.Equal(t, "ou_dryrun_user", gjson.Get(out, "api.0.body.user_id").String(), "stdout:\n%s", out)
	require.True(t, gjson.Get(out, "api.0.body.need_rsvp_status").Bool(), "stdout:\n%s", out)
	require.NotEmpty(t, gjson.Get(out, "api.0.body.time_min").String(), "stdout:\n%s", out)
	require.NotEmpty(t, gjson.Get(out, "api.0.body.time_max").String(), "stdout:\n%s", out)
}

func TestCalendar_RSVPDryRunUsesCalendarAndEventPath(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())
	t.Setenv("LARKSUITE_CLI_APP_ID", "app")
	t.Setenv("LARKSUITE_CLI_APP_SECRET", "secret")
	t.Setenv("LARKSUITE_CLI_BRAND", "feishu")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	result, err := clie2e.RunCmd(ctx, clie2e.Request{
		Args: []string{
			"calendar", "+rsvp",
			"--calendar-id", "cal_dry",
			"--event-id", "evt_dry",
			"--rsvp-status", "tentative",
			"--dry-run",
		},
		DefaultAs: "bot",
	})
	require.NoError(t, err)
	result.AssertExitCode(t, 0)

	out := result.Stdout
	require.Equal(t, "POST", gjson.Get(out, "api.0.method").String(), "stdout:\n%s", out)
	require.Equal(t, "/open-apis/calendar/v4/calendars/cal_dry/events/evt_dry/reply", gjson.Get(out, "api.0.url").String(), "stdout:\n%s", out)
	require.Equal(t, "tentative", gjson.Get(out, "api.0.body.rsvp_status").String(), "stdout:\n%s", out)
}
