package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// splitImportID splits an import ID by "/" and returns exactly n parts,
// or adds an error diagnostic if the count doesn't match.
func splitImportID(resp *resource.ImportStateResponse, id string, n int, example string) ([]string, diag.Diagnostics) {
	parts := strings.SplitN(id, "/", n)
	if len(parts) != n {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected %d parts separated by '/' (e.g. %q), got %q.", n, example, id),
		)
		return nil, resp.Diagnostics
	}
	for i, p := range parts {
		if p == "" {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Part %d of import ID %q must not be empty.", i+1, id),
			)
			return nil, resp.Diagnostics
		}
	}
	return parts, nil
}

// setattrs is a helper to set multiple path/value pairs in import state.
func setAttrs(ctx context.Context, resp *resource.ImportStateResponse, attrs map[string]string) {
	for p, v := range attrs {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(p), v)...)
	}
}
