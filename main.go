package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static NSStatusItem *statusItem = nil;
static NSMenuItem   *mainItem   = nil;

void initApp() {
    [NSApplication sharedApplication];
    [[NSApplication sharedApplication] setActivationPolicy:NSApplicationActivationPolicyAccessory];
    [NSApp finishLaunching];
}

static NSAttributedString* sectionTitle(NSString *title) {
    return [[NSAttributedString alloc]
        initWithString:[title stringByAppendingString:@"\n\n"]
            attributes:@{
                NSFontAttributeName: [NSFont monospacedSystemFontOfSize:11 weight:NSFontWeightSemibold],
                NSForegroundColorAttributeName: [NSColor colorWithWhite:0.6 alpha:1.0]
            }];
}

static NSAttributedString* buildHeader() {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];

    NSDate *now = [NSDate date];

    NSDateFormatter *timeFmt = [[NSDateFormatter alloc] init];
    [timeFmt setDateFormat:@"HH:mm:ss"];

    NSDateFormatter *dateFmt = [[NSDateFormatter alloc] init];
    [dateFmt setDateFormat:@"EEEE, MMM d, yyyy"];

    NSString *timeStr = [timeFmt stringFromDate:now];
    NSString *dateStr = [dateFmt stringFromDate:now];

    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:[timeStr stringByAppendingString:@"\n"]
            attributes:@{
                NSFontAttributeName: [NSFont monospacedDigitSystemFontOfSize:16 weight:NSFontWeightMedium],
                NSForegroundColorAttributeName: [NSColor whiteColor]
            }]];

    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:[dateStr stringByAppendingString:@"\n\n"]
            attributes:@{
                NSFontAttributeName: [NSFont systemFontOfSize:12],
                NSForegroundColorAttributeName: [NSColor colorWithWhite:0.7 alpha:1.0]
            }]];

    return s;
}

static NSColor* hourColor(int hoursLeft) {
    if (hoursLeft >= 12) return [NSColor colorWithRed:0.38 green:0.60 blue:0.14 alpha:1.0];
    if (hoursLeft >= 6)  return [NSColor colorWithRed:0.87 green:0.62 blue:0.15 alpha:1.0];
    if (hoursLeft >= 3)  return [NSColor colorWithRed:0.84 green:0.38 blue:0.19 alpha:1.0];
    return                      [NSColor colorWithRed:0.89 green:0.30 blue:0.29 alpha:1.0];
}

static NSAttributedString* buildHourViz(int currentH) {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];
    [s appendAttributedString:sectionTitle(@"TODAY'S 24 HOURS")];

    NSFont *font = [NSFont monospacedSystemFontOfSize:13.0 weight:NSFontWeightRegular];

    int hoursLeft = 23 - currentH;
    NSColor *fc = hourColor(hoursLeft);
    NSColor *dim = [NSColor colorWithWhite:0.4 alpha:1.0];

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

        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:label attributes:attrs]];

        NSString *sep = (h % 8 == 7) ? @"\n" : @" ";
        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:sep attributes:@{NSFontAttributeName: font}]];
    }

    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:@"\n\n" attributes:@{}]];
    return s;
}

static NSAttributedString* buildMinuteViz(int currentM) {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];
    [s appendAttributedString:sectionTitle(@"MINUTES THIS HOUR")];

    NSFont *font = [NSFont monospacedSystemFontOfSize:10 weight:NSFontWeightRegular];

    for (int i = 0; i < 60; i++) {
        NSColor *c = (i < currentM)
            ? [NSColor colorWithRed:0.55 green:0.75 blue:0.35 alpha:1.0]
            : [NSColor colorWithWhite:0.2 alpha:1.0];

        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:@"▮"
                attributes:@{
                    NSFontAttributeName: font,
                    NSForegroundColorAttributeName: c
                }]];
    }

    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:@"\n\n" attributes:@{}]];
    return s;
}

static NSAttributedString* buildSecViz(int currentS) {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];
    [s appendAttributedString:sectionTitle(@"SECONDS THIS MINUTE")];

    NSFont *font = [NSFont monospacedSystemFontOfSize:10 weight:NSFontWeightRegular];

    for (int i = 0; i < 60; i++) {
        NSColor *c = (i < currentS)
            ? [NSColor colorWithRed:0.4 green:0.8 blue:0.7 alpha:1.0]
            : [NSColor colorWithWhite:0.2 alpha:1.0];

        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:@"■"
                attributes:@{
                    NSFontAttributeName: font,
                    NSForegroundColorAttributeName: c
                }]];
    }

    [s appendAttributedString:[[NSAttributedString alloc]
        initWithString:@"\n\n" attributes:@{}]];
    return s;
}

static NSAttributedString* buildLegend() {
    NSMutableAttributedString *s = [[NSMutableAttributedString alloc] init];

    NSArray *items = @[
        @{@"label": @"past", @"color": [NSColor lightGrayColor]},
        @{@"label": @"12+ h", @"color": [NSColor colorWithRed:0.6 green:0.8 blue:0.4 alpha:1]},
        @{@"label": @"6–12 h", @"color": [NSColor colorWithRed:0.9 green:0.7 blue:0.3 alpha:1]},
        @{@"label": @"3–6 h", @"color": [NSColor colorWithRed:0.9 green:0.5 blue:0.3 alpha:1]},
        @{@"label": @"0–3 h", @"color": [NSColor colorWithRed:0.9 green:0.3 blue:0.3 alpha:1]},
    ];

    for (NSDictionary *item in items) {
        NSString *text = [NSString stringWithFormat:@"■ %@  ", item[@"label"]];
        [s appendAttributedString:[[NSAttributedString alloc]
            initWithString:text
                attributes:@{
                    NSFontAttributeName: [NSFont systemFontOfSize:11],
                    NSForegroundColorAttributeName: item[@"color"]
                }]];
    }

    return s;
}

void createStatusItem(const char *initialTitle) {
    NSStatusBar *bar = [NSStatusBar systemStatusBar];
    statusItem = [bar statusItemWithLength:NSVariableStatusItemLength];

    NSFont *mono = [NSFont monospacedDigitSystemFontOfSize:12 weight:NSFontWeightRegular];
    statusItem.button.attributedTitle = [[NSAttributedString alloc]
        initWithString:[NSString stringWithUTF8String:initialTitle]
        attributes:@{NSFontAttributeName: mono}];

    NSMenu *menu = [[NSMenu alloc] init];

    mainItem = [[NSMenuItem alloc] initWithTitle:@"" action:nil keyEquivalent:@""];
    [mainItem setEnabled:NO];
    [menu addItem:mainItem];

    [menu addItem:[NSMenuItem separatorItem]];

    NSMenuItem *quit = [[NSMenuItem alloc]
        initWithTitle:@"Quit" action:@selector(terminate:) keyEquivalent:@"q"];
    [quit setTarget:[NSApplication sharedApplication]];
    [menu addItem:quit];

    statusItem.menu = menu;
}

void updateUI(const char *title, int h, int m, int s, const char *footer) {
    NSString *t = [NSString stringWithUTF8String:title];
    NSString *f = [NSString stringWithUTF8String:footer];

    dispatch_async(dispatch_get_main_queue(), ^{
        NSFont *mono = [NSFont monospacedDigitSystemFontOfSize:12 weight:NSFontWeightRegular];

        statusItem.button.attributedTitle = [[NSAttributedString alloc]
            initWithString:t
                attributes:@{NSFontAttributeName: mono}];

        NSMutableAttributedString *full = [[NSMutableAttributedString alloc] init];

        [full appendAttributedString:buildHeader()];
        [full appendAttributedString:buildHourViz(h)];
        [full appendAttributedString:buildMinuteViz(m)];
        [full appendAttributedString:buildSecViz(s)];
        [full appendAttributedString:buildLegend()];

        [full appendAttributedString:[[NSAttributedString alloc]
            initWithString:[@"\n\n" stringByAppendingString:f]
                attributes:@{
                    NSFontAttributeName: [NSFont systemFontOfSize:12],
                    NSForegroundColorAttributeName: [NSColor whiteColor]
                }]];

        [mainItem setAttributedTitle:full];
    });
}

void runApp() { [[NSApplication sharedApplication] run]; }
*/
import "C"

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
	"unsafe"
)

func cstr(s string) *C.char { return C.CString(s) }
func cfree(p *C.char)       { C.free(unsafe.Pointer(p)) }

func main() {
	runtime.LockOSThread()

	if len(os.Args) < 2 {
		fmt.Println("Usage: countdown <YYYY-MM-DD> [label]")
		return
	}

	label := "Target"
	if len(os.Args) >= 3 && os.Args[2] != "" {
		label = os.Args[2]
	}

	targetDate, err := time.ParseInLocation("2006-01-02", os.Args[1], time.Local)
	if err != nil {
		fmt.Println("Invalid date. Use YYYY-MM-DD")
		return
	}

	target := time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		23, 59, 59, 0, time.Local,
	)

	C.initApp()

	init := cstr("⏳")
	C.createStatusItem(init)
	cfree(init)

	go func() {
		for {
			now := time.Now()
			H, M, S := now.Hour(), now.Minute(), now.Second()

			diff := target.Sub(now)

			var d, h, m, s int
			if diff > 0 {
				total := int(math.Ceil(diff.Seconds()))
				d = total / 86400
				h = (total % 86400) / 3600
				m = (total % 3600) / 60
				s = total % 60
			}

			pct := (H*3600 + M*60 + S) * 100 / 86400

			var title, footer string

			if diff <= 0 {
				title = "🎉 " + label
				footer = fmt.Sprintf("🎉 %s — Done!    TODAY USED %d%%", label, pct)
			} else {
				title = fmt.Sprintf("⏳ %dd %02dh %02dm %02ds", d, h, m, s)
				footer = fmt.Sprintf(
					"⏳ %dd %02dh %02dm %02ds · %s · TODAY %d%%",
					d, h, m, s,
					label,
					pct,
				)
			}

			t := cstr(title)
			f := cstr(footer)

			C.updateUI(t, C.int(H), C.int(M), C.int(S), f)

			cfree(t)
			cfree(f)

			time.Sleep(time.Second)
		}
	}()

	C.runApp()
}