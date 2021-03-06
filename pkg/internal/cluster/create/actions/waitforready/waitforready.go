/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package waitforready

import (
	"fmt"
	"time"

	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/nodes"
	"github.com/gosoon/kubernetes-operator/pkg/internal/cluster/create/actions"
)

// Action implements an action for waiting for the cluster to be ready
type Action struct {
	waitTime time.Duration
}

// NewAction returns a new action for waiting for the cluster to be ready
func NewAction(waitTime time.Duration) actions.Action {
	return &Action{
		waitTime: waitTime,
	}
}

// Execute runs the action
func (a *Action) Execute(ctx *actions.ActionContext) error {
	// skip entirely if the wait time is 0
	if a.waitTime == time.Duration(0) {
		return nil
	}
	ctx.Status.Start(
		fmt.Sprintf(
			"Waiting ≤ %s for control-plane = Ready",
			formatDuration(a.waitTime),
		),
	)

	// Wait for the nodes to reach Ready status.
	startTime := time.Now()
	isReady := nodes.WaitForReady(startTime.Add(a.waitTime))
	if !isReady {
		ctx.Status.End(false)
		fmt.Println(" • WARNING: Timed out waiting for Ready ⚠️")
		return nil
	}
	// mark success
	ctx.Status.End(true)
	fmt.Printf(" • Ready after %s 💚\n", formatDuration(time.Since(startTime)))
	return nil
}

func formatDuration(duration time.Duration) string {
	return duration.Round(time.Second).String()
}
