package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandReference_TopLevelCommandsExist(t *testing.T) {
	for _, name := range []string{"home", "search", "post", "categories", "serve", "version"} {
		if mustCommandByName(t, rootCmd, name) == nil {
			t.Fatalf("expected top-level command %q", name)
		}
	}
}

func TestCommandReference_GlobalFlagsExist(t *testing.T) {
	verbose := rootCmd.PersistentFlags().Lookup("verbose")
	if verbose == nil {
		t.Fatalf("expected --verbose global flag")
	}
	if verbose.Shorthand != "v" {
		t.Fatalf("expected --verbose shorthand -v, got %q", verbose.Shorthand)
	}

	format := rootCmd.PersistentFlags().Lookup("format")
	if format == nil {
		t.Fatalf("expected --format global flag")
	}
	if format.DefValue != "json" {
		t.Fatalf("expected --format default json, got %q", format.DefValue)
	}

	config := rootCmd.PersistentFlags().Lookup("config")
	if config == nil {
		t.Fatalf("expected --config global flag")
	}
	if config.DefValue != ".supost.yaml" {
		t.Fatalf("expected --config default .supost.yaml, got %q", config.DefValue)
	}
}

func TestCommandReference_SearchFlags(t *testing.T) {
	search := mustCommandByName(t, rootCmd, "search")

	category := search.Flags().Lookup("category")
	if category == nil || category.DefValue != "0" {
		t.Fatalf("expected search --category with default 0")
	}

	subcategory := search.Flags().Lookup("subcategory")
	if subcategory == nil || subcategory.DefValue != "0" {
		t.Fatalf("expected search --subcategory with default 0")
	}

	page := search.Flags().Lookup("page")
	if page == nil || page.DefValue != "1" {
		t.Fatalf("expected search --page with default 1")
	}

	perPage := search.Flags().Lookup("per-page")
	if perPage == nil || perPage.DefValue != "100" {
		t.Fatalf("expected search --per-page with default 100")
	}
}

func TestCommandReference_PostCreateFlags(t *testing.T) {
	post := mustCommandByName(t, rootCmd, "post")
	create := mustCommandByName(t, post, "create")

	for _, flagName := range []string{"category", "subcategory", "name", "body", "email", "price", "dry-run"} {
		if create.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected post create flag %q", flagName)
		}
	}
}

func TestCommandReference_PostRespondFlags(t *testing.T) {
	post := mustCommandByName(t, rootCmd, "post")
	respond := mustCommandByName(t, post, "respond")

	for _, flagName := range []string{"message", "reply-to", "dry-run"} {
		if respond.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected post respond flag %q", flagName)
		}
	}

	if !isRequiredFlag(respond, "message") {
		t.Fatalf("expected post respond --message to be required")
	}
	if !isRequiredFlag(respond, "reply-to") {
		t.Fatalf("expected post respond --reply-to to be required")
	}
}

func TestCommandReference_PostAndRespondArgs(t *testing.T) {
	post := mustCommandByName(t, rootCmd, "post")
	if err := post.Args(post, []string{}); err == nil {
		t.Fatalf("expected post command to require <post_id>")
	}
	if err := post.Args(post, []string{"130031605"}); err != nil {
		t.Fatalf("expected post command to accept a single <post_id>: %v", err)
	}

	respond := mustCommandByName(t, post, "respond")
	if err := respond.Args(respond, []string{}); err == nil {
		t.Fatalf("expected post respond command to require <post_id>")
	}
	if err := respond.Args(respond, []string{"130031802"}); err != nil {
		t.Fatalf("expected post respond command to accept a single <post_id>: %v", err)
	}
}

func TestCommandReference_ServePortDefault(t *testing.T) {
	serve := mustCommandByName(t, rootCmd, "serve")
	port := serve.Flags().Lookup("port")
	if port == nil {
		t.Fatalf("expected serve --port flag")
	}
	if port.DefValue != "8080" {
		t.Fatalf("expected serve --port default 8080, got %q", port.DefValue)
	}
}

func mustCommandByName(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()
	cmd := commandByName(parent, name)
	if cmd == nil {
		t.Fatalf("command %q not found under %q", name, parent.Name())
	}
	return cmd
}

func commandByName(parent *cobra.Command, name string) *cobra.Command {
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

func isRequiredFlag(cmd *cobra.Command, flagName string) bool {
	flag := cmd.Flags().Lookup(flagName)
	if flag == nil {
		return false
	}
	required, ok := flag.Annotations[cobra.BashCompOneRequiredFlag]
	return ok && len(required) > 0
}
