package cmd

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestCommandReference_NoLegacyListingsCommand(t *testing.T) {
	if commandByName(rootCmd, "listings") != nil {
		t.Fatalf("did not expect legacy listings command to be registered")
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

func TestCommandReference_SearchAllowsOptionalQueryArgs(t *testing.T) {
	search := mustCommandByName(t, rootCmd, "search")
	if err := search.Args(search, []string{}); err != nil {
		t.Fatalf("expected search command to allow no query args: %v", err)
	}
	if err := search.Args(search, []string{"red", "bike"}); err != nil {
		t.Fatalf("expected search command to allow query args: %v", err)
	}
}

func TestCommandReference_PostCreateFlags(t *testing.T) {
	post := mustCommandByName(t, rootCmd, "post")
	create := mustCommandByName(t, post, "create")

	for _, flagName := range []string{"category", "subcategory", "name", "body", "email", "price", "ip", "photo", "dry-run"} {
		if create.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected post create flag %q", flagName)
		}
	}
}

func TestCommandReference_PostRespondFlags(t *testing.T) {
	post := mustCommandByName(t, rootCmd, "post")
	respond := mustCommandByName(t, post, "respond")

	for _, flagName := range []string{"message", "reply-to", "ip", "dry-run"} {
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

func TestProjectStructure_ReadmeListedPathsExist(t *testing.T) {
	root := repoRoot(t)

	expectedFiles := []string{
		"AGENTS.md",
		"README.md",
		"Makefile",
		"main.go",
		"cmd/root.go",
		"cmd/version.go",
		"cmd/home.go",
		"cmd/search.go",
		"cmd/post.go",
		"cmd/post_create.go",
		"cmd/post_respond.go",
		"cmd/categories.go",
		"cmd/command_reference_test.go",
		"cmd/serve.go",
		"internal/config/config.go",
		"internal/domain/category.go",
		"internal/domain/category_rules.go",
		"internal/domain/home_category.go",
		"internal/domain/message.go",
		"internal/domain/post.go",
		"internal/domain/post_create_page.go",
		"internal/domain/post_create_submit.go",
		"internal/domain/post_respond.go",
		"internal/domain/search_result.go",
		"internal/domain/user.go",
		"internal/domain/errors.go",
		"internal/service/categories.go",
		"internal/service/home.go",
		"internal/service/post.go",
		"internal/service/post_create.go",
		"internal/service/post_create_submit.go",
		"internal/service/post_respond.go",
		"internal/service/search.go",
		"internal/repository/interfaces.go",
		"internal/repository/inmemory.go",
		"internal/repository/inmemory_post_create.go",
		"internal/repository/inmemory_post_respond.go",
		"internal/repository/inmemory_search.go",
		"internal/repository/postgres.go",
		"internal/repository/postgres_post_create.go",
		"internal/repository/postgres_post_respond.go",
		"internal/repository/postgres_search.go",
		"internal/adapters/output.go",
		"internal/adapters/mailgun.go",
		"internal/adapters/home_output.go",
		"internal/adapters/search_output.go",
		"internal/adapters/post_output.go",
		"internal/adapters/post_create_output.go",
		"internal/adapters/post_create_submit_output.go",
		"internal/adapters/post_respond_output.go",
		"internal/adapters/page_header.go",
		"internal/adapters/page_footer.go",
		"internal/adapters/home_cache.go",
		"internal/util/util.go",
		"configs/config.yaml.example",
		".env.example",
	}
	for _, rel := range expectedFiles {
		assertFileExists(t, filepath.Join(root, rel))
	}

	expectedDirs := []string{"supabase/migrations", "testdata/seed", "docs"}
	for _, rel := range expectedDirs {
		assertDirExists(t, filepath.Join(root, rel))
	}
}

func TestProjectStructure_LegacyListingPathsRemoved(t *testing.T) {
	root := repoRoot(t)
	removedPaths := []string{
		"cmd/listings.go",
		"internal/domain/listing.go",
		"internal/service/listings.go",
		"internal/service/listings_test.go",
		"testdata/seed/listings.json",
	}

	for _, rel := range removedPaths {
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			t.Fatalf("expected removed path %q to be absent", rel)
		} else if !os.IsNotExist(err) {
			t.Fatalf("unexpected stat error for %q: %v", rel, err)
		}
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

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve test file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected file at %q: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file at %q, got directory", path)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected directory at %q: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected directory at %q, got file", path)
	}
}
