package diagnostics

import (
	"fmt"
	"io"
	"strings"
	"theseus-target/license-checker/internal/cli"
)

const Version = "25.0.1"
const MutualExclusion = "--failOn and --onlyAllow can not be used at the same time. Choose one or the other."

func HelpText() string {
	return "license-checker@25.0.1\n\n" +
		"   --production only show production dependencies.\n" +
		"   --development only show development dependencies.\n" +
		"   --unknown report guessed licenses as unknown licenses.\n" +
		"   --start [path of the initial json to look for]\n" +
		"   --onlyunknown only list packages with unknown or guessed licenses.\n" +
		"   --json output in json format.\n" +
		"   --csv output in csv format.\n" +
		"   --csvComponentPrefix column prefix for components in csv file\n" +
		"   --out [filepath] write the data to a specific file.\n" +
		"   --customPath to add a custom Format file in JSON\n" +
		"   --exclude [list] exclude modules which licenses are in the comma-separated list from the output\n" +
		"   --relativeLicensePath output the location of the license files as relative paths\n" +
		"   --summary output a summary of the license usage\n" +
		"   --failOn [list] fail (exit with code 1) on the first occurrence of the licenses of the semicolon-separated list\n" +
		"   --onlyAllow [list] fail (exit with code 1) on the first occurrence of the licenses not in the semicolon-seperated list\n" +
		"   --direct look for direct dependencies only\n" +
		"   --packages [list] restrict output to the packages (package@version) in the semicolon-seperated list\n" +
		"   --excludePackages [list] restrict output to the packages (package@version) not in the semicolon-seperated list\n" +
		"   --excludePrivatePackages restrict output to not include any package marked as private\n\n" +
		"   --version The current version\n" +
		"   --help  The text you are reading right now :)\n"
}

func Preflight(o cli.Options, stderr io.Writer) (handled bool, code int) {
	if o.Help {
		fmt.Fprint(stderr, HelpText())
		return true, 0
	}
	if o.Version {
		fmt.Fprintln(stderr, Version)
		return true, 1
	}
	if o.FailOn != "" && o.OnlyAllow != "" {
		fmt.Fprintln(stderr, MutualExclusion)
		return true, 1
	}
	if strings.Contains(o.FailOn, ",") {
		fmt.Fprintln(stderr, CommaWarning("failOn"))
	}
	if strings.Contains(o.OnlyAllow, ",") {
		fmt.Fprintln(stderr, CommaWarning("onlyAllow"))
	}
	return false, 0
}
func CommaWarning(flag string) string {
	return "Warning: As of v17 the --" + flag + " argument takes semicolons as delimeters instead of commas (some license names can contain commas)"
}
func FailOn(lic string) string {
	return fmt.Sprintf("Found license defined by the --failOn flag: %q. Exiting.", lic)
}
func OnlyAllow(item, lic string) string {
	return fmt.Sprintf("Package %q is licensed under %q which is not permitted by the --onlyAllow flag. Exiting.", item, lic)
}
