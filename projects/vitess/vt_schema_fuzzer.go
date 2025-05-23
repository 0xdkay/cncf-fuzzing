// Copyright 2022 ADA Logics Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package fuzzing

import (
	"fmt"
	fuzz "github.com/AdaLogics/go-fuzz-headers"

	"vitess.io/vitess/go/vt/schema"
	"vitess.io/vitess/go/vt/sqlparser"
)

// FuzzOnlineDDLFromCommentedStatement implements a fuzzer
// that targets schema.OnlineDDLFromCommentedStatement
func FuzzOnlineDDLFromCommentedStatement(data []byte) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	parser := sqlparser.NewTestParser()
	stmt, err := parser.Parse(string(data))
	if err != nil {
		return 0
	}
	onlineDDL, err := schema.OnlineDDLFromCommentedStatement(stmt)
	if err != nil {
		return 0
	}
	_, _ = onlineDDL.GetAction(parser)
	_, _, _ = onlineDDL.GetActionStr(parser)
	_ = onlineDDL.GetGCUUID()
	return 1
}

// FuzzNewOnlineDDLs implements a fuzzer that
// targets schema.NewOnlineDDLs
func FuzzNewOnlineDDLs(data []byte) int {
	parser := sqlparser.NewTestParser()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	f := fuzz.NewConsumer(data)

	keyspace, err := f.GetString()
	if err != nil {
		return 0
	}

	ddlstmtString, err := f.GetString()
	if err != nil {
		return 0
	}
	ddlStmt, _, err := schema.ParseOnlineDDLStatement(ddlstmtString, parser)
	if err != nil {
		return 0
	}

	sql, err := f.GetString()
	if err != nil {
		return 0
	}

	ddlStrategySetting := &schema.DDLStrategySetting{}
	err = f.GenerateStruct(ddlStrategySetting)
	if err != nil {
		return 0
	}

	requestContext, err := f.GetString()
	if err != nil {
		return 0
	}
	providedUUID, err := f.GetString()
	if err != nil {
		return 0
	}
	onlineDDLs, err := schema.NewOnlineDDLs(keyspace, sql, ddlStmt, ddlStrategySetting, requestContext, providedUUID, parser)
	if err != nil {
		return 0
	}
	for _, onlineDDL := range onlineDDLs {
		_, _ = onlineDDL.GetAction(parser)
		_, _, _ = onlineDDL.GetActionStr(parser)
		_ = onlineDDL.GetGCUUID()
	}
	return 1
}
