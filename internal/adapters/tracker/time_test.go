package tracker_test

import (
	"fmt"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/adapters/tracker"
	"testing"
	"time"
)

// TestTimeTracker is a mock implementation of the TimeTracker interface.
func TestTimeTracker(t *testing.T) {
	tracker := tracker.NewDefaultTimeTracker()

	tracker.Start()
	<-time.After(1 * time.Second)
	tracker.Stop()

	if tracker.GetDuration() < 1*time.Second {
		t.Errorf("Expected duration to be at least 1 second, got %v", tracker.GetDuration())
	}
	fmt.Printf("Duration: %v\n", tracker.GetDuration())
}
