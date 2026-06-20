//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTimeScheduleCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.TimeScheduleList, client.TimeScheduleSummary]{
		Use:           "time-schedule",
		Short:         "Read Linear time schedules",
		ListShort:     "List visible Linear time schedules",
		LimitHelp:     "maximum time schedules to return",
		GetUse:        "get TIME_SCHEDULE_ID",
		GetShort:      "Get one time schedule by id",
		LoadList:      loadTimeScheduleList,
		PageWithItems: timeSchedulePageWithItems,
		LoadGet:       loadTimeSchedule,
		WriteItem:     writeTimeSchedule,
	})
}

func writeTimeSchedule(
	command *cobra.Command,
	options *rootOptions,
	schedule client.TimeScheduleSummary,
) error {
	if wrote, err := writeIDOnly(command, options, schedule.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, schedule)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s entries %d", schedule.ID, schedule.Name, schedule.EntryCount)
}

func loadTimeScheduleList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TimeScheduleList, []client.TimeScheduleSummary, error) {
	schedules, err := client.ListTimeSchedules(ctx, runtime.graphqlClient, limit)
	return schedules, schedules.TimeSchedules, err
}

func loadTimeSchedule(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.TimeScheduleSummary, error) {
	return client.GetTimeScheduleByID(ctx, runtime.graphqlClient, id)
}

func timeSchedulePageWithItems(
	page client.TimeScheduleList,
	schedules []client.TimeScheduleSummary,
) client.TimeScheduleList {
	page.TimeSchedules = schedules
	return page
}
