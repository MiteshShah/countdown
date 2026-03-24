package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static NSStatusItem  *statusItem    = nil;
static NSMenuItem    *vizMenuItem   = nil;
static NSMenuItem    *secsMenuItem  = nil;
static NSMenuItem    *countdownItem = nil;
static NSMenuItem    *targetItem    = nil;
static NSMenuItem    *pctItem       = nil;

void initApp() {
    [NSApplication sharedApplication];
    [[NSApplication sharedApplication] setActivationPolicy:NSApplicationActivationPolicyAccessory];
}

static NSColor* hourColor(int hoursLeft) {
    if (hoursLeft >= 12) return [NSColor colorWithRed:0.38 green:0.60 blue:0.14 alpha:1.0];
    if (hoursLeft >= 6)  return [NSColor colorWithRed:0.87 green:0.62 blue:0.15 alpha:1.0];
    if (hoursLeft >= 3)  return [NSColor colorWithRed:0.84 green:0.38 blue:0.19 alpha:1.0];
    return                      [NSColor colorWithRed:0.89 green:0.30 blue:0.29 alpha:1.0];
}

static NSAttributedString* buildHourViz(int currentH) {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];
    NSFont *font = [NSFont monospacedSystemFontOfSize:13.0 weight:NSFontWeightRegular];
    int hoursLeft = 23 - currentH;
    NSColor *fc = hourColor(hoursLeft);
    NSColor *dim = [NSColor colorWithWhite:0.40 alpha:1.0];

	// ✅ ADD THIS HEADER
    NSString *header = [NSString stringWithFormat:@"Hours left today: %d\n\n", hoursLeft];
    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:header
            attributes:@{
                NSFontAttributeName: [NSFont systemFontOfSize:12.0 weight:NSFontWeightSemibold],
                NSForegroundColorAttributeName: fc
            }]];

    for (int h = 0; h < 24; h++) {
        NSString *label = [NSString stringWithFormat:@"%02d", h];
        NSColor *color = (h < currentH) ? dim : fc;
        NSMutableDictionary *attrs = [@{
            NSFontAttributeName: font,
            NSForegroundColorAttributeName: color
        } mutableCopy];
        if (h == currentH) {
            attrs[NSUnderlineStyleAttributeName] = @(NSUnderlineStyleSingle);
            attrs[NSUnderlineColorAttributeName] = fc;
        }
        [s appendAttributedString:[[NSAttributedString alloc] initWithString:label attributes:attrs]];
        NSString *sep = (h % 8 == 7) ? @"\n" : @" ";
        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:sep
                attributes:@{NSFontAttributeName: font,
                              NSForegroundColorAttributeName: [NSColor clearColor]}]];
    }
    return s;
}

static NSAttributedString* buildSecViz(int currentH, int currentS) {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];
    NSFont *font = [NSFont monospacedSystemFontOfSize:10.0 weight:NSFontWeightRegular];
    NSColor *activeColor = hourColor(23 - currentH);
    NSColor *passedColor = [activeColor colorWithAlphaComponent:0.35];
    NSColor *dimC = [NSColor colorWithWhite:0.22 alpha:1.0];

	// ✅ HEADER (same style as hours)
	int secLeft = 59 - currentS;
    NSString *header = [NSString stringWithFormat:@"Seconds left: %d\n\n", secLeft];
    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:header
            attributes:@{
                NSFontAttributeName: [NSFont systemFontOfSize:11.0 weight:NSFontWeightSemibold],
                NSForegroundColorAttributeName: activeColor
            }]];

    for (int i = 0; i < 60; i++) {
        NSColor *c = (i < currentS) ? passedColor : (i == currentS) ? activeColor : dimC;
        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:@"■"
                attributes:@{NSFontAttributeName: font,
                              NSForegroundColorAttributeName: c}]];
        if (i == 29) {
            [s appendAttributedString:[[NSAttributedString alloc]
                initWithString:@"\n" attributes:@{NSFontAttributeName: font}]];
        }
    }
    return s;
}

void createStatusItem(const char *initialTitle) {
    NSStatusBar *bar = [NSStatusBar systemStatusBar];
    statusItem = [bar statusItemWithLength:NSVariableStatusItemLength];
    NSFont *mono = [NSFont monospacedSystemFontOfSize:12.0 weight:NSFontWeightRegular];
    statusItem.button.attributedTitle = [[NSAttributedString alloc]
        initWithString:[NSString stringWithUTF8String:initialTitle]
            attributes:@{NSFontAttributeName: mono}];

    NSMenu *menu = [[NSMenu alloc] init];
    [menu setAutoenablesItems:NO];

    vizMenuItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [vizMenuItem setEnabled:NO];
    [menu addItem:vizMenuItem];

    [menu addItem:[NSMenuItem separatorItem]];

    secsMenuItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [secsMenuItem setEnabled:NO];
    [menu addItem:secsMenuItem];

    [menu addItem:[NSMenuItem separatorItem]];

    pctItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [pctItem setEnabled:NO];
    [menu addItem:pctItem];

    [menu addItem:[NSMenuItem separatorItem]];

    countdownItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [countdownItem setEnabled:NO];
    [menu addItem:countdownItem];

    targetItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [targetItem setEnabled:NO];
    [menu addItem:targetItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *quit = [[NSMenuItem alloc]
        initWithTitle:@"Quit" action:@selector(terminate:) keyEquivalent:@"q"];
    [quit setTarget:[NSApplication sharedApplication]];
    [menu addItem:quit];

    statusItem.menu = menu;
}

void updateAll(const char *barTitle, int currentH, int currentS,
               const char *pctLabel, const char *countdownStr, const char *targetStr) {
    NSString *bt  = barTitle     ? [NSString stringWithUTF8String:barTitle]     : @"";
    NSString *pct = pctLabel     ? [NSString stringWithUTF8String:pctLabel]     : @"";
    NSString *cd  = countdownStr ? [NSString stringWithUTF8String:countdownStr] : @"";
    NSString *tl  = targetStr    ? [NSString stringWithUTF8String:targetStr]    : @"";
    int h = currentH, s = currentS;

    dispatch_async(dispatch_get_main_queue(), ^{
        NSFont *mono = [NSFont monospacedSystemFontOfSize:12.0 weight:NSFontWeightRegular];
        statusItem.button.attributedTitle = [[NSAttributedString alloc]
            initWithString:bt attributes:@{NSFontAttributeName: mono}];
        if (vizMenuItem)   [vizMenuItem   setAttributedTitle:buildHourViz(h)];
        if (secsMenuItem)  [secsMenuItem  setAttributedTitle:buildSecViz(h, s)];
        if (pctItem)       [pctItem       setTitle:pct];
        if (countdownItem) [countdownItem setTitle:cd];
        if (targetItem)    [targetItem    setTitle:tl];
    });
}

void runApp() { [[NSApplication sharedApplication] run]; }
*/
import "C"
import (
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"
)

func cstr(s string) *C.char { return C.CString(s) }
func cfree(p *C.char)       { C.free(unsafe.Pointer(p)) }

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: countdown <YYYY-MM-DD> [\"Label\"]")
		os.Exit(1)
	}
	targetDate, err := time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid date %q — use YYYY-MM-DD\n", os.Args[1])
		os.Exit(1)
	}
	target := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		23, 59, 59, 0, time.Local,
	)
	label := "Target"
	if len(os.Args) >= 3 {
		label = os.Args[2]
	}

	C.initApp()
	init := cstr("⏳ loading…")
	C.createStatusItem(init)
	cfree(init)

	targetStr := fmt.Sprintf("🎯  %s  —  %s", label, targetDate.Format("Jan 2, 2006"))

	go func() {
		for {
			now := time.Now()
			H, M, S := now.Hour(), now.Minute(), now.Second()

			diff := target.Sub(now)
			var barTitle, cdStr string
			if diff <= 0 {
				barTitle = "🎉 " + label + " — Done!"
				cdStr = "🏁  Target date has passed"
			} else {
				total := int(math.Ceil(diff.Seconds()))
				d := total / 86400
				h := (total % 86400) / 3600
				m := (total % 3600) / 60
				s := total % 60
				barTitle = fmt.Sprintf("⏳ %dd %02dh %02dm %02ds", d, h, m, s)
				cdStr = fmt.Sprintf("⏱  %d days  %02dh %02dm %02ds remaining", d, h, m, s)
			}

			secsToday := H*3600 + M*60 + S
			pct := secsToday * 100 / 86400
			hoursLeft := 23 - H
			pctStr := fmt.Sprintf("📅  Day used: %d%%  ·  %d hours left today", pct, hoursLeft)

			bt := cstr(barTitle)
			cd := cstr(cdStr)
			tl := cstr(targetStr)
			pl := cstr(pctStr)
			C.updateAll(bt, C.int(H), C.int(S), pl, cd, tl)
			cfree(bt)
			cfree(cd)
			cfree(tl)
			cfree(pl)

			time.Sleep(time.Second)
		}
	}()

	C.runApp()
}
