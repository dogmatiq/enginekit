package nocopy

import "sync"

type NoCopy [0]sync.Mutex
