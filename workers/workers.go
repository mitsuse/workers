package workers

import (
	"fmt"
	"os"
)

type Worker interface {
	Name() string
	Work()
}

func Log(worker Worker, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", worker.Name(), err)
}
