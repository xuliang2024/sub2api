package service

import "testing"

func TestNormalizeCampaignHost(t *testing.T) {
	cases := map[string]string{
		"":                              "",
		"Try.Example.COM":               "try.example.com",
		"https://try.example.com/path":  "try.example.com",
		"http://try.example.com:8080/x": "try.example.com",
		"try.example.com:443":           "try.example.com",
		"try.example.com.":              "try.example.com",
		"  try.example.com  ":           "try.example.com",
		"[::1]:8443":                    "::1",
	}
	for input, want := range cases {
		if got := NormalizeCampaignHost(input); got != want {
			t.Errorf("NormalizeCampaignHost(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestValidateCampaignSlug(t *testing.T) {
	valid := []string{"a", "deepseek-v4", "promo2026", "x-1"}
	for _, slug := range valid {
		if err := ValidateCampaignSlug(slug); err != nil {
			t.Errorf("ValidateCampaignSlug(%q) unexpected error: %v", slug, err)
		}
	}
	invalid := []string{"", "-lead", "trail-", "UPPER", "with space", "中文", "a_b"}
	for _, slug := range invalid {
		if err := ValidateCampaignSlug(slug); err == nil {
			t.Errorf("ValidateCampaignSlug(%q) should fail", slug)
		}
	}
}

func TestValidateCampaignNormalizes(t *testing.T) {
	campaign := &Campaign{
		Slug:   "promo-1",
		Name:   "  促销活动  ",
		Status: "",
		Host:   "HTTPS://Try.Example.com:8443/landing",
		Trial:  CampaignTrial{Balance: -5},
	}
	if err := ValidateCampaign(campaign); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if campaign.Name != "促销活动" {
		t.Errorf("name not trimmed: %q", campaign.Name)
	}
	if campaign.Status != CampaignStatusDraft {
		t.Errorf("empty status should default to draft, got %q", campaign.Status)
	}
	if campaign.Host != "try.example.com" {
		t.Errorf("host not normalized: %q", campaign.Host)
	}
	if campaign.Trial.Balance != 0 {
		t.Errorf("negative trial balance should clamp to 0, got %v", campaign.Trial.Balance)
	}
}

func TestValidateCampaignRejectsBadInput(t *testing.T) {
	if err := ValidateCampaign(&Campaign{Slug: "ok", Name: "  "}); err == nil {
		t.Error("empty name should be rejected")
	}
	if err := ValidateCampaign(&Campaign{Slug: "ok", Name: "n", Status: "live"}); err == nil {
		t.Error("invalid status should be rejected")
	}
	if err := ValidateCampaign(&Campaign{Slug: "bad slug", Name: "n"}); err == nil {
		t.Error("invalid slug should be rejected")
	}
	if err := ValidateCampaign(&Campaign{Slug: "ok", Name: "n", CTA: CampaignCTA{Mode: "buy"}}); err == nil {
		t.Error("invalid cta mode should be rejected")
	}
}

func TestMergeCampaignModelsUsesConfiguredOrder(t *testing.T) {
	available := []CampaignModelView{
		{Model: "model-a", DisplayName: "model-a", Available: true},
		{Model: "model-b", DisplayName: "model-b", Available: true},
		{Model: "model-c", DisplayName: "model-c", Available: true},
	}
	entries := []CampaignModelEntry{
		{Model: "model-c", Order: 1},
		{Model: "model-a", Order: 2, Featured: true, DisplayName: "A 展示名", Badge: "主推"},
		{Model: "model-x", Order: 3, Tagline: "分组里没有"},
	}

	got := mergeCampaignModels(entries, available)
	if len(got) != 3 {
		t.Fatalf("expected 3 models, got %d", len(got))
	}
	// featured 优先
	if got[0].Model != "model-a" {
		t.Errorf("featured model should be first, got %q", got[0].Model)
	}
	if got[0].DisplayName != "A 展示名" || got[0].Badge != "主推" {
		t.Errorf("entry overrides not applied: %+v", got[0])
	}
	if !got[0].Available {
		t.Error("model-a should be marked available")
	}
	if got[1].Model != "model-c" {
		t.Errorf("second should be model-c by order, got %q", got[1].Model)
	}
	if got[2].Model != "model-x" || got[2].Available {
		t.Errorf("model-x should stay listed but unavailable: %+v", got[2])
	}
}

func TestMergeCampaignModelsWithoutEntriesReturnsAll(t *testing.T) {
	available := []CampaignModelView{
		{Model: "model-b", Available: true},
		{Model: "model-a", Available: true},
	}
	got := mergeCampaignModels(nil, available)
	if len(got) != 2 {
		t.Fatalf("expected 2 models, got %d", len(got))
	}
	if got[0].Model != "model-a" || got[1].Model != "model-b" {
		t.Errorf("models should be sorted by name, got %v/%v", got[0].Model, got[1].Model)
	}
	if empty := mergeCampaignModels(nil, nil); empty != nil {
		t.Errorf("expected nil for empty input, got %v", empty)
	}
}
