package sdp

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(file), "resource")
}

func TestSdpBootstrapBranches(t *testing.T) {
	branchesFile := path.Join(resourceDir(), "branches_file.txt")
	os.Setenv(BranchesFileEnvVar, branchesFile)
	branches, err := BootstrapBranches()
	require.NoError(t, err)

	require.Equal(t, 3, len(branches))
	assert.Equal(t, "foo", branches[0])
	assert.Equal(t, "bar", branches[1])
	assert.Equal(t, "baz", branches[2])
}
