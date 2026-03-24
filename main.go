package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

// Forward declarations
void onQuit(id sender);
void onEditDate(id sender);

// Global state
static NSStatusItem *statusItem = nil;
static NSMenuItem *countdownMenuItem = nil;
static NSMenuItem *targetMenuItem = nil;

void initApp() {
    [NSApplication sharedApplication];
    [[NSApplication sharedApplication] setActivationPolicy:NSApplicationActivationPolicyAccessory];
}

void createStatusItem(const char *initialTitle) {
    NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
    statusItem = [statusBar statusItemWithLength:NSVariableStatusItemLength];

    NSFont *monoFont = [NSFont monospacedSystemFontOfSize:12.0 weight:NSFontWeightRegular];
    NSAttributedString *attrTitle = [[NSAttributedString alloc]
        initWithString:[NSString stringWithUTF8String:initialTitle]
            attributes:@{NSFontAttributeName: monoFont}];
    statusItem.button.attributedTitle = attrTitle;

    // Build menu
    NSMenu *menu = [[NSMenu alloc] init];

    // Countdown info item (disabled, just for display)
    countdownMenuItem = [[NSMenuItem alloc]
        initWithTitle:@"Calculating..."
               action:nil
        keyEquivalent:@""];
    [countdownMenuItem setEnabled:NO];
    [menu addItem:countdownMenuItem];

    // Target date item
    targetMenuItem = [[NSMenuItem alloc]
        initWithTitle:@""
               action:nil
        keyEquivalent:@""];
    [targetMenuItem setEnabled:NO];
    [menu addItem:targetMenuItem];

    [menu addItem:[NSMenuItem separatorItem]];

    // Quit
    NSMenuItem *quitItem = [[NSMenuItem alloc]
        initWithTitle:@"Quit Countdown"
               action:@selector(terminate:)
        keyEquivalent:@"q"];
    [quitItem setTarget:[NSApplication sharedApplication]];
    [menu addItem:quitItem];

    statusItem.menu = menu;
}

void updateMenuBar(const char *barTitle, const char *menuInfo, const char *targetLabel) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSFont *monoFont = [NSFont monospacedSystemFontOfSize:12.0 weight:NSFontWeightRegular];
        NSAttributedString *attrTitle = [[NSAttributedString alloc]
            initWithString:[NSString stringWithUTF8String:barTitle]
                attributes:@{NSFontAttributeName: monoFont}];
        statusItem.button.attributedTitle = attrTitle;

        if (countdownMenuItem) {
            [countdownMenuItem setTitle:[NSString stringWithUTF8String:menuInfo]];
        }
        if (targetMenuItem) {
            [targetMenuItem setTitle:[NSString stringWithUTF8String:targetLabel]];
        }
    });
}

void runApp() {
    [[NSApplication sharedApplication] run];
}
*/
import "C"
import (
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"
)

func freeStr(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {
	// ── Parse target date ──────────────────────────────────────────────────
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: countdown <YYYY-MM-DD> [\"Label\"]")
		fmt.Fprintln(os.Stderr, "  e.g. countdown 2025-12-31 \"New Year 2026\"")
		os.Exit(1)
	}

	targetDate, err := time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid date %q — use YYYY-MM-DD\n", os.Args[1])
		os.Exit(1)
	}
	// End of target day
	target := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		23, 59, 59, 0, time.Local,
	)

	label := "Target"
	if len(os.Args) >= 3 {
		label = os.Args[2]
	}

	// ── Init Cocoa ─────────────────────────────────────────────────────────
	C.initApp()

	init := C.CString("⏳ --d --h --m --s")
	C.createStatusItem(init)
	freeStr(init)

	targetStr := fmt.Sprintf("🎯  %s  ·  %s", label, targetDate.Format("Jan 2, 2006"))

	// ── Tick goroutine ─────────────────────────────────────────────────────
	go func() {
		for {
			now := time.Now()
			diff := target.Sub(now)

			var barTitle, menuInfo string

			if diff <= 0 {
				barTitle = "🎉 " + label + " — Done!"
				menuInfo = "🏁  Target date has passed"
			} else {
				totalSecs := int(math.Ceil(diff.Seconds()))
				days := totalSecs / 86400
				hours := (totalSecs % 86400) / 3600
				mins := (totalSecs % 3600) / 60
				secs := totalSecs % 60

				barTitle = fmt.Sprintf("⏳ %dd %02dh %02dm %02ds", days, hours, mins, secs)
				menuInfo = fmt.Sprintf("⏱  %d days  %02d hrs  %02d min  %02d sec", days, hours, mins, secs)
			}

			bt := C.CString(barTitle)
			mi := C.CString(menuInfo)
			tl := C.CString(targetStr)
			C.updateMenuBar(bt, mi, tl)
			freeStr(bt)
			freeStr(mi)
			freeStr(tl)

			time.Sleep(time.Second)
		}
	}()

	// ── Run the Cocoa event loop (blocks) ──────────────────────────────────
	C.runApp()
}
