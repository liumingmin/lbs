package main

import (
	"context"

	"github.com/liumingmin/goutils/container"
	"github.com/liumingmin/goutils/log"
	"github.com/xuri/excelize/v2"
)

func ReadExcel(ctx context.Context, f *excelize.File, sheetName string) *container.DataTable {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		log.Error(ctx, "GetRows excel err: %v", err)
		return nil
	}

	dataTable := container.NewDataTable(rows[0], rows[0][0], []string{}, 1000)
	dataTable.PushAll(rows[1:])
	return dataTable
}
