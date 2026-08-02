package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addTimeScheduleCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.TimeScheduleList, client.TimeScheduleSummary]{
			Use:       "time-schedule",
			Short:     "Read Linear time schedules",
			ListShort: "List visible Linear time schedules",
			LimitHelp: "maximum time schedules to return",
			GetUse:    "get TIME_SCHEDULE_ID",
			GetShort:  "Get one time schedule by id",
			LoadList:  clientList(client.ListTimeSchedules),
			LoadGet:   clientGet(client.GetTimeScheduleByID),
			WriteItem: writeTimeSchedule,
		},
	)
}

func writeTimeSchedule(command *cobra.Command, options *rootOptions, schedule client.TimeScheduleSummary) error {
	return writeItemLine(
		command, options, schedule, schedule.ID,
		"%s %s entries %d", schedule.ID, schedule.Name, schedule.EntryCount,
	)
}
