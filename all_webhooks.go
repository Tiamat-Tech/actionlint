// Code generated by actionlint/scripts/generate-webhook-events. DO NOT EDIT.

package actionlint

// AllWebhookTypes is a table of all webhooks with their types. This variable was generated by
// script at ./scripts/generate-webhook-events based on
// https://github.com/github/docs/blob/main/content/actions/learn-github-actions/events-that-trigger-workflows.md
var AllWebhookTypes = map[string][]string{
	"branch_protection_rule":      {"created", "edited", "deleted"},
	"check_run":                   {"created", "rerequested", "completed"},
	"check_suite":                 {"completed", "requested", "rerequested"},
	"create":                      {},
	"delete":                      {},
	"deployment":                  {},
	"deployment_status":           {},
	"discussion":                  {"created", "edited", "deleted", "transferred", "pinned", "unpinned", "labeled", "unlabeled", "locked", "unlocked", "category_changed", "answered", "unanswered"},
	"discussion_comment":          {"created", "edited", "deleted"},
	"fork":                        {},
	"gollum":                      {},
	"issue_comment":               {"created", "edited", "deleted"},
	"issues":                      {"opened", "edited", "deleted", "transferred", "pinned", "unpinned", "closed", "reopened", "assigned", "unassigned", "labeled", "unlabeled", "locked", "unlocked", "milestoned", "demilestoned"},
	"label":                       {"created", "edited", "deleted"},
	"milestone":                   {"created", "closed", "opened", "edited", "deleted"},
	"page_build":                  {},
	"project":                     {"created", "updated", "closed", "reopened", "edited", "deleted"},
	"project_card":                {"created", "moved", "converted", "edited", "deleted"},
	"project_column":              {"created", "updated", "moved", "deleted"},
	"public":                      {},
	"pull_request":                {"assigned", "unassigned", "labeled", "unlabeled", "opened", "edited", "closed", "reopened", "synchronize", "converted_to_draft", "ready_for_review", "locked", "unlocked", "review_requested", "review_request_removed", "auto_merge_enabled", "auto_merge_disabled"},
	"pull_request_review":         {"submitted", "edited", "dismissed"},
	"pull_request_review_comment": {"created", "edited", "deleted"},
	"pull_request_target":         {"assigned", "unassigned", "labeled", "unlabeled", "opened", "edited", "closed", "reopened", "synchronize", "converted_to_draft", "ready_for_review", "locked", "unlocked", "review_requested", "review_request_removed", "auto_merge_enabled", "auto_merge_disabled"},
	"push":                        {},
	"registry_package":            {"published", "updated"},
	"release":                     {"published", "unpublished", "created", "edited", "deleted", "prereleased", "released"},
	"status":                      {},
	"watch":                       {"started"},
	"workflow_run":                {"completed", "requested"},
}
