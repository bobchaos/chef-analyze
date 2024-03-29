//
// Copyright 2019 Chef Software, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// use this function to execute a real chef-analyze command, under the hood the function
// will detect the correct binary depending on the OS and Architecture this test is running
// from, if you need to override the binary to use, set the `CHEF_ANALYZE_BIN` environment
// variable to the path of the binary you wish to use.
//
// example:
//
//  func TestHelpCommand(t *testing.T) {
//    out, err, exitcode := ChefAnalyze("help")
//
//    assert.Contains(t,
//      out.String(),
//      "Use \"chef-analyze [command] --help\" for more information about a command.",
//      "STDOUT doesn't match")
//    assert.Empty(t,
//      err.String(),
//      "STDERR should be empty")
//    assert.Equal(t, 0, exitcode,
//      "EXITCODE is not the expected one")
//  }
//
func ChefAnalyze(args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitcode int) {
	cmd := exec.Command(findChefAnalyzeBinary(), args...)
	cmd.Env = os.Environ()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	exitcode = 0
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitcode = exitError.ExitCode()
		} else {
			exitcode = 999
			if _, err := stderr.WriteString(err.Error()); err != nil {
				// we should never get here but if we do, lets print the error
				fmt.Println(err)
			}
		}
	}
	return
}

func findChefAnalyzeBinary() string {
	if bin := os.Getenv("CHEF_ANALYZE_BIN"); bin != "" {
		return bin
	}

	if runtime.GOOS != "" && runtime.GOARCH != "" {
		return fmt.Sprintf("chef-analyze_%s_%s", runtime.GOOS, runtime.GOARCH)
	}

	return "chef-analyze"
}
